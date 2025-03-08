package handlers

import (
	"errors"
	"gin-gorm-rest-api/dtos"
	"gin-gorm-rest-api/middlewares"
	"gin-gorm-rest-api/services"
	"gin-gorm-rest-api/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	jwtService     *services.JwtService
	mailService    *services.MailService
	userService    *services.UserService
	roleService    *services.RoleService
	tokenService   *services.TokenService
	authMiddleware *middlewares.AuthMiddleware
}

func NewAuthHandler(
	jwtService *services.JwtService,
	mailService *services.MailService,
	userService *services.UserService,
	roleService *services.RoleService,
	tokenService *services.TokenService,
	authMiddleware *middlewares.AuthMiddleware,
) *AuthHandler {
	return &AuthHandler{jwtService, mailService, userService, roleService, tokenService, authMiddleware}
}

func (ah *AuthHandler) Login(ctx *gin.Context) {
	var loginDto dtos.LoginDto
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
	authUser, err := ah.userService.GetUserCredential(loginDto.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// for decoy
			utils.ComparePassword("some password", loginDto.Password)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "email or password is incorrect",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
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
	accessToken, refreshToken, ok := ah.jwtService.SigningToken(ctx, authUser)
	if !ok {
		return
	}

	ctx.SetCookie("accessToken", accessToken, 0, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login",
		"data": gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
			"grantType":    "credential",
		},
	})
}

func (ah *AuthHandler) Logout(ctx *gin.Context) {
	ctx.SetCookie("accessToken", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "logout",
	})
}

func (ah *AuthHandler) Register(ctx *gin.Context) {
	var registerDto dtos.RegisterDto
	ctx.ShouldBind(&registerDto)
	validationErrors := utils.Validate(registerDto)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	userRole, _ := ah.roleService.GetRoleBy(dtos.GetDataByOptions{
		Field: "name",
		Value: "user",
	})

	user, err := ah.userService.CreateUser(dtos.CreateUserDto{
		Name:     registerDto.Name,
		Email:    registerDto.Email,
		Password: registerDto.Password,
		RoleId:   userRole.Id.String(),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "register",
		"data":    user,
	})
}

func (ah *AuthHandler) RefreshToken(ctx *gin.Context) {
	authUser, ok := ah.authMiddleware.GetAuthUserFromAuthorization(ctx, "refresh")
	if !ok {
		return
	}

	// signing jwt
	accessToken, refreshToken, ok := ah.jwtService.SigningToken(ctx, authUser)
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

func (ah *AuthHandler) ForgotPassword(ctx *gin.Context) {
	// get request body
	var forgotPassword dtos.ForgotPasswordDto
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
	user, err := ah.userService.GetUserBy(dtos.GetDataByOptions{
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

	token, err := ah.tokenService.CreateToken(user.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// send link to reset password
	url := os.Getenv("BASE_URL") + "/auth/reset-password?token=" + token.Token
	go ah.mailService.SendMail(services.SendMailOptions{
		Subject: "Forgot Password",
		To:      user.Email,
		Message: "<a href=\"" + url + "\">" + url + "</a>",
	})

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

func (ah *AuthHandler) ResetPassword(ctx *gin.Context) {
	tokenString := ctx.Query("token")
	if tokenString == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid token",
		})
		return
	}
	token, err := ah.tokenService.GetToken(tokenString)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid token",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
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
	var resetPassword dtos.ResetPasswordDto
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
	user, err := ah.userService.GetUser(token.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = ah.userService.UpdateUserPassword(user, resetPassword.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "reset password",
	})
}
