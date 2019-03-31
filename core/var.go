package core

import (
	"github.com/spf13/cobra"
)

var (
	logLevel string

	rootCmd = &cobra.Command{
		Use:              "message",
		Short:            "",
		Long:             "",
		PersistentPreRun: setup,
	}
	startGRPCCmd = &cobra.Command{
		Use:   "grpc",
		Short: "Start the grpc server",
		Long:  "",
		Args:  initStartGRPC,
		Run:   startGRPC,
	}
	startSlackBotCmd = &cobra.Command{
		Use:   "slackbot",
		Short: "Start the slackbot server",
		Long:  "",
		Args:  initStartSlackBot,
		Run:   startSlackBot,
	}
)

func init() {
	// Logs require PersistentFlags to be available to all flag arguments
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "error", "Show log level output")
}