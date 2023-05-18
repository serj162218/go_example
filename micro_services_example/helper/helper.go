package helper

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/serj162218/go_example/micro_services_example/model"
)

func GenerateToken(user model.User, secretKey []byte) (string, error) {
	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string, secretKey []byte) (*jwt.Token, error) {
	// Verify JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check that the JWT token is valid
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// return JWT token
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
