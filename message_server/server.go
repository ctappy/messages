package main

import (
	"context"
	"github.com/ctaperts/messages/messagepb"
	"google.golang.org/grpc"
	"log"
	// "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	// "google.golang.org/grpc/status"
	"github.com/ctaperts/messages/message_email"
	"net"
	"os"
	"os/signal"
)

type server struct {
}

func (*server) CreateEmailMessage(ctx context.Context, req *messagepb.CreateEmailMessageRequest) (*messagepb.CreateEmailMessageResponse, error) {
	message := req.GetMessage()

	data := emailMessageItem{
		To:      message.GetTo(),
		From:    message.GetFrom(),
		Subject: message.GetSubject(),
		Body:    message.GetBody(),
	}

	log.Printf("From: %s, Subject: %s, Body: %s, To: %s\n", data.From, data.Subject, data.Body, data.To)
	if email.Send(data.From, data.Subject, data.Body, []string{data.To}) {
		log.Println("Email sent successfully")
	} else {
		log.Println("Email failed to send")
	}

	return &messagepb.CreateEmailMessageResponse{
		Message: &messagepb.EmailMessage{
			Id:      "",
			To:      message.GetTo(),
			From:    message.GetFrom(),
			Subject: message.GetSubject(),
			Body:    message.GetBody(),
		},
	}, nil
}

func (*server) CreateSlackMessage(ctx context.Context, req *messagepb.CreateSlackMessageRequest) (*messagepb.CreateSlackMessageResponse, error) {
	log.Println("Create slack message request")
	message := req.GetMessage()

	data := slackMessageItem{
		Subject: message.GetSubject(),
		Body:    message.GetBody(),
	}

	log.Println(data)

	return &messagepb.CreateSlackMessageResponse{
		Message: &messagepb.SlackMessage{
			Id:      "",
			Subject: message.GetSubject(),
			Body:    message.GetBody(),
		},
	}, nil
}

func dataToEmailMessage(data *emailMessageItem) *messagepb.EmailMessage {
	return &messagepb.EmailMessage{
		Id:      data.ID,
		To:      data.To,
		From:    data.From,
		Subject: data.Subject,
		Body:    data.Body,
	}
}

func dataToSlackMessage(data *slackMessageItem) *messagepb.SlackMessage {
	return &messagepb.SlackMessage{
		Id:      data.ID,
		Subject: data.Subject,
		Body:    data.Body,
	}
}

type slackMessageItem struct {
	ID      string `id`
	Subject string `subject`
	Body    string `body`
}

type emailMessageItem struct {
	ID      string `id`
	To      string `to`
	From    string `from`
	Subject string `subject`
	Body    string `body`
}

var debug *bool

func main() {
	log.Println("Message Service Started")

	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	tls := false
	opts := []grpc.ServerOption{}
	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	messagepb.RegisterMessageServiceServer(s, &server{})

	// evans cli
	// `evans -p 50052 -r`
	reflection.Register(s)

	go func() {
		log.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until signal is received
	<-ch
	log.Println("Stopping the server")
	s.Stop()
	log.Println("Closing the listener")
	lis.Close()
	log.Println("Application is stopped")
}
