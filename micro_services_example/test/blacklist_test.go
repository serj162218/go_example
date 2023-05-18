package test

import (
	"context"
	"testing"

	"github.com/serj162218/go_example/micro_services_example/initializer"
	"github.com/serj162218/go_example/micro_services_example/model"
)



/**
 * Delete the token from redis first then check if the result is false
 */
func TestIsTokenInBlackList(t *testing.T) {
	token := "testToken"
	initializer.RDB.SRem(context.Background(), "black_list", token)
	result := model.IsTokenInBlackList(token)
	if result != false {
		t.Errorf("expected %v, got %v", false, result)
	}
}

/**
 * Test that is the token add into black list successfully
 * Remember to delete the token from redis after the test is done
 */
func TestAddTokenToBlacklist(t *testing.T) {
	token := "test token"
	result := model.AddTokenToBlacklist(token)
	if result != nil {
		t.Errorf("expected no error, but got %v", result)
	}
	isExist, err := initializer.RDB.SIsMember(context.Background(), "black_list", token).Result()
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	if isExist != true {
		t.Errorf("expected %v, but got %v", true, isExist)
	}
	initializer.RDB.SRem(context.Background(), "black_list", token)
}
