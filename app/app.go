package app

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// App base structure for constructs a cli application
type App struct {
	// binary name
	basename    string
	name        string
	description string

	// command line options
	options CmdLineOptioner

	// cobra run funciton
	runFunc RunFunc

	silence  bool
	noConfig bool

	// pflag commands
	commands []*Command

	args cobra.PositionalArgs

	// cobra command
	cmd *cobra.Command
}

// Option optional parameters for initializing the applicaiton
type Option func(*App)

// WithOptions reads from command line
func WithOptions(opt CmdLineOptioner) Option {
	return func(a *App) {
		a.options = opt
	}
}

func WithDescription(desc string) Option {
	return func(a *App) {
		a.description = desc
	}
}

func WithSilence() Option {
	return func(a *App) {
		a.silence = true
	}
}

// WithNoConfig not support config flag
func WithNoConfig() Option {
	return func(a *App) {
		a.noConfig = true
	}
}

// WithValidArgs validate the non-flag args
func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(a *App) {
		a.args = args
	}
}

// WithDefaultValidArgs default args validation funciton
func WithDefaultValidArgs() Option {
	return func(a *App) {
		a.args = cobra.NoArgs
	}
}

// RunCommand the app's startup callback function
type RunFunc func(basename string) error

func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}

// NewApp creates a new application instance
// with applicaiton name, binary name, and initial options
func NewApp(name, basename string, opts ...Option) *App {
	a := &App{
		name:     name,
		basename: basename,
	}

	for _, o := range opts {
		o(a)
	}

	a.buildCommand()

	return a
}

func (a *App) buildCommand() {
	rootCmd := cobra.Command{
		Use:           FormatBaseName(a.basename),
		Short:         a.name,
		Long:          a.description,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          a.args,
	}

	// set usage message print
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	// indicate command message
	rootCmd.Flags().SortFlags = true

	InitFlags(rootCmd.Flags())

	// add sub command
	if len(a.commands) > 0 {
		for _, command := range a.commands {
			rootCmd.AddCommand(command.cobraCommand())
		}
		rootCmd.SetHelpCommand(helpCommand(FormatBaseName(a.basename)))
	}

	// runFunc
	if a.runFunc != nil {
		rootCmd.RunE = a.runCommand
	}

	// add opts to the rootCmd flagsets
	var namedFlagSets NamedFlagSets
	if a.options != nil {
		namedFlagSets = a.options.Flags()
		fs := rootCmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}
	}

	// add config to the global flagsets
	if !a.noConfig {
		addConfigFlag(a.basename, namedFlagSets.FlagSet("global"))
	}
	// add help
	addHelpFlag(rootCmd.Name(), namedFlagSets.FlagSet("global"))
	// add global flageset
	rootCmd.Flags().AddFlagSet(namedFlagSets.FlagSet("global"))

	addCmdTemplate(&rootCmd, namedFlagSets)

	a.cmd = &rootCmd
}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	printWorkingDir()
	printFalgs(cmd.Flags())

	if !a.noConfig {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	if a.options != nil {
		return a.applyOptionRules()
	}

	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}

	return nil
}

func (a *App) applyOptionRules() error {
	if completeableOptions, ok := a.options.(CompleteableOptions); ok {
		return completeableOptions.Complete()
	}

	if errs := a.options.Validate(); len(errs) > 0 {
		// todo: https://github1s.com/kubernetes/apimachinery/blob/HEAD/pkg/util/errors/errors.go
		return errs[0]
	}

	return nil
}

// Run launches the application
func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}
