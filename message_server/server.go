package grpc

import (
	"context"
	"github.com/ctaperts/messages/messagepb"
	"google.golang.org/grpc"
	"net/http"
	// "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	// "google.golang.org/grpc/status"
	"github.com/ctaperts/messages/log"
	"github.com/ctaperts/messages/message_email"
	"github.com/ctaperts/messages/message_slack"
	"github.com/ctaperts/messages/message_slack/message"
	"github.com/ctaperts/messages/src"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nlopes/slack"
	"log"
	"net"
	"os"
	"os/signal"
)

const grpcIPPort = "0.0.0.0:50052"
const httpIPPort = "0.0.0.0:8080"

type server struct {
}

var (
	LocalConfig configuration.Config
	Log         logging.Logs
	Logs        *logging.Logs
	api         *slack.Client
	channelID   string
	trace       bool
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

	slackMessage.PostAttachment(Log, api, "grpc title", "grpc type", data.Subject, data.Body, channelID)
	// slackMessage.PostOptions(Log, api, data.Subject, data.Body, channelID)

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

func Exec(logLevel string, Logs logging.Logs) {
	if logLevel == "trace" {
		trace = true
	} else {
		trace = false
	}
	// Set global variable
	Log = Logs

	LocalConfig = configuration.LoadConfig()
	getApi := slack.New(
		LocalConfig.Slack.BotUserToken,
		slack.OptionDebug(trace),
		slack.OptionLog(log.New(os.Stdout, "TRACE-SLACK-BOT: ", log.Lshortfile|log.LstdFlags)),
	)
	api = getApi
	_, _, getChannelID := bot.GetInfo(LocalConfig.Slack.ChannelName, api)
	channelID = getChannelID
	Log.Info.Println("Message Service Started")

	lis, err := net.Listen("tcp", grpcIPPort)
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
		Log.Info.Println("Starting GRPC Server...")
		if err := s.Serve(lis); err != nil {
			Log.Err.Fatalf("failed to serve: %v", err)
		}
	}()

	Log.Info.Println("Starting HTTP Server...")
	go run()
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

func run() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := messagepb.RegisterMessageServiceHandlerFromEndpoint(ctx, mux, "localhost:50052", opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("HTTP Listening on %s\n", httpIPPort)
	log.Fatal(http.ListenAndServe(httpIPPort, mux))

}
