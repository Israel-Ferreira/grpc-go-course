package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/Israel-Ferreira/grpc-go-course/greet/proto"
)

func (c *Server) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {

	log.Printf("Requested Name: %s \n", in.FirstName)
	name := in.FirstName

	return &pb.GreetResponse{Result: fmt.Sprintf("Hello, %s", name)}, nil
}
