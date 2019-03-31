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

func getId(api *slack.Client) (userName, userID string) {

	// Find the bot info
	authTest, err := api.AuthTest()
	if err != nil {
		log.Panicf("Error getting info: %s\n", err)
		return
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
	api := slack.New(
		LocalConfig.Slack.SlackKey,
		slack.OptionDebug(trace),
		slack.OptionLog(log.New(os.Stdout, "trace-slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	args := configuration.DefaultArgs()
	userName, userID := getId(api)
	args.BotName = userName
	args.BotID = userID
	if args.Debug {
		log.Println("arguments:", args)
	}
	Log.Info.Println("Starting RTM slackbot")
	Log.Debug.Println("Bot Name:", userName)
	Log.Debug.Println("Bot ID:", userID)

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
