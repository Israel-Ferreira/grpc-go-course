package main

import (
	"log"
	"net"

	pb "github.com/Israel-Ferreira/grpc-go-course/greet/proto"
	"google.golang.org/grpc"
)

var addr string = "0.0.0.0:5051"

type Server struct {
	pb.GreetServiceServer
}

func main() {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Error on Listen on Address: %v \n", err.Error())
	}

	log.Println("Server started on Address: ", addr)

	server := grpc.NewServer()

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Error on start gRPC server: %v \n", err)
	}

}
