package app

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Command a sub command structure of the cli application
type Command struct {
	usage string
	desc  string
	// commands is the list of commands supported by this program
	commands []*Command

	// options command line configuration
	options CmdLineOptioner

	// runFun
	runFunc RunCommandFunc
}

// CommandOption optional parameters for initializing the command
type CommandOption func(*Command)

// WithCommandOptions allows the app to read from the flags from command line
func WithCommandOptions(opt CmdLineOptioner) CommandOption {
	return func(c *Command) {
		c.options = opt
	}
}

// RunCommandFunc the app's command startup callback function
type RunCommandFunc func(args []string) error

func WithCommandRunFunc(runFunc RunCommandFunc) CommandOption {
	return func(c *Command) {
		c.runFunc = runFunc
	}
}

// NewCommand creates a new sub command instance
// with command name and others options
func NewCommand(usage, desc string, opts ...CommandOption) *Command {
	c := &Command{
		usage: usage,
		desc:  desc,
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

// cobraCommand create sub command and generate cmd tree
func (c *Command) cobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.usage,
		Short: c.desc,
	}

	cmd.SetOut(os.Stdout)
	cmd.Flags().SortFlags = true

	// add sub commands
	if len(c.commands) > 0 {
		for _, command := range c.commands {
			cmd.AddCommand(command.cobraCommand())
		}
	}

	// runFunc
	if c.runFunc != nil {
		cmd.Run = c.runCommand
	}

	// sub command flags
	if c.options != nil {
		for _, f := range c.options.Flags().FlagSets {
			cmd.Flags().AddFlagSet(f)
		}
	}

	addHelpCommandFlag(c.usage, cmd.Flags())

	return cmd
}

// runCommand sub cobra command run
func (c *Command) runCommand(cmd *cobra.Command, args []string) {
	if c.runFunc != nil {
		if err := c.runFunc(args); err != nil {
			fmt.Printf("%v %v\n", color.RedString("Error:"), err)
			os.Exit(1)
		}
	}
}

// FormatBaseName generates a executable file name under different os
// based on basename
func FormatBaseName(basename string) string {
	if runtime.GOOS == "windows" {
		// case-insensitive
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}

	return basename
}

func addCmdTemplate(cmd *cobra.Command, namedFlagSets NamedFlagSets) {
	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)

		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
}
