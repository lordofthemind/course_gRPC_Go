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
	cc, err := grpc.NewClient("0.0.0.0:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	client := blogpb.NewBlogServiceClient(cc)
	fmt.Println("Created client!")

	// Create blog
	blogID, err := createBlog(client)
	if err != nil {
		log.Fatalf("Error while creating blog: %v", err)
	}

	// Read blog
	err = readBlog(client, blogID)
	if err != nil {
		log.Fatalf("Error while reading blog: %v", err)
	}

	// Update blog
	err = updateBlog(client, blogID)
	if err != nil {
		log.Fatalf("Error while updating blog: %v", err)
	}

	// Delete blog
	err = deleteBlog(client, blogID)
	if err != nil {
		log.Fatalf("Error while deleting blog: %v", err)
	}

}

func createBlog(client blogpb.BlogServiceClient) (string, error) {
	fmt.Println("Creating the blog")

	blog := &blogpb.Blog{
		AuthorId: "Stephen King",
		Title:    "The Shining",
		Content:  "Here's Johnny!",
	}

	res, err := client.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		return "", fmt.Errorf("unexpected error: %v", err)
	}

	fmt.Printf("Blog has been created: %v\n", res)
	return res.GetBlog().GetId(), nil
}

func readBlog(client blogpb.BlogServiceClient, blogID string) error {
	fmt.Println("Reading the blog")

	// Trying to read a blog with an invalid ID to test error handling
	_, err := client.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "invalidID"})
	if err != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	}

	// Reading the blog with the correct ID
	res, err := client.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: blogID})
	if err != nil {
		return fmt.Errorf("error happened while reading: %v", err)
	}

	fmt.Printf("Blog was read: %v\n", res)
	return nil
}

func updateBlog(client blogpb.BlogServiceClient, blogID string) error {
	fmt.Println("Updating the blog")

	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Stephen King",
		Title:    "The Shining",
		Content:  "Here's Johnny! Redrum!",
	}

	res, err := client.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if err != nil {
		return fmt.Errorf("error happened while updating: %v", err)
	}

	fmt.Printf("Blog was updated: %v\n", res)
	return nil

}

func deleteBlog(client blogpb.BlogServiceClient, blogID string) error {
	fmt.Println("Deleting the blog")

	_, err := client.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if err != nil {
		return fmt.Errorf("error happened while deleting: %v", err)
	}

	fmt.Println("Blog was deleted")
	return nil
}
