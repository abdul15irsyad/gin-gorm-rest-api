package middlewares

import (
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {
	// check authorization
	authorization := ctx.GetHeader("Authorization")
	if authorization == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid credential",
		})
		return
	}
	accessToken := strings.Split(authorization, " ")[1]
	var payload utils.CustomClaims
	err := utils.ValidateJwt(&payload, accessToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// check to database
	var authUser models.User
	result := database.DB.Where("id = ?", payload.Id).First(&authUser)
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "invalid credential",
		})
		return
	}

	ctx.Set("authUser", authUser)
	ctx.Next()
}
