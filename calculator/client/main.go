package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	doDecompositePrime(service)

	doGetAvg(service)

	doMax(service)
}

func doGetAvg(c pb.CalculatorServiceClient) {
	log.Println("Invoke Avg Function")

	numbers := []int64{1, 2, 10, 9, 8, 7, 6}

	stream, err := c.Avg(context.Background())

	if err != nil {
		log.Fatalln(err)
	}

	for _, num := range numbers {
		stream.Send(&pb.AvgRequest{Num: num})
		time.Sleep(500 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.Result)

}

func doDecompositePrime(c pb.CalculatorServiceClient) {
	req := &pb.PrimeDecompositionRequest{
		Num: 13,
	}

	stream, err := c.DecompositePrime(context.Background(), req)

	if err != nil {
		log.Fatalln(err.Error())
	}

	var factors []int64

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln(err.Error())
		}

		factors = append(factors, msg.Factor)
	}

	if len(factors) == 1 && factors[0] == req.Num {
		fmt.Printf("%d is prime \n", req.Num)
	} else {
		fmt.Printf("%d is not a prime, factors : %v", req.Num, factors)
	}

}

func doMax(c pb.CalculatorServiceClient) {
	log.Println("Invoking Max function...")

	stream, err := c.Max(context.Background())

	if err != nil {
		log.Fatalln(err)
	}

	waitch := make(chan struct{})

	reqs := []*pb.MaxMsgRequest{
		{Num: 1},
		{Num: 5},
		{Num: 3},
		{Num: 6},
		{Num: 2},
		{Num: 20},
	}

	go func() {
		for _, num := range reqs {
			log.Println("Send Request to GRPC")
			stream.Send(num)

			time.Sleep(1 * time.Second)
		}

		stream.CloseSend()
	}()

	go func() {
		for {
			rec, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Println(err.Error())
				break
			}

			fmt.Printf("Received Num: %d \n", rec.Result)
		}

		close(waitch)
	}()

	<-waitch

}
