package controllers

import (
	"belajar-gin/database"
	"belajar-gin/dto"
	"belajar-gin/models"
	"belajar-gin/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllUser(ctx *gin.Context) {
	var users []models.User
	database.DB.Find(&users)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "get all user",
		"data":    users,
	})
}

func GetUser(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var user models.User
	result := database.DB.First(&user, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "get user",
		"data":    user,
	})
}

func CreateUser(ctx *gin.Context) {
	var createUserDto dto.CreateUserDto
	ctx.ShouldBindJSON(&createUserDto)
	validationErrors := utils.Validate(createUserDto)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	hashedPassword, _ := utils.HashPassword(createUserDto.Password)
	user := models.User{
		Name:     createUserDto.Name,
		Email:    createUserDto.Email,
		Password: string(hashedPassword),
	}
	database.DB.Save(&user)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "create user",
		"data":    user,
	})
}

func UpdateUser(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var updateUserDto dto.UpdateUserDto
	ctx.ShouldBindJSON(&updateUserDto)
	validationErrors := utils.Validate(updateUserDto)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	var user models.User
	result := database.DB.First(&user, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
			"data":    nil,
		})
		return
	}

	user.Name = updateUserDto.Name
	user.Email = updateUserDto.Email
	if updateUserDto.Password != nil {
		fmt.Print(updateUserDto.Password)
		hashedPassword, _ := utils.HashPassword(*updateUserDto.Password)
		user.Password = string(hashedPassword)
	}
	database.DB.Save(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "update user",
		"data":    user,
	})
}

func DeleteUser(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var user models.User
	result := database.DB.First(&user, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
			"data":    nil,
		})
		return
	}

	database.DB.Delete(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "delete user",
		"data":    nil,
	})
}
