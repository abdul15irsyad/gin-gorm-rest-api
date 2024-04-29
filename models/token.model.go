package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Token struct {
	BaseModel
	Token     string    `json:"token" gorm:"index,not null"`
	Type      TokenType `json:"type" gorm:"not null"`
	ExpiredAt time.Time `json:"expiredAt"`
	UserId    uuid.UUID `json:"-" gorm:"not null"`
	User      *User     `json:"user" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type TokenType string

const (
	TokenForgotPassword TokenType = "forgot-password"
	TokenVerifyEmail    TokenType = "verify-email"
)

func (token *Token) AfterLoad() {
	if token != nil {
		token.User.AfterLoad()
	}
}

func GetToken(db *gorm.DB, id uuid.UUID) (Token, error) {
	var token Token
	result := db.Preload(clause.Associations).First(&token, "id = ?", id)
	if result.Error != nil {
		return Token{}, result.Error
	}
	token.AfterLoad()
	return token, nil
}
