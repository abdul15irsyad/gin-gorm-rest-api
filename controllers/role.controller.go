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
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func GetAllRoles(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	search := ctx.Query("search")
	roles, total, totalPage, err := models.GetPaginatedRoles(database.DB, page, limit, &search)
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

func GetRole(ctx *gin.Context) {
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
	role, err := models.GetRole(database.DB, id)
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

func CreateRole(ctx *gin.Context) {
	var createRoleDto dto.CreateRoleDto
	ctx.ShouldBind(&createRoleDto)
	validationErrors := utils.Validate(createRoleDto)
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
	role, _ = models.GetRole(database.DB, role.Id)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "create role",
		"data":    role,
	})
}

func UpdateRole(ctx *gin.Context) {
	paramId := ctx.Param("id")
	var updateRoleDto dto.UpdateRoleDto
	ctx.ShouldBind(&updateRoleDto)
	updateRoleDto.Id = paramId
	validationErrors := utils.Validate(updateRoleDto)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  validationErrors,
		})
		return
	}

	// save to database
	id, _ := uuid.Parse(paramId)
	role, err := models.GetRole(database.DB, id)
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
	role, _ = models.GetRole(database.DB, role.Id)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "update role",
		"data":    role,
	})
}

func DeleteRole(ctx *gin.Context) {
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
	role, err := models.GetRole(database.DB, id)
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
