package main

import (
	"context"
	"fmt"
	"io"
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

func (c *Server) LongGreet(stream pb.GreetService_LongGreetServer) error {
	log.Println("Long Greet Function was Invoked!!")

	res := ""

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&pb.GreetResponse{Result: res})
		}

		if err != nil {
			return err
		}

		fmt.Println(req)

		res += fmt.Sprintf("Hello %s \n", req.FirstName)

	}

}

func (c *Server) GreetEveryone(stream pb.GreetService_GreetEveryoneServer) error {
	log.Println("GreetEveryone was invoked!!!")

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			fmt.Printf("Error on While Reading client stream: %s \n", err.Error())
			return err
		}

		res := fmt.Sprintf("Hello %s !!", req.FirstName)

		if err = stream.Send(&pb.GreetResponse{Result: res}); err != nil {
			return err
		}
	}

}
