package main

import (
	"context"
	"fmt"
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
	log.Println("Update Function was called")

	fmt.Println(blog.Id)
	oid, err := primitive.ObjectIDFromHex(blog.Id)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Cannot parse id to objectid",
		)
	}

	post := models.NewBlogItem(blog)

	res, err := collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": post})

	if err != nil {
		fmt.Println(err.Error())
		return nil, status.Errorf(codes.NotFound, "Id não encontrado")
	}

	fmt.Println(res)

	return &pb.EmptyMessage{}, nil
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
	log.Println("DeleteBlog was invoked!!...")

	oid, err := primitive.ObjectIDFromHex(id.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Error on Parse ID: %v", err))
	}

	if _, err := collection.DeleteOne(context.Background(), bson.M{"_id": oid}); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Error: blog post not found, err: %v", err))
	}

	return &pb.EmptyMessage{}, nil
}

func (s *Server) ListBlog(em *pb.EmptyMessage, stream pb.BlogService_ListBlogServer) error {
	log.Println("List blog was invoked")

	cursor, err := collection.Find(context.Background(), primitive.D{{}})

	if err != nil {
		return status.Errorf(codes.Internal, "Error on consulting db: %v", err)
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		data := &models.BlogItem{}

		if err := cursor.Decode(&data); err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Error on decode data: %s", err))
		}

		stream.Send(data.DocumentToBlog())
	}

	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Error in Cursor: %s", err.Error()))
	}

	return nil
}
