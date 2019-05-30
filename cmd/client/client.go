package client

import (
	"github.com/gary-kim/cmdctrl/client"
	"github.com/gary-kim/cmdctrl/cmd"
	"github.com/spf13/cobra"
)

// ClientCmd is the client cobra command
var ClientCmd *cobra.Command

func init() {
	// Set up all flag vars
	RESTMode := false
	RESTUpdateInterval := 60
	sharedPass := ""
	logLocation := "cmdctrl.log"

	ClientCmd = &cobra.Command{
		Use:   "client",
		Short: "Run cmdctrl in client mode",
		Long: `Start cmdctrl in client mode

Use client mode on devices to be managed by cmdctrl.
cmdctrl will connect to the specified server
and follow its instructions`,
		Run: func(command *cobra.Command, args []string) {
			cmd.CheckArgs(1, 1, command, args)
			client.RunClient(args[0], client.Options{
				RESTMode:           RESTMode,
				RESTUpdateInterval: RESTUpdateInterval,
				SharedPass:         sharedPass,
				LogFile:            logLocation,
			})
		},
	}
	cmd.Root.AddCommand(ClientCmd)

	// Set flags
	ClientCmd.Flags().BoolVar(&RESTMode, "rest-mode", false, "Contact the server with RESTful HTTP requests rather than using a websocket connection.")
	ClientCmd.Flags().IntVar(&RESTUpdateInterval, "rest-update-interval", 60, "How often to query the server for updates when in rest mode")
	ClientCmd.Flags().StringVar(&sharedPass, "shared-pass", "", "A shared pass for the server and client. Must be the same between the server and client. This is used by the client to authenticate the server.")
	ClientCmd.Flags().StringVar(&logLocation, "log-file", "cmdctrl.log", "Specify a location in which to save cmdctrl client's logs")

}
