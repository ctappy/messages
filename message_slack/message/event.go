package message

import (
	"github.com/nlopes/slack"
	"log"
)

const (
	shutDownMessage = "shutdown messagebot"
)

func Event(ev *slack.MessageEvent, rtm *slack.RTM, done chan struct{}) bool {
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
