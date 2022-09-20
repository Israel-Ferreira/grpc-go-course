package models

import (
	pb "github.com/Israel-Ferreira/grpc-go-course/blog/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorId string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}

func (blg *BlogItem) DocumentToBlog() *pb.Blog {
	return &pb.Blog{
		Id:       blg.ID.Hex(),
		AuthorId: blg.AuthorId,
		Title:    blg.Title,
		Content:  blg.Content,
	}
}

func NewBlogItem(protoBlog *pb.Blog) *BlogItem {
	return &BlogItem{
		AuthorId: protoBlog.AuthorId,
		Title:    protoBlog.Title,
		Content:  protoBlog.Content,
	}
}
