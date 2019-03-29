package core

import (
	"log"
	"os"

	"github.com/ctaperts/messages/message_slack"
	"github.com/spf13/cobra"
)

func Exec() {
	RootCmd.AddCommand(StartGRPCCmd)
	RootCmd.AddCommand(StartSlackBotCmd)

	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func startGRPC(cmd *cobra.Command, args []string) {
	log.Println("Starting GRPC Message Server")
}

func startSlackBot(cmd *cobra.Command, args []string) {
	log.Println("Starting Slack Message Server")
	bot.SlackBot()
}

func initStartSlackBot(cmd *cobra.Command, args []string) error {
	log.Println("Initializing Slackbot")
	return nil
}
func initStartGRPC(cmd *cobra.Command, args []string) error {
	log.Println("Initializing GRPC")
	return nil
}
