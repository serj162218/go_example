package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/serj162218/go_example/grpc_twoway_practice/chatService"
)

type chatServer struct {
	pb.UnsafeChatServiceServer
}

var streamMap map[string]pb.ChatService_ChatServer

func (s *chatServer) Chat(stream pb.ChatService_ChatServer) error {
	// Use metadata to get user's information
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		log.Fatalf("metadata is nil")
	}

	// Get user's id
	clientID := md.Get("client-id")[0]
	log.Printf("%s enter the room", clientID)

	// Store in Map
	streamMap[clientID] = stream
Loop:
	for {
		// Receive message from stream
		chatMessage, err := stream.Recv()
		if err != nil {
			log.Printf("Failed to receive message from client: %v", err)
			break Loop
		}

		log.Printf("Received message from client: From=%s, Message=%s", chatMessage.From, chatMessage.Message)

		response := &pb.ChatMessage{
			Message: chatMessage.Message,
			From:    chatMessage.From,
		}

		for _, v := range streamMap {
			if err := v.Send(response); err != nil {
				log.Fatalf("Failed to send response to client: %v", err)
			}
		}
	}
	// Remember to delete the closed stream in map
	delete(streamMap, clientID)
	return nil
}

func main() {
	streamMap = map[string]pb.ChatService_ChatServer{}
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterChatServiceServer(server, &chatServer{})

	log.Println("Starting gRPC server on port :50051...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
