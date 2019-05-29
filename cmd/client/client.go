package client

import (
	"github.com/gary-kim/cmdctrl/client"
	"github.com/gary-kim/cmdctrl/cmd"
	"github.com/spf13/cobra"
)

func init() {
	// Set up all flag vars
	RESTMode := false
	RESTUpdateInterval := 60
	sharedPass := ""

	clientCmd := &cobra.Command{
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
			})
		},
	}
	cmd.Root.AddCommand(clientCmd)

	// Set flags
	clientCmd.PersistentFlags().BoolVar(&RESTMode, "rest-mode", false, "Contact the server with RESTful HTTP requests rather than using a websocket connection.")
	clientCmd.PersistentFlags().IntVar(&RESTUpdateInterval, "rest-update-interval", 60, "How often to query the server for updates when in rest mode")
	clientCmd.PersistentFlags().StringVar(&sharedPass, "shared-pass", "", "A shared pass for the server and client. Must be the same between the server and client. This is used by the client to authenticate the server.")

}
