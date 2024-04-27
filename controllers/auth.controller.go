package controllers

import (
	"errors"
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/dto"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Login(ctx *gin.Context) {
	var loginDto dto.LoginDto
	ctx.ShouldBind(&loginDto)
	validationErrors := utils.Validate(loginDto)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	// verify credential
	var authUser models.User
	result := database.DB.Select([]string{"id", "email", "password"}).First(&authUser, "email = ?", loginDto.Email)
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
	if validationErrors != nil {
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

func AuthUser(ctx *gin.Context) {
	authUser, _ := ctx.Get("authUser")
	user := authUser.(models.User)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "get auth user",
		"data":    user,
	})
}

func UpdateAuthUser(ctx *gin.Context) {
	var updateAuthUserDto dto.UpdateAuthUserDto
	ctx.ShouldBind(&updateAuthUserDto)
	validationErrors := utils.Validate(updateAuthUserDto)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	authUser, _ := ctx.Get("authUser")
	user := authUser.(models.User)
	user.Name = updateAuthUserDto.Name
	user.Email = updateAuthUserDto.Email
	database.DB.Save(&user)
	user, _ = models.GetUser(database.DB, user.Id)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "update auth user",
		"data":    user,
	})
}

func UpdateAuthUserPassword(ctx *gin.Context) {
	var updateAuthUserPasswordDto dto.UpdateAuthUserPasswordDto
	ctx.ShouldBind(&updateAuthUserPasswordDto)
	validationErrors := utils.Validate(updateAuthUserPasswordDto)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	authUser, _ := ctx.Get("authUser")
	user := authUser.(models.User)
	correctPassword, err := utils.ComparePassword(user.Password, updateAuthUserPasswordDto.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	if !correctPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "password is incorrect",
		})
		return
	}
	hashedPassword, err := utils.HashPassword(updateAuthUserPasswordDto.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	user.Password = string(hashedPassword)
	database.DB.Save(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "update auth user password",
	})
}
