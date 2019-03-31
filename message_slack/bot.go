package slack

import (
	"github.com/ctaperts/messages/log"
	"github.com/ctaperts/messages/message_slack/message"
	"github.com/ctaperts/messages/src"
	"github.com/nlopes/slack"
	"log"
	"os"
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
	Log.Trace.Println("Config:", LocalConfig.Slack)
	Log.Info.Println("Slack OAuth Key:", LocalConfig.Slack.BotUserToken)
	Log.Info.Println("Channel Name:", LocalConfig.Slack.ChannelName)
	api := slack.New(
		LocalConfig.Slack.BotUserToken,
		slack.OptionDebug(trace),
		slack.OptionLog(log.New(os.Stdout, "trace-slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	args := configuration.DefaultArgs()
	userName, userID, channelID := getInfo(LocalConfig.Slack.ChannelName, api)
	args.BotName = userName
	args.BotID = userID
	Log.Trace.Println("arguments:", args)
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
				if message.Event(ev, rtm, done) == false {
					Log.Debug.Printf("Discarding message with content %+v\n", ev)
				}
			default:
				Log.Debug.Printf("Discarded event of type '%s' with content '%#v'\n", msg.Type, ev)
			}
		}
	}()

	<-done

}
