package services

import (
	"fmt"
	"gin-gorm-rest-api/models"
	"net/http"
	"os"
	"strconv"
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

type JwtService struct {
}

func NewJwtService() *JwtService {
	return &JwtService{}
}

func (js *JwtService) GenerateJwt(id uuid.UUID, tokenType string, exp time.Duration) (string, error) {
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
		return "", err
	}
	return tokenString, nil
}

func (js *JwtService) ValidateJwt(tokenString string) (CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		jwtSecretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
		return jwtSecretKey, nil
	})
	if err != nil {
		return CustomClaims{}, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return CustomClaims{}, err
	}

	return *claims, nil
}

func (js *JwtService) SigningToken(ctx *gin.Context, authUser models.User) (accessToken string, refreshToken string, ok bool) {
	// signing jwt
	jwtAccessTokenExpiredInHour, _ := strconv.ParseFloat(os.Getenv("JWT_ACCESS_TOKEN_EXPIRED_IN_HOUR"), 64)
	accessToken, err := js.GenerateJwt(authUser.Id, "access", time.Duration(jwtAccessTokenExpiredInHour*60*60)*time.Second)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return "", "", false
	}
	jwtRefreshTokenExpiredInHour, _ := strconv.ParseFloat(os.Getenv("JWT_REFRESH_TOKEN_EXPIRED_IN_HOUR"), 64)
	refreshToken, err = js.GenerateJwt(authUser.Id, "refresh", time.Duration(jwtRefreshTokenExpiredInHour*60*60)*time.Second)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return "", "", false
	}
	return accessToken, refreshToken, true
}
