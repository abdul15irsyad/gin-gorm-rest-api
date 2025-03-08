package services

import (
	"errors"
	"gin-gorm-rest-api/lib"
	"gin-gorm-rest-api/models"
	"gin-gorm-rest-api/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenService struct {
	db *gorm.DB
}

func NewTokenService(libDB *lib.LibDatabase) *TokenService {
	return &TokenService{libDB.Database}
}

func (ts *TokenService) GetToken(tokenString string) (models.Token, error) {
	var token models.Token
	result := ts.db.Where("token = ? AND type = ? AND used_at IS NULL", tokenString, models.TokenForgotPassword).First(&token)
	if result.Error != nil {
		return models.Token{}, result.Error
	}

	return token, nil
}

func (ts *TokenService) CreateToken(userId uuid.UUID) (models.Token, error) {
	var randomString string
	for {
		var token models.Token
		randomString = utils.GenerateRandomString(64)
		result := ts.db.Where("token = ?", randomString).First(&token)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			break
		}
	}
	randomUuid, _ := uuid.NewRandom()
	token := models.Token{
		BaseModel: models.BaseModel{Id: randomUuid},
		Token:     randomString,
		Type:      models.TokenForgotPassword,
		UserId:    userId,
		ExpiredAt: time.Now().Add(time.Hour),
	}
	result := ts.db.Save(&token)
	if result.Error != nil {
		return models.Token{}, result.Error
	}
	token, _ = models.GetToken(ts.db, token.Id)

	return token, nil
}

func (ts *TokenService) UpdateTokenUsedAt(token models.Token) error {
	now := time.Now()
	token.UsedAt = &now
	result := ts.db.Save(&token)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
