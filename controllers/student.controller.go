package controllers

import (
	"belajar-gin/database"
	"belajar-gin/models"
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
		"msg":  "get all student",
		"data": students,
	})
}

func GetStudent(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var student models.Student
	result := database.DB.First(&student, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg":  "data not found",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "get student",
		"data": student,
	})
}

func CreateStudent(ctx *gin.Context) {
	var student models.Student
	ctx.ShouldBindJSON(&student)

	database.DB.Save(&student)
	ctx.JSON(http.StatusCreated, gin.H{
		"msg":  "create student",
		"data": student,
	})
}

func UpdateStudent(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var student models.Student
	result := database.DB.First(&student, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg":  "data not found",
			"data": nil,
		})
		return
	}

	var updatedStudent models.Student
	ctx.ShouldBindJSON(&updatedStudent)
	student.Name = updatedStudent.Name
	student.Year = updatedStudent.Year

	database.DB.Save(&student)

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "update student",
		"data": student,
	})
}

func DeleteStudent(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var student models.Student
	result := database.DB.First(&student, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg":  "data not found",
			"data": nil,
		})
		return
	}

	database.DB.Delete(&student)

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "delete student",
		"data": nil,
	})
}
