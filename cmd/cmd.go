package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Root is cmdctrl's root command
var Root = &cobra.Command{
	Use: "cmdctrl",
	Short: "cmdctrl is a very simple device management system",
	Long: `cmdctrl is a very simple device management system

Allows offloading of various tasks to mobile devices,
running commands remotely on devices,
sending messages to mobile devices,
and so much more`,
	Version: Version,
}

// Execute executes the program
func Execute() {
	if err := Root.Execute(); err != nil {
		fmt.Printf("Error parsing command: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	Root.SetVersionTemplate("{{.Name}} - Version {{.Version}}\n")
}

// CheckArgs here is copied from github.com/ncw/rclone/cmd
// CheckArgs checks there are enough arguments and prints a message if not
func CheckArgs(MinArgs, MaxArgs int, cmd *cobra.Command, args []string) {
	if len(args) < MinArgs {
		_ = cmd.Usage()
		fmt.Println()
		_, _ = fmt.Fprintf(os.Stderr, "Command %s needs %d arguments minimum: you provided %d non flag arguments: %q\n", cmd.Name(), MinArgs, len(args), args)
		os.Exit(1)
	} else if len(args) > MaxArgs {
		_ = cmd.Usage()
		fmt.Println()
		_, _ = fmt.Fprintf(os.Stderr, "Command %s needs %d arguments maximum: you provided %d non flag arguments: %q\n", cmd.Name(), MaxArgs, len(args), args)
		os.Exit(1)
	}
}
