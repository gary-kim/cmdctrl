package server

import (
	"github.com/gary-kim/cmdctrl/cmd"
	"github.com/gary-kim/cmdctrl/server"
	"github.com/spf13/cobra"
)

func init() {
	address := "localhost:80"
	sharedPass := ""
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Run cmdctrl in server mode",
		Long: `Start cmdctrl in server mode

Use server mode on the command and control server
running on a central server. Clients will connect
to this system to recieve commands`,
		Run: func(command *cobra.Command, args []string) {
			opt := server.Options{
				SharedPass: sharedPass,
			}
			server.RunServer(address, opt)
		},
	}
	cmd.Root.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVar(&address, "addr", ":80", "Indicates the address you would like the cmdctrl server to bind to. Use (:port) to specify just a port for a full (host:addr) for the entire bind address")
	serverCmd.PersistentFlags().StringVar(&sharedPass, "shared-pass", "", "A shared pass for the server and client. Must be the same between the server and client. This is used by the client to authenticate the server.")
}
