package test

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/serj162218/go_example/micro_services_example/helper"
	"github.com/serj162218/go_example/micro_services_example/model"
)

func TestGenerateToken(t *testing.T) {
	//Generate a JWT token with custom key
	user := model.User{}
	secretKey := []byte("test-key")
	result, err := helper.GenerateToken(user, secretKey)

	if err != nil {
		t.Errorf("expect no error, but got %v", err)
	}

	if result == "" {
		t.Errorf("expect result not empty")
	}
}

func TestVerifyToken(t *testing.T) {
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJJRCI6InRlc3QiLCJFbWFpbCI6InRlc3RAZXhhbXBsZS5jb20ifQ.F8jgR95R3FAs6L1Cn4-eZpLv7N2tKJ7r5VmuqBDYD6k"
	secretKey := []byte("test-key")

	result, err := helper.VerifyToken(token, secretKey)

	if err != nil {
		t.Errorf("expect no error, but got %v", err)
	}

	if !result.Valid {
		t.Errorf("expect result is valided")
	}

	claims, ok := result.Claims.(jwt.MapClaims)
	if !ok {
		t.Errorf("expect result claims to be MapClaims")
	}

	if claims["Email"] != "test@example.com" {
		t.Errorf("expect result Email to be test@example.com")
	}
	if claims["ID"] != "test" {
		t.Errorf("expect result ID to be test")
	}
}
