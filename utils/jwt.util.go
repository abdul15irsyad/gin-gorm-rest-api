package utils

import (
	"errors"
	"fmt"
	"gin-gorm-rest-api/database"
	"gin-gorm-rest-api/models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	Id   uuid.UUID `json:"id"`
	Type string
	jwt.RegisteredClaims
}

func GenerateJwt(id uuid.UUID, tokenType string, exp time.Duration) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		Id:   id,
		Type: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	})
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func ValidateJwt(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		jwtSecretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
		return jwtSecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}

func GetAuthUserFromAuthorization(ctx *gin.Context, tokenType string) (user *models.User, ok bool) {
	authorization := ctx.GetHeader("Authorization")
	if authorization == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid credential",
		})
		return nil, false
	}
	accessToken := strings.Split(authorization, " ")[1]
	payload, err := ValidateJwt(accessToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "TOKEN_EXPIRED",
				"message": "token expired",
			})
			return nil, false
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid credential",
			})
			return nil, false
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return nil, false
	}

	if (*payload).Type != tokenType {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "token is not " + tokenType,
		})
		return nil, false
	}

	// check to database
	var authUser models.User
	result := database.DB.Where("id = ?", (*payload).Id).First(&authUser)
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "invalid credential",
		})
		return nil, false
	}
	return &authUser, true
}

func SigningToken(ctx *gin.Context, accessToken *string, refreshToken *string, authUser *models.User) (ok bool) {
	// signing jwt
	jwtAccessTokenExpiredInHour, _ := strconv.ParseFloat(os.Getenv("JWT_ACCESS_TOKEN_EXPIRED_IN_HOUR"), 64)
	token, err := GenerateJwt((*authUser).Id, "access", time.Duration(jwtAccessTokenExpiredInHour*60*60)*time.Second)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return false
	}
	*accessToken = *token
	jwtRefreshTokenExpiredInHour, _ := strconv.ParseFloat(os.Getenv("JWT_REFRESH_TOKEN_EXPIRED_IN_HOUR"), 64)
	token, err = GenerateJwt((*authUser).Id, "refresh", time.Duration(jwtRefreshTokenExpiredInHour*60*60)*time.Second)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return false
	}
	*refreshToken = *token
	return true
}
