package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lordofthemind/course_gRPC_Go/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	// Create a connection to the server
	opts := grpc.WithInsecure()
	cc, err := grpc.NewClient("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)
	fmt.Println("Created client!")

	// Create blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Stephen King",
		Title:    "The Shining",
		Content:  "Here's Johnny!",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v", res)
}
