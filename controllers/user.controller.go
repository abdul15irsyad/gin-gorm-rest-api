package controllers

import (
	"errors"
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/dto"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetAllUsers(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	search := ctx.Query("search")
	users, total, totalPage, err := models.GetPaginatedUsers(database.DB, page, limit, &search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "get all user",
		"data":    users,
		"meta": gin.H{
			"currentPage": page,
			"pageSize":    limit,
			"totalPages":  totalPage,
			"totalItems":  total,
		},
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
	result, err := models.GetUser(database.DB, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}
	user = *result

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
	hashedPassword, err := utils.HashPassword(createUserDto.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
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
	result, err := models.GetUser(database.DB, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}
	user = *result
	user.Name = updateUserDto.Name
	user.Email = updateUserDto.Email
	if updateUserDto.Password != nil {
		hashedPassword, err := utils.HashPassword(*updateUserDto.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
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
	result, err := models.GetUser(database.DB, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}
	user = *result

	database.DB.Delete(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "delete user",
	})
}
