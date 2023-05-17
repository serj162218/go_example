package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/serj162218/go_example/grpc_twoway_practice/chatService"
)

func main() {
	// Set user's name
	fmt.Print("user name:")
	var userID string
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() {
		userID = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error while reading name", err)
	}

	// Connect to server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	// Set metadata
	md := metadata.New(map[string]string{"client-id": userID})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Create stream with metadata
	stream, err := client.Chat(ctx)
	if err != nil {
		log.Fatalf("Failed to open stream: %v", err)
	}

	fmt.Print("connecting to chat service...")

	// receive message
	go func() {
		for {
			response, err := stream.Recv()
			if err != nil {
				log.Fatalf("Failed to receive response from server: %v", err)
			}

			log.Printf("Received response from server: From=%s, Message=%s", response.From, response.Message)
		}
	}()

	// send message to server
	for {
		// reading input
		fmt.Print("Input:")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input := scanner.Text()

			// send message
			message := &pb.ChatMessage{
				From:    userID,
				Message: input,
			}
			if err := stream.Send(message); err != nil {
				log.Fatalf("Failed to send message to server: %v", err)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error while reading input", err)
		}
	}
}
