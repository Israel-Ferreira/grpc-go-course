package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/Israel-Ferreira/grpc-go-course/blog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := "localhost:9090"

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	service := pb.NewBlogServiceClient(conn)

	createPost(service, &pb.Blog{
		AuthorId: "Israel",
		Title:    "Teste",
		Content:  "Criando um post qualquer",
	})
}

func createPost(service pb.BlogServiceClient, blog *pb.Blog) {
	res, err := service.CreateBlog(context.Background(), blog)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.Id)
}
