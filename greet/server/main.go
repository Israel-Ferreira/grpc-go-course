package main

import (
	"log"
	"net"

	pb "github.com/Israel-Ferreira/grpc-go-course/greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var addr string = "0.0.0.0:50051"

type Server struct {
	pb.GreetServiceServer
}

func main() {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Error on Listen on Address: %v \n", err.Error())
	}

	log.Println("Server started on Address: ", addr)

	tls := true
	opts := []grpc.ServerOption{}

	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"

		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)

		if err != nil {
			log.Fatalf("Failed on Loading certificates: %v \n", err)
		}

		opts = append(opts, grpc.Creds(creds))
	}

	server := grpc.NewServer(opts...)
	pb.RegisterGreetServiceServer(server, &Server{})

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Error on start gRPC server: %v \n", err)
	}

}
