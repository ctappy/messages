package message

import (
	"fmt"
	"github.com/ctaperts/messages/log"
	"github.com/ctaperts/messages/message_email"
	"github.com/ctaperts/messages/src"
	"github.com/nlopes/slack"
	"log"
	"regexp"
	"strings"
)

const (
	helpMessage = "```" + `
Message Bot

Usage:
  @%s [command]

Available Commands:
  email      Email from slackbot
  shutdown   Shutdown slackbot
  help       Show this command
` + "```"
	helpEmailMessage = "```" + `
Send email

Usage:
  @%s email <from address> <to address> "<subject>" "<body>"

Additional info:
  "\n" in body for newline
` + "```"
)

func formatEmailText(text string) (emailFrom, emailSubject, emailBody string, emailTo []string) {
	r := regexp.MustCompile("'.+'|\".+\"|\\S+")
	rSubjectBody := regexp.MustCompile("\"(.*?)\"")
	m := r.FindAllString(text, -1)
	sb := rSubjectBody.FindAllString(text, -1)
	if len(m) != 5 {
		return
	}
	if len(sb) != 2 {
		return
	}
	// check for email typo
	if !strings.Contains(m[2], "mailto") {
		return
	} else if !strings.Contains(m[3], "mailto") {
		return
	}
	rEmail := regexp.MustCompile("<mailto:(.*)\\|")
	emailFromSetup := rEmail.FindStringSubmatch(m[2])
	emailFrom = emailFromSetup[1]
	emailToSetup := m[3]
	emailTo = rEmail.FindStringSubmatch(emailToSetup)[1:]
	if strings.Contains(emailToSetup, ";") {
		emailTo = append(emailTo, strings.Split(emailToSetup, ";")[1:]...)
	}
	emailSubject = strings.Replace(sb[0], "\"", "", 2)
	emailBody = strings.Replace(sb[1], "\"", "", 2)
	return
}

func Event(Log logging.Logs, LocalConfig configuration.Config, bot configuration.BotInfo, ev *slack.MessageEvent, rtm *slack.RTM, done chan struct{}) bool {
	// Check text for exact message
	switch text := ev.Text; text {
	case fmt.Sprintf("<@%s> shutdown", bot.ID):
		log.Println("Shutting down message received")
		rtm.SendMessage(rtm.NewOutgoingMessage("Shutting down now...", bot.ChannelID))
		rtm.Disconnect()
		done <- struct{}{}
		return true
	case fmt.Sprintf("<@%s> help", bot.ID):
		rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf(helpMessage, bot.Name), bot.ChannelID))
		return true
	case fmt.Sprintf("<@%s> email help", bot.ID):
		rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf(helpEmailMessage, bot.Name), bot.ChannelID))
		return true
	}
	// Check text for commands with args
	switch {
	case strings.HasPrefix(ev.Text, fmt.Sprintf("<@%s> email", bot.ID)):
		rtm.SendMessage(rtm.NewOutgoingMessage("Sending email..", bot.ChannelID))
		emailFrom, emailSubject, emailBody, emailTo := formatEmailText(ev.Text)
		if email.Send(LocalConfig, emailFrom, emailSubject, emailBody, emailTo) {
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("Email Sent to %s", emailTo), bot.ChannelID))
			Log.Email.Println(fmt.Sprintf("TO: %s, FROM: %s, SUBJECT: %s, BODY: %s, SLACK-CHANNEL: %s", emailTo, emailFrom, emailSubject, emailBody, bot.ChannelName))
		} else {
			rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("Failed to send email run `<@%s> email help` for more info", bot.ID), bot.ChannelID))
		}
		return true
	default:
		return false
	}
}
