package handlers

import (
	"belajar-gin/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllMahasiswa(c *gin.Context) {
	ms := structs.Mahasiswas{}.LoadData()

	c.JSON(http.StatusOK, gin.H{
		"msg":  "get all mahasiswa",
		"data": ms.Mahasiswas,
	})
}

func GetMahasiswa(c *gin.Context) {
	// convert string to int
	id, _ := strconv.Atoi(c.Param("id"))

	ms := structs.Mahasiswas{}.LoadData()
	mahasiswa := ms.FindMahasiswa(func(m structs.Mahasiswa) bool {
		return m.Id == id
	})

	c.JSON(http.StatusOK, gin.H{
		"msg":  "get mahasiswa",
		"data": mahasiswa,
	})
}
