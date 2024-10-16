package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)
func GenerateToken(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	return token.SignedString(jwtSecret)
}