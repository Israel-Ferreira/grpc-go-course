package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	pb "github.com/Israel-Ferreira/grpc-go-course/calculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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

func (s *SumService) DecompositePrime(in *pb.PrimeDecompositionRequest, stream pb.CalculatorService_DecompositePrimeServer) error {
	k := 2

	n := int(in.Num)

	var factors []int64

	for n > 1 {
		if n%k == 0 {
			fmt.Println(k)
			n = n / k

			factors = append(factors, int64(k))
		} else {
			k += 1
		}
	}

	for _, factor := range factors {
		res := &pb.PrimeDecompositionResponse{Factor: factor}

		if err := stream.Send(res); err != nil {
			log.Fatalf("Error on Sending response to client: %v \n", err)
		}
	}

	return nil
}

func (s *SumService) Avg(stream pb.CalculatorService_AvgServer) error {
	var numbers []int64

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			sumVal := sum(numbers)

			fmt.Println(sumVal)
			total := float64(sumVal / float64(len(numbers)))

			return stream.SendAndClose(&pb.AvgResponse{Result: total})
		}

		if err != nil {
			return err
		}

		fmt.Printf("Número Enviado pela Stream: %d \n", req.Num)

		numbers = append(numbers, req.Num)
	}
}

func (s *SumService) Max(stream pb.CalculatorService_MaxServer) error {
	var lastMaxNum int64

	for {
		rec, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		num := rec.Num

		if num > lastMaxNum {
			stream.Send(&pb.MaxMsgResponse{Result: num})
			lastMaxNum = num
		}
	}

}

func (s *SumService) Sqrt(ctx context.Context, in *pb.SqrtRequest) (*pb.SqrtResponse, error) {
	log.Printf("Sqrt Function Was Invoked: %v \n", in)

	num := in.Num

	if num < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number %d \n", num),
		)
	}

	return &pb.SqrtResponse{
		Result: math.Sqrt(float64(num)),
	}, nil
}

func sum(numbers []int64) float64 {
	var sumResult float64

	for _, num := range numbers {
		sumResult += float64(num)
	}

	return sumResult
}

func main() {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed on Create listener: %v \n", err)
	}

	server := grpc.NewServer()

	pb.RegisterCalculatorServiceServer(server, &SumService{})

	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Error on Create gRPC: %v \n", err)
	}
}
