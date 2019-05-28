package server

import (
	"github.com/gary-kim/cmdctrl/cmd"
	"github.com/gary-kim/cmdctrl/server"
	"github.com/spf13/cobra"
)

func init() {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Run cmdctrl in server mode",
		Long: `Start cmdctrl in server mode

Use server mode on the command and control server
running on a central server. Clients will connect
to this system to recieve commands`,
		Run: func(command *cobra.Command, args []string) {
			opt := server.Options{

			}
			server.Run("", opt)
		},
	}
	cmd.Root.AddCommand(serverCmd)
}
