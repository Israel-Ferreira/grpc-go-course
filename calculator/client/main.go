package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/Israel-Ferreira/grpc-go-course/calculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr string = "0.0.0.0:50091"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Error on Dial gRPC Address: %v", err)
	}

	defer conn.Close()

	service := pb.NewCalculatorServiceClient(conn)

	for i := 0; i <= 10; i++ {
		result, err := service.Sum(context.Background(), &pb.SumRequest{
			Num1: int64(i - 1),
			Num2: int64(i),
		})

		if err != nil {
			log.Fatalln(err.Error())
		}

		fmt.Println(result.Result)
	}
}
