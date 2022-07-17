package routers

import (
	v1 "belajar-gin/routers/v1"

	"github.com/gin-gonic/gin"
)

func V1(router *gin.Engine) {
	v1Router := router.Group("/v1")
	v1.Root(v1Router.Group("/"))
	v1.Mahasiswa(v1Router.Group("/mahasiswa"))
}
