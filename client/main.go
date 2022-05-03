package main

import (
	"context"
	"fmt"
	"go-proto-crud/post"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type client struct {
	client post.PostServiceClient
}

func (c *client) ListPosts() ([]post.Post, error) {
	stream, err := c.client.ListPosts(context.Background(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	posts := []post.Post{}
	for {
		post, err := stream.Recv()
		if err == io.EOF {
			break // end of the stream
		}
		if err != nil {
			return nil, err
		}
		posts = append(posts, *post)
	}
	return posts, nil
}

func (c *client) ReadPost(id string) (*post.Post, error) {
	post, err := c.client.ReadPost(context.Background(), &post.PostId{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return post, nil
}

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}
	defer connection.Close()

	client := client{
		client: post.NewPostServiceClient(connection),
	}

	posts, err := client.ListPosts()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	post, err := client.ReadPost(posts[0].Id)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("ListPosts: ", posts)
	fmt.Println("ReadPost: ", post)
}
