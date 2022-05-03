package main

import (
	"context"
	"go-proto-crud/post"
	"go-proto-crud/store"
	"log"
	"net"
	"os"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	store store.Store
}

func (s *server) ListPosts(in *emptypb.Empty, stream post.PostService_ListPostsServer) error {
	for _, post := range s.store.List() {
		stream.Send(&post)
	}
	return nil
}

func (s *server) ReadPost(c context.Context, req *post.PostId) (*post.Post, error) {
	id := req.GetId()
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id is not a type of uuid")
	}

	post, err := s.store.Read(id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return post, nil
}

func (s *server) CreatePost(c context.Context, req *post.Post) (*emptypb.Empty, error) {
	title := req.GetTitle()
	if len(title) < 5 || len(title) > 255 {
		return nil, status.Errorf(codes.InvalidArgument, "title has to be between 5 and 255")
	}

	rate := req.GetRate()
	if rate < 1 || rate > 5 {
		return nil, status.Errorf(codes.InvalidArgument, "rate has to be in range 1 to 5")
	}

	post := post.Post{
		Id:      uuid.New().String(),
		Title:   title,
		Rate:    rate,
		IsDone:  false,
		Updated: nil,
	}

	err := s.store.Create(post)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error occurred")
	}

	return &emptypb.Empty{}, nil
}

func (s *server) UpdatePost(c context.Context, req *post.Post) (*emptypb.Empty, error) {
	id := req.GetId()
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id is not a type of uuid")
	}

	title := req.GetTitle()
	if len(title) < 5 || len(title) > 255 {
		return nil, status.Errorf(codes.InvalidArgument, "title has to be between 5 and 255")
	}

	rate := req.GetRate()
	if rate < 1 || rate > 5 {
		return nil, status.Errorf(codes.InvalidArgument, "rate has to be in range 1 to 5")
	}

	post := post.Post{
		Id:     id,
		Title:  title,
		Rate:   rate,
		IsDone: req.GetIsDone(),
	}

	err = s.store.Update(post)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *server) RemovePost(c context.Context, req *post.PostId) (*emptypb.Empty, error) {
	id := req.GetId()
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id is not a type of uuid")
	}

	err = s.store.Remove(id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func main() {
	str := store.NewStore()

	connection, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	post.RegisterPostServiceServer(s, &server{
		store: str,
	})
	log.Println("Server listen on 50051")

	if err := s.Serve(connection); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
