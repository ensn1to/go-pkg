package app

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// package variable
const (
	flagHelp          = "help"
	flagHelpShorthand = "h"
)

func helpCommand(name string) *cobra.Command {
	return &cobra.Command{
		Use:   "help [command]",
		Short: "Help info about any command",
		Long: `Help provides help for any command in the application.
	Simply type` + name + ` help [path to command] for more detail`,

		Run: func(cmd *cobra.Command, args []string) {
			cmd, _, e := cmd.Root().Find(args)
			if cmd == nil || e != nil {
				cmd.Printf("Unknown help topic %#q\n", args)
				cmd.Root().Usage()
			} else {
				cmd.InitDefaultHelpCmd()
				cmd.Help()
			}
		},
	}
}

// addHelpFlag for applicaiton
func addHelpFlag(name string, fs *pflag.FlagSet) {
	fs.BoolP(flagHelp, flagHelpShorthand, false, fmt.Sprintf("Help for %s", name))
}

// addHelpCommandFlag for command
func addHelpCommandFlag(usage string, fs *pflag.FlagSet) {
	fs.BoolP(flagHelp, flagHelpShorthand, false, fmt.Sprintf("Help for the %s command", color.GreenString(strings.Split(usage, " ")[0])))
}
