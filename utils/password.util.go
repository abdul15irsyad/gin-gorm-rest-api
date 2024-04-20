package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	cost := bcrypt.DefaultCost // Adjust cost as needed (higher = more secure, slower)
	return bcrypt.GenerateFromPassword([]byte(password), cost)
}
