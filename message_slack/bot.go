package main

import (
	"encoding/json"
	"flag"
	"github.com/nlopes/slack"
	"io"
	"io/ioutil"
	"log"
	"os"
)

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

var LocalConfig Config
var debug *bool

const (
	shutDownMessage = "shutdown messagebot"
)

func loadConfig(jsonFile io.Reader) Config {
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	err := json.Unmarshal(byteValue, &config)
	if err != nil {
		log.Fatalf("Failed to load json file %v", err)
	}
	return config
}

func getId(api *slack.Client) (userName, userID string) {

	// Find the bot info
	authTest, err := api.AuthTest()
	if err != nil {
		log.Printf("Error getting info: %s\n", err)
		return
	}

	userID = authTest.UserID
	userName = authTest.User
	return
}

func SlackBot() {
	api := slack.New(
		LocalConfig.Slack.SlackKey,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	userName, userID := getId(api)
	log.Println("Starting RTM slackbot")
	log.Println("Bot Name:", userName)
	log.Println("Bot ID:", userID)

	rtm := api.NewRTM(slack.RTMOptionUseStart(true))
	go rtm.ManageConnection()

	// Start notice
	rtm.SendMessage(rtm.NewOutgoingMessage("Starting messagebot", LocalConfig.Slack.ChannelID))

	// Observe incoming messages.
	done := make(chan struct{})
	connectingReceived := false
	connectedReceived := false
	go func() {
		for msg := range rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.ConnectingEvent:
				if connectingReceived {
					log.Panicf("Received multiple connecting events.\n")
				}
				connectingReceived = true
			case *slack.ConnectedEvent:
				if connectedReceived {
					log.Panicf("Received multiple connecting events.\n")
				}
				connectedReceived = true
			// Check messages in channel
			case *slack.MessageEvent:
				if messageEvent(ev, rtm, done) == false {
					log.Printf("Discarding message with content %+v\n", ev)
				}
			default:
				log.Printf("Discarded event of type '%s' with content '%#v'\n", msg.Type, ev)
			}
		}
	}()

	<-done

}
func messageEvent(ev *slack.MessageEvent, rtm *slack.RTM, done chan struct{}) bool {
	if ev.Text == shutDownMessage {
		log.Println("Shutting down message received")
		rtm.SendMessage(rtm.NewOutgoingMessage("Shutting down now...", "CANSP3PRR"))
		rtm.Disconnect()
		done <- struct{}{}
		return true
	} else {
		return false
	}
}

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
			log.Printf("Loading configuration from %q\n", *configPtr)
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
		log.Printf("Successfully Opened %q\n", *configPtr)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	LocalConfig = loadConfig(jsonFile)
}

func main() {
	SlackBot()
}
