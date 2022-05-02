package main

import (
	"context"
	"log"
	"net"
	"os"

	"go-proto-crud/blog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct{}

func (*server) ListBlogs(in *emptypb.Empty, stream blog.BlogService_ListBlogsServer) error {

	blog := &blog.Blog{
		Id:      "132",
		Title:   "Blog",
		Rate:    3,
		IsDone:  false,
		Updated: &timestamppb.Timestamp{},
	}

	stream.Send(blog)
	stream.Send(blog)
	stream.Send(blog)

	return nil
}

func (*server) ReadBlog(c context.Context, req *blog.BlogId) (*blog.Blog, error) {
	id := req.GetId()

	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Received empty string for id")
	}

	res := &blog.Blog{
		Id: id,
	}

	return res, nil
}

func (*server) CreateBlog(c context.Context, req *blog.Blog) (*emptypb.Empty, error) {
	return nil, nil
}

func (*server) UpdateBlog(c context.Context, req *blog.Blog) (*emptypb.Empty, error) {
	return nil, nil
}

func (*server) RemoveBlog(c context.Context, req *blog.BlogId) (*emptypb.Empty, error) {
	return nil, nil
}

func main() {
	connection, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	blog.RegisterBlogServiceServer(s, &server{})
	log.Println("Server listen on 50051")

	if err := s.Serve(connection); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
