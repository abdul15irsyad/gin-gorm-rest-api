package v1

import (
	"belajar-gin/handlers"

	"github.com/gin-gonic/gin"
)

func Mahasiswa(route *gin.RouterGroup) {
	route.GET("/", handlers.GetAllMahasiswa)
	route.GET("/:id", handlers.GetMahasiswa)
}
