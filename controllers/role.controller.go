package controllers

import (
	"errors"
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/dto"
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

type RoleController struct {
	roleService *services.RoleService
}

func NewRoleController(roleService *services.RoleService) *RoleController {
	return &RoleController{roleService: roleService}
}

func (rc *RoleController) GetAllRoles(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	search := ctx.Query("search")
	roles, total, totalPage, err := rc.roleService.GetPaginatedRoles(page, limit, &search)
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

func (rc *RoleController) GetRole(ctx *gin.Context) {
	paramId := ctx.Param("id")
	var getRoleDto dto.GetRoleDto
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
	role, err := rc.roleService.GetRole(id)
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

func (rc *RoleController) CreateRole(ctx *gin.Context) {
	var createRoleDto dto.CreateRoleDto
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
		_, err := rc.roleService.GetRoleBy(dto.GetDataByOptions{
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
	database.DB.Save(&role)
	role, _ = rc.roleService.GetRole(role.Id)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "create role",
		"data":    role,
	})
}

func (rc *RoleController) UpdateRole(ctx *gin.Context) {
	paramId := ctx.Param("id")
	var updateRoleDto dto.UpdateRoleDto
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
		_, err := rc.roleService.GetRoleBy(dto.GetDataByOptions{
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
	role, err := rc.roleService.GetRole(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}
	role.Name = updateRoleDto.Name
	role.Slug = slug.Make(updateRoleDto.Name)
	role.Desc = updateRoleDto.Desc
	database.DB.Save(&role)
	role, _ = rc.roleService.GetRole(role.Id)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "update role",
		"data":    role,
	})
}

func (rc *RoleController) DeleteRole(ctx *gin.Context) {
	paramId := ctx.Param("id")
	var getRoleDto dto.GetRoleDto
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
	role, err := rc.roleService.GetRole(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}

	database.DB.Delete(&role)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "delete role",
	})
}
