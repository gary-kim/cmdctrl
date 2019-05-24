package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version number
var Version = "v0.0.1"

func init() {
	versionCmd := &cobra.Command{
		Use: "version",
		Short: "Print version number of cmdctrl",
		Run: func(command *cobra.Command, args []string) {
			fmt.Printf("cmdctrl - Version %s\n", Version)
		},
	}
	Root.AddCommand(versionCmd)
}