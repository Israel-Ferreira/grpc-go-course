package main

import (
	"context"
	"log"

	pb "github.com/Israel-Ferreira/grpc-go-course/blog/proto"
	"github.com/Israel-Ferreira/grpc-go-course/blog/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.BlogServiceServer
}

func (s *Server) CreateBlog(ctx context.Context, blog *pb.Blog) (*pb.BlogId, error) {
	log.Println("Criando Post no Blog.....")

	post := models.NewBlogItem(blog)

	result, err := collection.InsertOne(ctx, post)

	if err != nil {
		return nil, err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)

	if !ok {
		return nil, status.Error(codes.Aborted, "Erro ao criar o POST")
	}

	blogId := &pb.BlogId{Id: oid.Hex()}

	return blogId, nil
}

func (s *Server) UpdateBlog(ctx context.Context, blog *pb.Blog) (*pb.EmptyMessage, error) {
	return nil, nil
}

func (s *Server) ReadBlog(ctx context.Context, blogId *pb.BlogId) (*pb.Blog, error) {
	postId := blogId.Id

	oid, err := primitive.ObjectIDFromHex(postId)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Cannot parse ID",
		)
	}

	var result models.BlogItem

	if err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&result); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			"Cannot found post by ID",
		)
	}

	return result.DocumentToBlog(), nil
}

func (s *Server) DeleteBlog(ctx context.Context, id *pb.BlogId) (*pb.EmptyMessage, error) {
	return nil, nil
}

func (s *Server) ListBlog(em *pb.EmptyMessage, stream pb.BlogService_ListBlogServer) error {
	return nil
}
