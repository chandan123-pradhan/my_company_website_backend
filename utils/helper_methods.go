package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// GenerateToken generates a new JWT token for a user.
//
// It takes the user ID and email as parameters and creates a token
// that includes these claims along with an expiration time of 72 hours.
//
// Returns the signed token as a string and an error if any occurs
// during the signing process.
func GenerateToken(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	return token.SignedString(jwtSecret) // Return the signed token
}

// CheckPasswordHash checks if the provided password matches the hashed password.
// 
// It takes the plain text password and the hashed password as parameters.
// Returns an error if the passwords do not match or nil if they do match.
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) // Compare the password with the hash
}

// HashPassword hashes a plain text password using bcrypt.
// 
// It takes the plain text password as a parameter and returns the hashed
// password as a string along with any error that may occur during the hashing process.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // Generate the hash
	return string(bytes), err // Return the hashed password and any error
}


// ParseToken parses the JWT token and returns the claims and any error encountered.
func ParseToken(tokenString string) (int, error) {
	
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify that the token method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64)) // convert float64 to int
		return userID, nil
	}
	return 0, fmt.Errorf("invalid token")
}
