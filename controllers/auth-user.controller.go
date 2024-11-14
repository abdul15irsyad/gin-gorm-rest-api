package controllers

import (
	"gin-gorm-rest-api/config"
	"gin-gorm-rest-api/dto"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/services"
	"gin-gorm-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthUserController struct {
	userService    *services.UserService
	databaseConfig *config.DatabaseConfig
}

func NewAuthUserController(userService *services.UserService, databaseConfig *config.DatabaseConfig) *AuthUserController {
	return &AuthUserController{userService: userService, databaseConfig: databaseConfig}
}

func (auc *AuthUserController) AuthUser(ctx *gin.Context) {
	authUser, _ := ctx.Get("authUser")
	user := authUser.(models.User)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "get auth user",
		"data":    user,
	})
}

func (auc *AuthUserController) UpdateAuthUser(ctx *gin.Context) {
	var updateAuthUserDto dto.UpdateAuthUserDto
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
	auc.databaseConfig.DB.Save(&user)
	user, _ = auc.userService.GetUser(user.Id)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "update auth user",
		"data":    user,
	})
}

func (auc *AuthUserController) UpdateAuthUserPassword(ctx *gin.Context) {
	var updateAuthUserPasswordDto dto.UpdateAuthUserPasswordDto
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
	user.Password = string(hashedPassword)
	auc.databaseConfig.DB.Save(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "update auth user password",
	})
}
