package slack

import (
	"fmt"
	"github.com/ctaperts/messages/log"
	"github.com/ctaperts/messages/message_slack/message"
	"github.com/ctaperts/messages/src"
	"github.com/nlopes/slack"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	LocalConfig configuration.Config
	trace       bool
)

func getInfo(channelName string, api *slack.Client) (userName, userID, channelID string) {

	// Find the bot info
	authTest, err := api.AuthTest()
	if err != nil {
		log.Panicf("Error getting info: %s\n", err)
		return
	}
	channels, err := api.GetChannels(false)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}
	for _, channel := range channels {
		if channelName == channel.Name {
			channelID = channel.ID
		}
	}

	userID = authTest.UserID
	userName = authTest.User
	return
}

func Bot(logLevel string, Log logging.Logs) {
	if logLevel == "trace" {
		trace = true
		Log.Info.Println("Running slackbot in trace mode")
	} else {
		trace = false
	}

	LocalConfig = configuration.LoadConfig()
	Log.Info.Println("Slack OAuth Key:", LocalConfig.Slack.BotUserToken)
	Log.Info.Println("Channel Name:", LocalConfig.Slack.ChannelName)
	api := slack.New(
		LocalConfig.Slack.BotUserToken,
		slack.OptionDebug(trace),
		slack.OptionLog(log.New(os.Stdout, "trace-slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	bot := configuration.DefaultArgs()
	userName, userID, channelID := getInfo(LocalConfig.Slack.ChannelName, api)
	// setup struct with slackbot info
	bot.Name = userName
	bot.ID = userID
	bot.ChannelID = channelID
	bot.ChannelName = LocalConfig.Slack.ChannelName
	Log.Trace.Println("Bot Info:", bot)
	Log.Info.Println("Starting RTM slackbot")
	Log.Debug.Println("Bot Name:", userName)
	Log.Debug.Println("Bot ID:", userID)

	rtm := api.NewRTM(slack.RTMOptionUseStart(true))
	go rtm.ManageConnection()

	// Start notice
	rtm.SendMessage(rtm.NewOutgoingMessage("Starting messagebot", channelID))

	// Observe incoming messages.
	done := make(chan struct{})
	connectingReceived := false
	connectedReceived := false
	go func() {
		for msg := range rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.ConnectingEvent:
				if connectingReceived {
					Log.Fatal.Panicf("Received multiple connecting events.\n")
				}
				connectingReceived = true
			case *slack.ConnectedEvent:
				if connectedReceived {
					Log.Fatal.Panicf("Received multiple connecting events.\n")
				}
				connectedReceived = true
			// Check messages in channel
			case *slack.MessageEvent:
				if message.Event(Log, LocalConfig, bot, ev, rtm, done) == false {
					Log.Debug.Printf("Discarding message with content %+v\n", ev)
					Log.Info.Printf("Text: %+v\n", ev.Text)
				}
			default:
				Log.Debug.Printf("Discarded event of type '%s' with content '%#v'\n", msg.Type, ev)
			}
		}
	}()

	// handle interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			Log.Debug.Printf("%s signal received, shutting down\n", sig)
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("%s is shutting down from %s signal", bot.Name, sig), channelID))
			rtm.Disconnect()
			done <- struct{}{}
		}
	}()

	<-done

}
