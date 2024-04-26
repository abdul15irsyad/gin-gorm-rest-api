package middlewares

import (
	"gin-gorm-rest-api/utils"

	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {
	// check authorization
	authUser, ok := utils.GetAuthUserFromAuthorization(ctx, "access")
	if !ok {
		return
	}

	ctx.Set("authUser", authUser)
	ctx.Next()
}
