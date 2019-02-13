package main

import (
	"context"
	"fmt"
	"github.com/ctaperts/grpc-messages/messagepb"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	// "google.golang.org/grpc/status"
	// "io"
	"log"
	// "time"
)

func main() {

	fmt.Println("Message Client")

	tls := false
	opts := grpc.WithInsecure()
	if tls {
		certFile := "ssl/ca.crt" // ca auth trust cert
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
		}
		opts = grpc.WithTransportCredentials(creds)
	}
	conn, err := grpc.Dial("localhost:50052", opts)

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	c := messagepb.NewMessageServiceClient(conn)

	fmt.Println("Creating the message")
	message := &messagepb.Message{
		Slack:   false,
		Email:   true,
		To:      "Colby",
		From:    "My first message",
		Subject: "Content of message",
		Body:    "Content of message",
	}
	createMessageRes, err := c.CreateMessage(context.Background(), &messagepb.CreateMessageRequest{Message: message})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Message has been created: %v\n", createMessageRes)
	messageID := createMessageRes.GetMessage().GetId()
	fmt.Println("Message ID:", messageID)

}
