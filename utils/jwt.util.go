package utils

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJwt(id uuid.UUID, isRefreshToken bool) (*string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":             id,
		"isRefreshToken": isRefreshToken,
		"jti":            uuid.String(),
	})
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}
