package main

import (
	"context"
	"fmt"
	"github.com/ctaperts/messages/messagepb"
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
	conn, err := grpc.Dial("127.0.0.1:50052", opts)

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	c := messagepb.NewMessageServiceClient(conn)

	fmt.Println("Creating the message")
	message := &messagepb.EmailMessage{
		To:      "colbytaperts@gmail.com",
		From:    "My first Email message",
		Subject: "Content of Email message",
		Body:    "Content of Email message",
	}
	createEmailMessageRes, err := c.CreateEmailMessage(context.Background(), &messagepb.CreateEmailMessageRequest{Message: message})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Email Message has been created: %v\n", createEmailMessageRes)

	messageID := createEmailMessageRes.GetMessage().GetId()
	fmt.Println("SlackMessage ID:", messageID)
	slackMessage := &messagepb.SlackMessage{
		Subject: "Content of slack message",
		Body:    "Content of slack message",
	}
	createMessageRes, err := c.CreateSlackMessage(context.Background(), &messagepb.CreateSlackMessageRequest{Message: slackMessage})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Slack Message has been created: %v\n", createMessageRes)
	messageID = createMessageRes.GetMessage().GetId()
	fmt.Println("Message ID:", messageID)

}
