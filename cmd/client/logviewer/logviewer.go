package logviewer

import (
	"github.com/gary-kim/cmdctrl/client/logviewer"
	"github.com/gary-kim/cmdctrl/cmd/client"
	"github.com/spf13/cobra"
)

func init() {
	// Set up all flag vars
	logLocation := "cmdctrl.log"

	viewLogCmd := &cobra.Command{
		Use:   "logviewer",
		Short: "Run cmdctrl in client mode",
		Long: `Start cmdctrl in client mode

Use client mode on devices to be managed by cmdctrl.
cmdctrl will connect to the specified server
and follow its instructions`,
		Run: func(command *cobra.Command, args []string) {
			err := logviewer.RunLogViewer(logviewer.Options{
				LogFile: logLocation,
			})
			if err != nil {
				return
			}
		},
	}
	client.ClientCmd.AddCommand(viewLogCmd)

	// Set flags
	viewLogCmd.Flags().StringVar(&logLocation, "log-file", "cmdctrl.log", "Specify a location in which to save cmdctrl client's logs")

}
