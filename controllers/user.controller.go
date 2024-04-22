package controllers

import (
	"errors"
	"fmt"
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/dto"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	paramId := ctx.Param("id")
	var getUserDto dto.GetUserDto
	ctx.ShouldBind(&getUserDto)
	getUserDto.Id = paramId
	validationErrors := utils.Validate(getUserDto)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	var user models.User
	id, _ := uuid.Parse(paramId)
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
	ctx.ShouldBind(&createUserDto)
	validationErrors := utils.Validate(createUserDto)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	// save to database
	hashedPassword, _ := utils.HashPassword(createUserDto.Password)
	randomUuid, err := uuid.NewRandom()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	user := models.User{
		BaseModel: models.BaseModel{Id: randomUuid},
		Name:      createUserDto.Name,
		Email:     createUserDto.Email,
		Password:  string(hashedPassword),
	}
	database.DB.Save(&user)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "create user",
		"data":    user,
	})
}

func UpdateUser(ctx *gin.Context) {
	paramId := ctx.Param("id")
	var updateUserDto dto.UpdateUserDto
	ctx.ShouldBind(&updateUserDto)
	updateUserDto.Id = paramId
	validationErrors := utils.Validate(updateUserDto)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	// save to database
	var user models.User
	id, _ := uuid.Parse(paramId)
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
	paramId := ctx.Param("id")
	var getUserDto dto.GetUserDto
	ctx.ShouldBind(&getUserDto)
	getUserDto.Id = paramId
	validationErrors := utils.Validate(getUserDto)
	if validationErrors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	var user models.User
	id, _ := uuid.Parse(paramId)
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
