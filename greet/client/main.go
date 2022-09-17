package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	pb "github.com/Israel-Ferreira/grpc-go-course/greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr string = "0.0.0.0:50051"

type GreetClientService struct {
	pb.GreetServiceClient
}

func doGreetManyTimes(c pb.GreetServiceClient) {
	log.Printf("Make Request to gRPC server....")

	names := []string{
		"Israel",
		"Nickolas",
		"Milena",
		"Ciro",
		"Gabriel",
		"Amanda",
		"Marina",
	}

	randNum := rand.Intn(len(names))

	fmt.Println(randNum)

	req := &pb.GreetRequest{
		FirstName: names[randNum],
	}

	stream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalln(err.Error())
	}

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln(err.Error())
		}

		log.Printf("Received Message: %v \n", msg)
	}

}

func main() {

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect in gRPC server: %v \n ", err.Error())
	}

	defer conn.Close()

	c := pb.NewGreetServiceClient(conn)

	result, err := c.Greet(context.Background(), &pb.GreetRequest{FirstName: "Israel"})

	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println(result)

	doGreetManyTimes(c)

	doLongGreet(c)
}

func doLongGreet(c pb.GreetServiceClient) {
	reqArr := []*pb.GreetRequest{
		{FirstName: "Israel"},
		{FirstName: "Marina"},
		{FirstName: "Matheus"},
		{FirstName: "Mariana"},
		{FirstName: "Denilson"},
	}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalln(err)
	}

	for _, person := range reqArr {
		log.Println("Sending stream to long greet")

		stream.Send(person)
		time.Sleep(time.Second * 1)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.Result)
}
