package slackMessage

import (
	"fmt"
	"github.com/ctaperts/messages/log"
	"github.com/nlopes/slack"
)

func PostOptions(Log logging.Logs, api *slack.Client, title, text, channelID string) {
	attachment := slack.Attachment{
		Text:       "Send email to:",
		Color:      "#6676f4",
		CallbackID: "email",
		Actions: []slack.AttachmentAction{
			{
				Name: "Name",
				Type: "select",
				Options: []slack.AttachmentActionOption{
					{
						Text:  "Colby Taperts",
						Value: "colbytaperts@gmail.com",
					},
				},
			},
			{
				Name:  "cancel Name",
				Text:  "Clear",
				Type:  "button",
				Style: "danger",
			},
		},
	}

	channelID, timestamp, err := api.PostMessage(channelID, slack.MsgOptionText(fmt.Sprintf("%s", title), false), slack.MsgOptionAttachments(attachment))
	if err != nil {
		Log.Err.Printf("%s\n", err)
	}
	Log.Email.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func PostAttachment(Log logging.Logs, api *slack.Client, title, text, subject, body, channelID string) {
	attachment := slack.Attachment{
		Text: fmt.Sprintf("%s", text),
		// Uncomment the following part to send a field too
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: fmt.Sprintf("%s", subject),
				Value: fmt.Sprintf("%s", body),
				Short: true,
			},
		},
	}

	channelID, timestamp, err := api.PostMessage(channelID, slack.MsgOptionText(fmt.Sprintf("%s", title), false), slack.MsgOptionAttachments(attachment))
	if err != nil {
		Log.Err.Printf("%s\n", err)
	}
	Log.Email.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
