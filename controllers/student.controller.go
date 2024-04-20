package controllers

import (
	"belajar-gin/database"
	"belajar-gin/dto"
	"belajar-gin/models"
	"belajar-gin/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllStudent(ctx *gin.Context) {
	var students []models.Student
	database.DB.Find(&students)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "get all student",
		"data":    students,
	})
}

func GetStudent(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var student models.Student
	result := database.DB.First(&student, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "get student",
		"data":    student,
	})
}

func CreateStudent(ctx *gin.Context) {
	var createStudentDto dto.CreateStudentDto
	ctx.ShouldBindJSON(&createStudentDto)
	errors := utils.Validate(createStudentDto)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  errors,
		})
		return
	}

	var student models.Student
	student.Name = createStudentDto.Name
	student.Year = createStudentDto.Year
	database.DB.Save(&student)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "create student",
		"data":    student,
	})
}

func UpdateStudent(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var student models.Student
	result := database.DB.First(&student, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
			"data":    nil,
		})
		return
	}

	var updatedStudent dto.UpdateStudentDto
	ctx.ShouldBindJSON(&updatedStudent)
	errors := utils.Validate(updatedStudent)
	if errors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors validation",
			"errors":  errors,
		})
		return
	}

	student.Name = updatedStudent.Name
	student.Year = updatedStudent.Year
	database.DB.Save(&student)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "update student",
		"data":    student,
	})
}

func DeleteStudent(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var student models.Student
	result := database.DB.First(&student, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
			"data":    nil,
		})
		return
	}

	database.DB.Delete(&student)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "delete student",
		"data":    nil,
	})
}
