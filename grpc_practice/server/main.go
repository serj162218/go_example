package main

import (
	"context"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/serj162218/go_example/grpc_practice/orderService"
)

type server struct {
	orderMap map[int64]*pb.Order
	mu       sync.Mutex
	pb.UnsafeOrderServiceServer
}

func (s *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	fmt.Printf("CreateOrder: %+v\n", req.Order)

	// Lock on
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a new order and add it to the map
	orderID := req.Order.Id
	order := req.Order
	s.orderMap[orderID] = order

	// Return success response
	return &pb.CreateOrderResponse{
		OrderId: orderID,
	}, nil
}

func (s *server) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	fmt.Printf("CancelOrder: %+v\n", req.OrderId)
	s.mu.Lock()
	defer s.mu.Unlock()

	// Search for the order in the map
	orderID := req.OrderId
	order, exists := s.orderMap[orderID]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "Order not found: %s", orderID)
	}

	// Change the status of the order
	order.Status = pb.OrderStatus_CANCELLED
	return &pb.CancelOrderResponse{Success: true}, nil
}

func (s *server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	fmt.Printf("GetOrder: %+v\n", req.OrderId)
	// Lock on
	s.mu.Lock()
	defer s.mu.Unlock()

	// Search for the order in the map
	orderID := req.OrderId
	order, exists := s.orderMap[orderID]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "Order not found: %s", orderID)
	}

	// Return success response
	return &pb.GetOrderResponse{
		Order: order,
	}, nil
}

func (s *server) GetAllOrders(ctx context.Context, _ *pb.GetAllOrdersRequest) (*pb.GetAllOrdersResponse, error) {
	fmt.Print("GetAllOrders")
	// Lock on
	s.mu.Lock()
	defer s.mu.Unlock()

	// Loop through the map and save all orders in orderList
	orderList := make([]*pb.Order, 0)
	for _, v := range s.orderMap {
		orderList = append(orderList, v)
	}

	// Return success response
	return &pb.GetAllOrdersResponse{
		Order: orderList,
	}, nil
}
func main() {
	// New a gRPC server
	s := grpc.NewServer()

	// RegisterServiceServer
	pb.RegisterOrderServiceServer(s, &server{
		orderMap: make(map[int64]*pb.Order),
	})

	// Listen for requests
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	fmt.Println("server listening at :50051")

	// Active Server
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
