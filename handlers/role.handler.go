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
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type RoleHandler struct {
	roleService    *services.RoleService
	databaseConfig *configs.DatabaseConfig
}

func NewRoleHandler(roleService *services.RoleService, databaseConfig *configs.DatabaseConfig) *RoleHandler {
	return &RoleHandler{roleService: roleService, databaseConfig: databaseConfig}
}

func (rh *RoleHandler) GetAllRoles(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	search := ctx.Query("search")
	roles, total, totalPage, err := rh.roleService.GetPaginatedRoles(page, limit, &search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "get all role",
		"data":    roles,
		"meta": gin.H{
			"currentPage": page,
			"pageSize":    limit,
			"totalPages":  totalPage,
			"totalItems":  total,
		},
	})
}

func (rh *RoleHandler) GetRole(ctx *gin.Context) {
	paramId := ctx.Param("id")
	var getRoleDto dtos.GetRoleDto
	ctx.ShouldBind(&getRoleDto)
	getRoleDto.Id = paramId
	validationErrors := utils.Validate(getRoleDto)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	id, _ := uuid.Parse(paramId)
	role, err := rh.roleService.GetRole(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "get role",
		"data":    role,
	})
}

func (rh *RoleHandler) CreateRole(ctx *gin.Context) {
	var createRoleDto dtos.CreateRoleDto
	ctx.ShouldBind(&createRoleDto)
	validationErrors := utils.Validate(createRoleDto)
	// check is name unique in database
	nameErrorExists := false
	for _, validationError := range validationErrors {
		if validationError.Field == "Name" {
			nameErrorExists = true
		}
	}
	if !nameErrorExists {
		_, err := rh.roleService.GetRoleBy(dtos.GetDataByOptions{
			Field:     "slug",
			Value:     slug.Make(createRoleDto.Name),
			ExcludeId: nil,
		})
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			validationErrors = append(validationErrors, utils.ErrorResponse{
				Field:   "Name",
				Message: "Key: 'createRoleDto.Name' Error:Field validation for 'Name' failed on the 'unique' tag",
				Tag:     "unique",
				Value:   createRoleDto.Name,
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
	randomUuid, err := uuid.NewRandom()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	role := models.Role{
		BaseModel: models.BaseModel{Id: randomUuid},
		Name:      createRoleDto.Name,
		Slug:      slug.Make(createRoleDto.Name),
		Desc:      createRoleDto.Desc,
	}
	rh.databaseConfig.DB.Save(&role)
	role, _ = rh.roleService.GetRole(role.Id)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "create role",
		"data":    role,
	})
}

func (rh *RoleHandler) UpdateRole(ctx *gin.Context) {
	paramId := ctx.Param("id")
	var updateRoleDto dtos.UpdateRoleDto
	ctx.ShouldBind(&updateRoleDto)
	updateRoleDto.Id = paramId
	validationErrors := utils.Validate(updateRoleDto)
	// check is name unique in database
	nameErrorExists := false
	for _, validationError := range validationErrors {
		if validationError.Field == "Name" {
			nameErrorExists = true
		}
	}
	if !nameErrorExists {
		_, err := rh.roleService.GetRoleBy(dtos.GetDataByOptions{
			Field:     "slug",
			Value:     slug.Make(updateRoleDto.Name),
			ExcludeId: &updateRoleDto.Id,
		})
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			validationErrors = append(validationErrors, utils.ErrorResponse{
				Field:   "Name",
				Message: "Key: 'updateRoleDto.Name' Error:Field validation for 'Name' failed on the 'unique' tag",
				Tag:     "unique",
				Value:   updateRoleDto.Name,
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
	role, err := rh.roleService.GetRole(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}
	role.Name = updateRoleDto.Name
	role.Slug = slug.Make(updateRoleDto.Name)
	role.Desc = updateRoleDto.Desc
	rh.databaseConfig.DB.Save(&role)
	role, _ = rh.roleService.GetRole(role.Id)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "update role",
		"data":    role,
	})
}

func (rh *RoleHandler) DeleteRole(ctx *gin.Context) {
	paramId := ctx.Param("id")
	var getRoleDto dtos.GetRoleDto
	ctx.ShouldBind(&getRoleDto)
	getRoleDto.Id = paramId
	validationErrors := utils.Validate(getRoleDto)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	id, _ := uuid.Parse(paramId)
	role, err := rh.roleService.GetRole(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}

	rh.databaseConfig.DB.Delete(&role)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "delete role",
	})
}
