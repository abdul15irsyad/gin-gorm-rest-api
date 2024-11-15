package handlers

import (
	"gin-gorm-rest-api/configs"
	"gin-gorm-rest-api/dtos"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/services"
	"gin-gorm-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthUserHandler struct {
	userService    *services.UserService
	databaseConfig *configs.DatabaseConfig
}

func NewAuthUserHandler(userService *services.UserService, databaseConfig *configs.DatabaseConfig) *AuthUserHandler {
	return &AuthUserHandler{userService: userService, databaseConfig: databaseConfig}
}

func (auh *AuthUserHandler) AuthUser(ctx *gin.Context) {
	authUser, _ := ctx.Get("authUser")
	user := authUser.(models.User)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "get auth user",
		"data":    user,
	})
}

func (auh *AuthUserHandler) UpdateAuthUser(ctx *gin.Context) {
	var updateAuthUserDto dtos.UpdateAuthUserDto
	ctx.ShouldBind(&updateAuthUserDto)
	validationErrors := utils.Validate(updateAuthUserDto)
	if len(validationErrors) > 0 {
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
	auh.databaseConfig.DB.Save(&user)
	user, _ = auh.userService.GetUser(user.Id)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "update auth user",
		"data":    user,
	})
}

func (auh *AuthUserHandler) UpdateAuthUserPassword(ctx *gin.Context) {
	var updateAuthUserPasswordDto dtos.UpdateAuthUserPasswordDto
	ctx.ShouldBind(&updateAuthUserPasswordDto)
	validationErrors := utils.Validate(updateAuthUserPasswordDto)
	if len(validationErrors) > 0 {
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
	user.Password = hashedPassword
	auh.databaseConfig.DB.Save(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "update auth user password",
	})
}
