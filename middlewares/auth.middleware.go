package middlewares

import (
	"errors"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	jwtService  *services.JwtService
	userService *services.UserService
}

func NewAuthMiddleware(jwtService *services.JwtService, userService *services.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService,
		userService,
	}
}

func (am *AuthMiddleware) Auth(ctx *gin.Context) {
	// check authorization
	// authUser, ok := am.GetAuthUserFromAuthorization(ctx, "access")
	authUser, ok := am.GetAuthUserFromCookies(ctx, "access")
	if !ok {
		return
	}

	ctx.Set("authUser", authUser)
	ctx.Next()
}

func (am *AuthMiddleware) GetAuthUserFromAuthorization(ctx *gin.Context, tokenType string) (models.User, bool) {
	authorization := ctx.GetHeader("Authorization")
	if authorization == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid credential",
		})
		return models.User{}, false
	}
	accessToken := strings.Split(authorization, " ")[1]
	payload, err := am.jwtService.ValidateJwt(accessToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "TOKEN_EXPIRED",
				"message": "token expired",
			})
			return models.User{}, false
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid credential",
			})
			return models.User{}, false
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return models.User{}, false
	}

	if payload.Type != tokenType {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "token is not " + tokenType + " token",
		})
		return models.User{}, false
	}

	// check to database
	authUser, err := am.userService.GetUser(payload.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "invalid credential",
		})
		return models.User{}, false
	}
	return authUser, true
}

func (am *AuthMiddleware) GetAuthUserFromCookies(ctx *gin.Context, tokenType string) (models.User, bool) {
	accessToken, err := ctx.Cookie("accessToken")
	if err != nil || accessToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid credential",
		})
		return models.User{}, false
	}

	payload, err := am.jwtService.ValidateJwt(accessToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "TOKEN_EXPIRED",
				"message": "token expired",
			})
			return models.User{}, false
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid credential",
			})
			return models.User{}, false
		}
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return models.User{}, false
	}

	if payload.Type != tokenType {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "token is not " + tokenType + " token",
		})
		return models.User{}, false
	}

	// check to database
	authUser, err := am.userService.GetUser(payload.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "invalid credential",
		})
		return models.User{}, false
	}
	return authUser, true
}
