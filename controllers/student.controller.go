package controllers

import (
	"belajar-gin/configs"
	"belajar-gin/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllStudent(c *gin.Context) {
	var students []models.Student
	configs.DB.Find(&students)
	c.JSON(http.StatusOK, gin.H{
		"msg":  "get all student",
		"data": students,
	})
}

func GetStudent(c *gin.Context) {
	// convert string to int
	id, _ := strconv.Atoi(c.Param("id"))
	var student models.Student
	configs.DB.Where("id = ?",id).First(&student)

	c.JSON(http.StatusOK, gin.H{
		"msg":  "get student",
		"data": student,
	})
}
