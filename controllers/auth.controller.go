package controllers

import (
	"errors"
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/dto"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/utils"
	"net/http"
	"time"

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
			"data":    nil,
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
			"data":    nil,
		})
		return
	}

	// signing jwt
	accessToken, err := utils.GenerateJwt(authUser.Id, false, 3*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	refreshToken, err := utils.GenerateJwt(authUser.Id, false, 3*24*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login",
		"data": gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
			"createdAt":    time.Now(),
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
	user := models.User{
		BaseModel: models.BaseModel{Id: randomUuid},
		Name:      registerDto.Name,
		Email:     registerDto.Email,
		Password:  string(hashedPassword),
	}
	database.DB.Save(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "register",
		"data":    user,
	})
}

func AuthUser(ctx *gin.Context) {
	user, _ := ctx.Get("authUser")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "get auth user",
		"data":    user,
	})
}
