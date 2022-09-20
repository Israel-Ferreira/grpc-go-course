package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Israel-Ferreira/grpc-go-course/blog/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var addr string = "0.0.0.0:9090"

var collection *mongo.Collection

func init() {
	connectionStr := "mongodb://root:root@localhost:27017/"

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionStr))

	if err != nil {
		log.Fatalln(err)
	}

	err = client.Connect(context.Background())

	if err != nil {
		log.Fatalln(err)
	}

	collection = client.Database("blogdb").Collection("blog")
}

func main() {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalln(err.Error())
	}

	server := grpc.NewServer()

	pb.RegisterBlogServiceServer(server, &Server{})

	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Error on start gRPC server: %s", err.Error())
	}
}
