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

func (c *Server) GreetManyTimes(in *pb.GreetRequest, stream pb.GreetService_GreetManyTimesServer) error {
	log.Printf("GreetManyTimes function was invoked with: %v \n", in)

	for i := 0; i < 10; i++ {
		res := fmt.Sprintf("Hello %s, %d \n", in.FirstName, i)
		err := stream.Send(&pb.GreetResponse{Result: res})

		if err != nil {
			return err
		}
	}

	return nil
}
