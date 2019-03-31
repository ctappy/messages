package core

import (
	"github.com/ctaperts/messages/message_slack"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func Exec() {
	rootCmd.AddCommand(startGRPCCmd)
	rootCmd.AddCommand(startSlackBotCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func setup(cmd *cobra.Command, args []string) {
	// Setup logging
	setLogLevel(logLevel)
	logInit(*logSettings)
	Log.info.Printf("Log level set to %s\n", logLevel)
}

func startGRPC(cmd *cobra.Command, args []string) {
	Log.info.Println("Starting GRPC Message Server")
}

func startSlackBot(cmd *cobra.Command, args []string) {
	Log.info.Println("Starting Slack Message Server")
	slack.Bot()
}

func initStartSlackBot(cmd *cobra.Command, args []string) error {
	log.Println("Initializing Slackbot")
	return nil
}
func initStartGRPC(cmd *cobra.Command, args []string) error {
	log.Println("Initializing GRPC")
	return nil
}
