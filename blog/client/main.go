package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/Israel-Ferreira/grpc-go-course/blog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	addr := "localhost:9090"

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	service := pb.NewBlogServiceClient(conn)

	id := createPost(service, &pb.Blog{
		AuthorId: "Israel",
		Title:    "Teste",
		Content:  "Criando um post qualquer",
	})

	blogId := &pb.BlogId{Id: id}

	fmt.Println(blogId)

	res := findById(service, &pb.BlogId{Id: blogId.Id})
	fmt.Println(res)

}

func createPost(service pb.BlogServiceClient, blog *pb.Blog) string {
	res, err := service.CreateBlog(context.Background(), blog)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.Id)

	return res.Id
}

func findById(client pb.BlogServiceClient, blogId *pb.BlogId) *pb.Blog {
	res, err := client.ReadBlog(context.Background(), blogId)

	if err != nil {
		status, ok := status.FromError(err)

		if ok {
			log.Fatalf("Status: %s, Error: %v \n", status.Code(), err)
		} else {
			log.Fatalln(err)
		}
	}

	return res
}
