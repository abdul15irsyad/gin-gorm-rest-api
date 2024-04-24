package middlewares

import (
	"errors"
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	payload, err := utils.ValidateJwt(accessToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "TOKEN_EXPIRED",
				"message": "token expired",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if (*payload).Type != "accessToken" {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "token is not access token",
		})
		return
	}

	// check to database
	var authUser models.User
	result := database.DB.Where("id = ?", (*payload).Id).First(&authUser)
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "invalid credential",
		})
		return
	}

	ctx.Set("authUser", authUser)
	ctx.Next()
}
