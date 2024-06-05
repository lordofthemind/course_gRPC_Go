package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/lordofthemind/course_gRPC_Go/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello from client")

	// Establishing connection to the server
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer cc.Close()

	// Creating a client for the GreetService
	c := greetpb.NewGreetServiceClient(cc)

	// Perform Unary RPC
	// doUnary(c)

	// Perform Server Streaming RPC
	// doServerStreaming(c)

	// Perform Client Streaming RPC
	// doClientStreaming(c)

	// Perform BiDi Streaming RPC
	// doBiDiStreaming(c)

	// Perform Unary RPC with deadline
	doUnaryWithDeadline(c, 5*time.Second)
	doUnaryWithDeadline(c, 1*time.Second)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Lord",
			LastName:  "OfTheMind",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Lord",
			LastName:  "OfTheMind",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// We've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	// Creating multiple requests to be sent in a stream
	requests := []*greetpb.LongGreetRequest{
		{Greeting: &greetpb.Greeting{FirstName: "Alice", LastName: "Smith"}},
		{Greeting: &greetpb.Greeting{FirstName: "Bob", LastName: "Johnson"}},
		{Greeting: &greetpb.Greeting{FirstName: "Carol", LastName: "Williams"}},
		{Greeting: &greetpb.Greeting{FirstName: "David", LastName: "Brown"}},
		{Greeting: &greetpb.Greeting{FirstName: "Eve", LastName: "Davis"}},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet RPC: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(1 * time.Second) // Simulating some delay
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet response: %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	// Creating a stream
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GreetEveryone RPC: %v", err)
		return
	}

	// Creating multiple requests to be sent in a stream
	requests := []*greetpb.GreetEveryoneRequest{
		{Greeting: &greetpb.Greeting{FirstName: "Alice", LastName: "Smith"}},
		{Greeting: &greetpb.Greeting{FirstName: "Bob", LastName: "Johnson"}},
		{Greeting: &greetpb.Greeting{FirstName: "Carol", LastName: "Williams"}},
		{Greeting: &greetpb.Greeting{FirstName: "David", LastName: "Brown"}},
		{Greeting: &greetpb.Greeting{FirstName: "Eve", LastName: "Davis"}},
	}

	waitc := make(chan struct{})

	// Sending a bunch of messages to the server
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending request: %v\n", req)
			stream.Send(req)
			time.Sleep(1 * time.Second) // Simulating some delay
		}
		stream.CloseSend()
	}()

	// Receiving a bunch of messages from the server
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving response from GreetEveryone: %v", err)
				close(waitc)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
	}()

	<-waitc
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a Unary with deadline RPC...")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Lord",
			LastName:  "OfTheMind",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Fatalf("Timeout was hit! Deadline was exceeded")
			} else {
				log.Fatalf("Unexpected error: %v", statusErr)
			}
		} else {
			log.Fatalf("Error while calling GreetWithDeadline RPC: %v", err)
		}
		return
	}
	log.Printf("Response from GreetWithDeadline: %v", res.Result)
}
