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

func main() {
	userID := 1234
	ctx := context.WithValue(context.Background(), "userID", userID)
	val, err := fetchData(ctx, userID)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Println(val)
}

func fetchData(ctx context.Context, userID int) (string, error) {
	fmt.Println("userID:", ctx.Value("userID"))
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

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

func fetchDataFromThirdParty(userID int) (string, error) {
	time.Sleep(time.Millisecond * 500)
	return "fetch success", nil
}