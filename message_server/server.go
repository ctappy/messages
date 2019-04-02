package grpc

import (
	"context"
	"github.com/ctaperts/messages/messagepb"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	// "google.golang.org/grpc/status"
	"github.com/ctaperts/messages/log"
	"github.com/ctaperts/messages/message_email"
	// "github.com/ctaperts/messages/message_slack/message"
	"github.com/ctaperts/messages/src"
	"net"
	"os"
	"os/signal"
)

type server struct {
}

var (
	LocalConfig configuration.Config
	Log         logging.Logs
	Logs        *logging.Logs
)

func (*server) CreateEmailMessage(ctx context.Context, req *messagepb.CreateEmailMessageRequest) (*messagepb.CreateEmailMessageResponse, error) {
	message := req.GetMessage()

	data := emailMessageItem{
		To:      message.GetTo(),
		From:    message.GetFrom(),
		Subject: message.GetSubject(),
		Body:    message.GetBody(),
	}

	Log.Email.Printf("GRPC: From: %s, Subject: %s, Body: %s, To: %s\n", data.From, data.Subject, data.Body, data.To)
	if email.Send(LocalConfig, data.From, data.Subject, data.Body, []string{data.To}) {
		Log.Debug.Println("Email sent successfully")
	} else {
		Log.Debug.Println("Email failed to send")
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
	Log.Info.Println("Create slack message request")
	message := req.GetMessage()

	data := slackMessageItem{
		Subject: message.GetSubject(),
		Body:    message.GetBody(),
	}

	Log.Debug.Println(data)

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

func Exec(Logs logging.Logs) {
	// Set global variable
	Log = Logs

	LocalConfig = configuration.LoadConfig()
	Log.Info.Println("Message Service Started")

	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		Log.Err.Fatalf("Failed to listen: %v", err)
	}

	tls := false
	opts := []grpc.ServerOption{}
	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			Log.Err.Fatalf("Failed loading certificates: %v", sslErr)
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
		Log.Info.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			Log.Err.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until signal is received
	<-ch
	Log.General.Println("Stopping the grpc server")
	s.Stop()
	Log.General.Println("Closing the listener")
	lis.Close()
	Log.General.Println("Application is stopped")
}
