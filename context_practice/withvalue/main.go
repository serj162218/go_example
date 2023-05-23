package main

import (
	"context"
	"fmt"
	"time"
)

type Response struct {
	value string
	err   error
}

type User struct {
	id string
}

func main() {
	userID := "1234"
	user := User{
		id: userID,
	}

	ctx := context.WithValue(context.Background(), user, userID)

	val, err := fetchData(ctx, user)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Println(val)
}

func fetchData(ctx context.Context, user User) (string, error) {
	var userID string = ctx.Value(user).(string)
	fmt.Println("userID:", userID)
	respch := make(chan Response)
	go func() {
		val, err := fetchDataFromThirdParty(userID)
		respch <- Response{val, err}
	}()

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("%v", ctx.Err())
		case resp := <-respch:
			return resp.value, nil

		}
	}
}

func fetchDataFromThirdParty(userID string) (string, error) {
	time.Sleep(time.Millisecond * 500)
	return "fetch success", nil
}
