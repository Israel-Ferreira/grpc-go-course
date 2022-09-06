package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/Israel-Ferreira/grpc-go-course/greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr string = "0.0.0.0:50051"

type GreetClientService struct {
	pb.GreetServiceClient
}

func main() {

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect in gRPC server: %v \n ", err.Error())
	}

	defer conn.Close()

	c := pb.NewGreetServiceClient(conn)

	result, err  := c.Greet(context.Background(), &pb.GreetRequest{FirstName: "Israel"})

	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println(result)

}
