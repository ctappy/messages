package core

import "github.com/spf13/cobra"

var (
	Debug bool

	RootCmd = &cobra.Command{
		Use:   "message",
		Short: "",
		Long:  "",
	}
	StartGRPCCmd = &cobra.Command{
		Use:   "grpc",
		Short: "Start the grpc server",
		Long:  "",
		Args:  initStartGRPC,
		Run:   startGRPC,
	}
	StartSlackBotCmd = &cobra.Command{
		Use:   "slackbot",
		Short: "Start the slackbot server",
		Long:  "",
		Args:  initStartSlackBot,
		Run:   startSlackBot,
	}
)

func init() {
	StartGRPCCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Show debug output")
	StartSlackBotCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Show debug output")
}
