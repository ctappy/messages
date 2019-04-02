package core

import (
	"github.com/ctaperts/messages/log"
	"github.com/ctaperts/messages/message_server"
	"github.com/ctaperts/messages/message_slack"
	"github.com/spf13/cobra"
	"os"
)

var (
	Log logging.Logs
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
	logging.SetupLog(logDir, logLevel)
	logging.LogInit(*logging.LogSettings)
	Log = logging.Log
	Log.Info.Printf("Log level set to %s\n", logLevel)
}

func startGRPC(cmd *cobra.Command, args []string) {
	Log.Info.Println("Starting GRPC Message Server")
	grpc.Exec(Log)
}

func startSlackBot(cmd *cobra.Command, args []string) {
	Log.Info.Println("Starting Slack Message Server")
	slack.Bot(logLevel, Log)
}

func initStartSlackBot(cmd *cobra.Command, args []string) error {
	return nil
}
func initStartGRPC(cmd *cobra.Command, args []string) error {
	return nil
}
