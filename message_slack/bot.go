package main

import (
	"encoding/json"
	"flag"
	"github.com/ctaperts/messages/message_slack/message/message"
	"github.com/ctaperts/messages/message_slack/src/variables"
	"github.com/nlopes/slack"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	LocalConfig Config
)

const (
	shutDownMessage = "shutdown messagebot"
)

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
	args := Arguments.Load()
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

func init() {
	LocalConfig = LoadConfig()
}

func main() {
	SlackBot()
}
