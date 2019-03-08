package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ctaperts/grpc-messages/messagepb"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	// "net/smtp"
	// "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	// "google.golang.org/grpc/status"
	"flag"
	"log"
	"os"
	"os/signal"
	// "io"
	"net"
	// "strconv"
	// "time"
)

type server struct {
}

func (*server) CreateEmailMessage(ctx context.Context, req *messagepb.CreateEmailMessageRequest) (*messagepb.CreateEmailMessageResponse, error) {
	fmt.Println("Create message request")
	message := req.GetMessage()

	data := emailMessageItem{
		To:      message.GetTo(),
		From:    message.GetFrom(),
		Subject: message.GetSubject(),
		Body:    message.GetBody(),
	}

	fmt.Println(data)
	fmt.Println(LocalConfig)

	// // Connect to the remote SMTP server.
	// c, err := smtp.Dial("smtp.gmail.com:465")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer c.Close()
	// // Set the sender and recipient.
	// c.Mail("colbytaperts@gmail.com")
	// c.Rcpt("colbytaperts@gmail.com")
	// // Send the email body.
	// wc, err := c.Data()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer wc.Close()
	// buf := bytes.NewBufferString("This is the email body.")
	// if _, err = buf.WriteTo(wc); err != nil {
	// 	log.Fatal(err)
	// }

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
	fmt.Println("Create message request")
	message := req.GetMessage()

	data := slackMessageItem{
		Subject: message.GetSubject(),
		Body:    message.GetBody(),
	}

	fmt.Println(data)

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

type Config struct {
	SMTP struct {
		Server   string `json:"server"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"smtp"`
	Slack struct {
		SlackKey  string `json:"slack_key"`
		ChannelID string `json:"channel_id"`
	} `json:"slack"`
}

func loadConfig(jsonFile io.Reader) Config {
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	err := json.Unmarshal(byteValue, &config)
	if err != nil {
		log.Fatalf("Failed to load json file %v", err)
	}
	return config
}

// TODO replace with context
/* global variable declaration */
var LocalConfig Config
var debug *bool

// init
func init() {
	// if we crash the go code, output file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// setup flags
	configPtr := flag.String("config", "./config.json", "JSON config file location")
	debug = flag.Bool("debug", false, "debug option")
	flag.Parse()

	// load json
	if _, err := os.Stat(*configPtr); err == nil {
		if *debug {
			fmt.Printf("Loading configuration from %q\n", *configPtr)
		}
	} else if os.IsNotExist(err) {
		log.Fatalf("File not found %q %v\n", *configPtr, err)
	} else {
		log.Fatalf("Issue finding file %q %v\n", *configPtr, err)
	}
	jsonFile, err := os.Open(*configPtr)
	if err != nil {
		log.Fatalf("Failed to open %q %v", *configPtr, err)
	}
	if *debug {
		fmt.Printf("Successfully Opened %q\n", *configPtr)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	LocalConfig = loadConfig(jsonFile)
}

func main() {
	if *debug {
		fmt.Println("loaded:", LocalConfig)
	}

	fmt.Println("Message Service Started")

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
		fmt.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("Application is stopped")
}
