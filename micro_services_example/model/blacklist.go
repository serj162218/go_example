package model

import (
	"context"
	"log"

	"github.com/serj162218/go_example/micro_services_example/initializer"
)

type BlackList struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

func AddTokenToBlacklist(token string) error {
	// store this token to the database
	err := initializer.RDB.SAdd(context.TODO(), "black_list", token).Err()
	if err != nil {
		return err
	}
	return nil
}

func IsTokenInBlackList(token string) bool {
	//check if the token is in redis
	isExist, err := initializer.RDB.SIsMember(context.TODO(), "black_list", token).Result()
	if err != nil {
		log.Fatal(err.Error())
		return true
	}
	return isExist
}
