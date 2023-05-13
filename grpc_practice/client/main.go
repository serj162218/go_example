package main

import (
	"context"
	"log"

	pb "github.com/serj162218/go_example/grpc_practice/orderService"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createOrder(client pb.OrderServiceClient, order *pb.Order) {
	resp, err := client.CreateOrder(context.Background(), &pb.CreateOrderRequest{Order: order})
	if err != nil {
		log.Fatalf("Error occurred while calling CreateOrder: %v", err)
	}

	// receive response
	log.Printf("response: %v", resp)
}

func getOrder(client pb.OrderServiceClient, orderId int64) {
	resp, err := client.GetOrder(context.Background(), &pb.GetOrderRequest{OrderId: 123})
	if err != nil {
		log.Fatalf("Error occurred while calling GetOrder: %v", err)
	}

	// receive response
	log.Printf("response: %v", resp)
}

func cancelOrder(client pb.OrderServiceClient, orderId int64) {
	resp, err := client.CancelOrder(context.Background(), &pb.CancelOrderRequest{OrderId: 123})
	if err != nil {
		log.Fatalf("Error occurred while calling CancelOrder: %v", err)
	}

	// receive response
	log.Printf("response: %v", resp)
}

func getAllOrders(client pb.OrderServiceClient) {
	resp, err := client.GetAllOrders(context.Background(), &pb.GetAllOrdersRequest{})
	if err != nil {
		log.Fatalf("Error occurred while calling GetAllOrders: %v", err)
	}

	// receive response
	log.Printf("response: %v", resp)
}

func main() {
	// Connect to grpc server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Cannot connect to gRPC server side: %v", err)
	}
	defer conn.Close()

	// new a gRPC client
	client := pb.NewOrderServiceClient(conn)

	//call CreateOrder
	createOrder(client, &pb.Order{
		Id:         123,
		CreateTime: 10000,
		Status:     pb.OrderStatus_CREATED,
		Items: []*pb.OrderItem{
			{
				ProductId:   1,
				ProductName: "fish",
				Quantity:    50,
			}, {
				ProductId:   2,
				ProductName: "cookies",
				Quantity:    30,
			},
		},
	})

	//call CreateOrder
	createOrder(client, &pb.Order{
		Id:         456,
		CreateTime: 20000,
		Status:     pb.OrderStatus_CREATED,
		Items: []*pb.OrderItem{
			{
				ProductId:   3,
				ProductName: "cakes",
				Quantity:    30,
			}, {
				ProductId:   4,
				ProductName: "juices",
				Quantity:    10,
			},
		},
	})

	//call GetOrder
	getOrder(client, 123)

	//call CancelOrder
	cancelOrder(client, 123)

	//call GetAllOrders
	getAllOrders(client)
}
