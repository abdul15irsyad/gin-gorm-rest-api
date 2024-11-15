package handlers

import (
	"errors"
	"gin-gorm-rest-api/configs"
	"gin-gorm-rest-api/dtos"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/services"
	"gin-gorm-rest-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService    *services.UserService
	databaseConfig *configs.DatabaseConfig
}

func NewUserHandler(userService *services.UserService, databaseConfig *configs.DatabaseConfig) *UserHandler {
	return &UserHandler{userService, databaseConfig}
}

func (uh *UserHandler) GetAllUsers(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	search := ctx.Query("search")
	users, total, totalPage, err := uh.userService.GetPaginatedUsers(page, limit, &search)
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

func (uh *UserHandler) GetUser(ctx *gin.Context) {
	paramId := ctx.Param("id")
	var getUserDto dtos.GetUserDto
	ctx.ShouldBind(&getUserDto)
	getUserDto.Id = paramId
	validationErrors := utils.Validate(getUserDto)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	id, _ := uuid.Parse(paramId)
	user, err := uh.userService.GetUser(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "get user",
		"data":    user,
	})
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var createUserDto dtos.CreateUserDto
	ctx.ShouldBind(&createUserDto)
	validationErrors := utils.Validate(createUserDto)
	// check is email unique in database
	emailErrorExists := false
	for _, validationError := range validationErrors {
		if validationError.Field == "Email" {
			emailErrorExists = true
		}
	}
	if !emailErrorExists {
		_, err := uh.userService.GetUserBy(dtos.GetDataByOptions{
			Field:     "email",
			Value:     createUserDto.Email,
			ExcludeId: nil,
		})
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			validationErrors = append(validationErrors, utils.ErrorResponse{
				Field:   "Email",
				Message: "Key: 'CreateUserDto.Email' Error:Field validation for 'Email' failed on the 'unique' tag",
				Tag:     "unique",
				Value:   createUserDto.Email,
			})
		}
	}
	if len(validationErrors) > 0 {
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
	roleId, _ := uuid.Parse(createUserDto.RoleId)
	user := models.User{
		BaseModel: models.BaseModel{Id: randomUuid},
		Name:      createUserDto.Name,
		Email:     createUserDto.Email,
		Password:  hashedPassword,
		RoleId:    roleId,
	}
	result := uh.databaseConfig.DB.Save(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
		return
	}
	user, _ = uh.userService.GetUser(user.Id)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "create user",
		"data":    user,
	})
}

func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	paramId := ctx.Param("id")
	var updateUserDto dtos.UpdateUserDto
	ctx.ShouldBind(&updateUserDto)
	updateUserDto.Id = paramId
	validationErrors := utils.Validate(updateUserDto)
	// check is email unique in database
	emailErrorExists := false
	for _, validationError := range validationErrors {
		if validationError.Field == "Email" {
			emailErrorExists = true
		}
	}
	if !emailErrorExists {
		_, err := uh.userService.GetUserBy(dtos.GetDataByOptions{
			Field:     "email",
			Value:     updateUserDto.Email,
			ExcludeId: &updateUserDto.Id,
		})
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			validationErrors = append(validationErrors, utils.ErrorResponse{
				Field:   "Email",
				Message: "Key: 'UpdateUserDto.Email' Error:Field validation for 'Email' failed on the 'unique' tag",
				Tag:     "unique",
				Value:   updateUserDto.Email,
			})
		}
	}
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	// save to database
	id, _ := uuid.Parse(paramId)
	user, err := uh.userService.GetUser(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}
	user.Name = updateUserDto.Name
	user.Email = updateUserDto.Email
	roleId, _ := uuid.Parse(updateUserDto.RoleId)
	user.RoleId = roleId
	if updateUserDto.Password != nil {
		hashedPassword, err := utils.HashPassword(*updateUserDto.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		user.Password = hashedPassword
	}
	uh.databaseConfig.DB.Save(&user)
	user, _ = uh.userService.GetUser(user.Id)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "update user",
		"data":    user,
	})
}

func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	paramId := ctx.Param("id")
	var getUserDto dtos.GetUserDto
	ctx.ShouldBind(&getUserDto)
	getUserDto.Id = paramId
	validationErrors := utils.Validate(getUserDto)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	id, _ := uuid.Parse(paramId)
	user, err := uh.userService.GetUser(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}

	uh.databaseConfig.DB.Delete(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "delete user",
	})
}
