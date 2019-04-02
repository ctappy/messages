package message

import (
	"fmt"
	"github.com/ctaperts/messages/log"
	"github.com/nlopes/slack"
)

func postAttachment(Log logging.Logs, api *slack.Client, title, message_type, subject, body, channelID string) {
	attachment := slack.Attachment{
		Text: fmt.Sprintf("%s", "type"),
		// Uncomment the following part to send a field too
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: fmt.Sprintf("%s", "subject"),
				Value: fmt.Sprintf("%s", "body"),
			},
		},
	}

	channelID, timestamp, err := api.PostMessage(channelID, slack.MsgOptionText(fmt.Sprintf("%s", "title"), false), slack.MsgOptionAttachments(attachment))
	if err != nil {
		Log.Err.Printf("%s\n", err)
	}
	Log.Email.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}
