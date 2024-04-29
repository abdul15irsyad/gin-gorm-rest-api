package controllers

import (
	"errors"
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/dto"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Login(ctx *gin.Context) {
	var loginDto dto.LoginDto
	ctx.ShouldBind(&loginDto)
	validationErrors := utils.Validate(loginDto)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	// verify credential
	var authUser models.User
	result := database.DB.Select([]string{"id", "password"}).First(&authUser, "email = ?", loginDto.Email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// for decoy
		utils.ComparePassword("some password", loginDto.Password)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "email or password is incorrect",
		})
		return
	}

	correctPassword, err := utils.ComparePassword(authUser.Password, loginDto.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	if !correctPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "email or password is incorrect",
		})
		return
	}

	// signing jwt
	accessToken, refreshToken, ok := utils.SigningToken(ctx, authUser)
	if !ok {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login",
		"data": gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
			"grantType":    "credential",
		},
	})
}

func Register(ctx *gin.Context) {
	var registerDto dto.RegisterDto
	ctx.ShouldBind(&registerDto)
	validationErrors := utils.Validate(registerDto)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	// save to database
	hashedPassword, _ := utils.HashPassword(registerDto.Password)
	randomUuid, err := uuid.NewRandom()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	userRoleId, _ := uuid.Parse("3ed4e622-4642-499a-b711-fb86a458f098")
	user := models.User{
		BaseModel: models.BaseModel{Id: randomUuid},
		Name:      registerDto.Name,
		Email:     registerDto.Email,
		Password:  string(hashedPassword),
		RoleId:    userRoleId,
	}
	database.DB.Save(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "register",
		"data":    user,
	})
}

func RefreshToken(ctx *gin.Context) {
	authUser, ok := utils.GetAuthUserFromAuthorization(ctx, "refresh")
	if !ok {
		return
	}

	// signing jwt
	accessToken, refreshToken, ok := utils.SigningToken(ctx, authUser)
	if !ok {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "refresh",
		"data": gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
			"grantType":    "refresh token",
		},
	})
}

func ForgotPassword(ctx *gin.Context) {
	// get request body
	var forgotPassword dto.ForgotPasswordDto
	ctx.ShouldBind(&forgotPassword)
	validationErrors := utils.Validate(forgotPassword)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	// check user
	user, err := models.GetUserBy(models.GetDataByOptions{
		DB:        database.DB,
		Field:     "email",
		Value:     forgotPassword.Email,
		ExcludeId: nil,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "there is no user with email '" + forgotPassword.Email + "'",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// create token
	var randomString string
	for {
		var token models.Token
		randomString = utils.GenerateRandomString(64)
		result := database.DB.Where("token = ?", randomString).First(&token)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			break
		}
	}
	randomUuid, _ := uuid.NewRandom()
	token := models.Token{
		BaseModel: models.BaseModel{Id: randomUuid},
		Token:     randomString,
		Type:      models.TokenForgotPassword,
		UserId:    user.Id,
		ExpiredAt: time.Now().Add(time.Hour),
	}
	result := database.DB.Save(&token)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
		return
	}
	token, _ = models.GetToken(database.DB, token.Id)

	// send link to reset password
	url := os.Getenv("BASE_URL") + "/auth/reset-password?token=" + token.Token
	err = utils.SendMail(utils.SendMailOptions{
		Subject: "Forgot Password",
		To:      user.Email,
		Message: "<a href=\"" + url + "\">" + url + "</a>",
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if Env := os.Getenv("ENV"); Env == "production" {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "forgot password",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "forgot password",
			"data": gin.H{
				"token": token.Token,
			},
		})
	}
}

func ResetPassword(ctx *gin.Context) {
	tokenString := ctx.Query("token")
	if tokenString == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid token",
		})
		return
	}
	var token models.Token
	result := database.DB.Where("token = ? AND type = ?", tokenString, models.TokenForgotPassword).First(&token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid token",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
		return
	}
	if token.ExpiredAt.Before(time.Now()) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    "TOKEN_EXPIRED",
			"message": "token expired",
		})
		return
	}

	// get request body
	var resetPassword dto.ResetPasswordDto
	ctx.ShouldBind(&resetPassword)
	validationErrors := utils.Validate(resetPassword)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	// change password in database
	user, err := models.GetUser(database.DB, token.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	hashedPassword, _ := utils.HashPassword(resetPassword.Password)
	user.Password = string(hashedPassword)
	result = database.DB.Save(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "reset password",
	})
}
