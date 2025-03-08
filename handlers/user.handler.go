package handlers

import (
	"errors"
	"fmt"
	"gin-gorm-rest-api/dtos"
	"gin-gorm-rest-api/services"
	"gin-gorm-rest-api/utils"
	"gin-gorm-rest-api/utils/validations"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService *services.UserService
	fileService *services.FileService
}

func NewUserHandler(userService *services.UserService, fileService *services.FileService) *UserHandler {
	return &UserHandler{userService, fileService}
}

func (uh *UserHandler) GetAllUsers(ctx *gin.Context) {
	var getUsersDto dtos.GetUsersDto
	ctx.ShouldBind(&getUsersDto)
	validationErrors := utils.Validate(getUsersDto)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}
	if getUsersDto.Page == nil {
		defaultPage := 1
		getUsersDto.Page = &defaultPage
	}
	if getUsersDto.Limit == nil {
		defaultLimit := 10
		getUsersDto.Limit = &defaultLimit
	}
	users, total, totalPage, err := uh.userService.GetPaginatedUsers(getUsersDto)
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
			"currentPage": *getUsersDto.Page,
			"pageSize":    *getUsersDto.Limit,
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

	user, err := uh.userService.CreateUser(createUserDto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

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
	// validate image
	if updateUserDto.Image != nil {
		validations.ValidateImage(validations.ValidateImageOptions{
			Validations: []struct {
				Validate func() bool
				Message  string
				Tag      string
			}{
				{
					Validate: func() bool {
						return validations.MaxFileSize(validations.MaxFileSizeOptions{
							File:    updateUserDto.Image,
							MaxSize: &validations.DefaultMaxSize,
						})
					},
					Message: "Key: 'UpdateUserDto.Image' Error:Field validation for 'Image' failed on the 'maxFileSize' tag",
					Tag:     fmt.Sprintf("maxFileSize=%dKB", validations.DefaultMaxSize/1024),
				},
				{
					Validate: func() bool {
						return validations.FileMime(validations.FileMimeOptions{
							File:         updateUserDto.Image,
							AllowedMimes: []string{"image/jpg", "image/png"},
						})
					},
					Message: "Key: 'UpdateUserDto.Image' Error:Field validation for 'Image' failed on the 'fileMime' tag",
					Tag:     fmt.Sprintf("fileMime=%s", strings.Join([]string{"image/jpg", "image/png"}, ",")),
				},
			},
			Field:           "image",
			Value:           updateUserDto.Image,
			ErrorValidation: &validationErrors,
		})
	}
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
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "data not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// save uploaded file
	var newImageFileId *uuid.UUID = nil
	if updateUserDto.Image != nil {
		newFile, ok := uh.fileService.UploadAndCreateFile(ctx, updateUserDto.Image)
		if !ok {
			return
		}
		newImageFileId = &newFile.Id
	}

	user, err = uh.userService.UpdateUser(id, updateUserDto, newImageFileId)

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

	uh.userService.DeleteUser(user.Id)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "delete user",
	})
}
