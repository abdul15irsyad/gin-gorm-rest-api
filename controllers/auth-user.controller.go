package controllers

import (
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/dto"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	database.DB.Save(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "update auth user password",
	})
}
