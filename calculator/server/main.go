package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Israel-Ferreira/grpc-go-course/calculator/proto"
	"google.golang.org/grpc"
)

var addr string = "0.0.0.0:50091"

type SumService struct {
	pb.CalculatorServiceServer
}

func (s *SumService) Sum(ctx context.Context, in *pb.SumRequest) (*pb.SumResponse, error) {
	a := in.Num1
	b := in.Num2

	result := a + b

	return &pb.SumResponse{Result: result}, nil
}

func main() {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed on Create listener: %v \n", err)
	}

	server := grpc.NewServer()

	pb.RegisterCalculatorServiceServer(server, &SumService{})

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Error on Create gRPC: %v \n", err)
	}
}
