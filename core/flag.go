package core

import (
	"github.com/ctaperts/messages/log"
	"github.com/ctaperts/messages/message_slack"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	Log logging.Logs
)

func Exec() {
	rootCmd.AddCommand(startGRPCCmd)
	rootCmd.AddCommand(startSlackBotCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func setup(cmd *cobra.Command, args []string) {
	// Setup logging
	logging.SetupLogLevel(logDir, logLevel)
	logging.LogInit(logDir, *logging.LogSettings)
	Log = logging.Log
	Log.Info.Printf("Log level set to %s\n", logLevel)
}

func startGRPC(cmd *cobra.Command, args []string) {
	Log.Info.Println("Starting GRPC Message Server")
}

func startSlackBot(cmd *cobra.Command, args []string) {
	Log.Info.Println("Starting Slack Message Server")
	slack.Bot(logLevel, Log)
}

func initStartSlackBot(cmd *cobra.Command, args []string) error {
	log.Println("Initializing Slackbot")
	return nil
}
func initStartGRPC(cmd *cobra.Command, args []string) error {
	log.Println("Initializing GRPC")
	return nil
}
