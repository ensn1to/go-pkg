package log

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

const (
	flagLevel             = "log.level"
	flagDisableCaller     = "log.disable-caller"
	flagDisableStacktrace = "log.disable-stacktrace"
	flagFormat            = "log.format"
	flagEnableColor       = "log.enable-color"
	flagOutputPaths       = "log.output-paths"
	flagErrorOutputPaths  = "log.error-output-paths"
	flagDevelopment       = "log.development"
	flagName              = "log.name"

	consoleFormat = "console" // txt
	jsonFormat    = "json"

	keyRequestID = "requestID"
)

// Options
type Options struct {
	OutputPaths       []string `json:"output-paths" mapstructure:"output-paths"`
	ErrorOutputPaths  []string `json:"error-output-paths" mapstructure:"error-output-paths"`
	Level             string   `json:"level" mapstructure:"level"`                   // log-level
	Format            string   `json:"format" mapstructure:"format"`                 // log file output format, JSON or Console(txt)
	DisableCaller     bool     `json:"disable-caller" mapstructure:"disable-caller"` // show name,location and line No. of the funcation called
	DisableStacktrace bool     `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`
	EnableColor       bool     `json:"enable-color" mapstructure:"enable-color"`
	Development       bool     `json:"development" mapstructure:"development"`
	Name              string   `json:"name" mapstructure:"name"`                   // logger name
	CommonFields      []string `json:"common-fields" mapstructure:"common-fields"` // common log fields, eg: requestId, username
}

// NewOption creates Options with default params
func NewOptions() *Options {
	return &Options{
		Level:             zapcore.InfoLevel.String(),
		DisableCaller:     false,
		DisableStacktrace: false,
		Format:            consoleFormat,
		EnableColor:       false,
		Development:       false,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		CommonFields:      []string{keyRequestID},
	}
}

// Validate validates the options fields
func (o *Options) Validate() []error {
	var errs []error

	var zapLevel zapcore.Level
	// string to zaplevel type
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}

	return errs
}

// AddFlags adds flags for log to the specified FlagSet object
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, flagLevel, o.Level, "Minimum log output `LEVEL`.")
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, flagDisableStacktrace,
		o.DisableStacktrace, "Disable the log to record a stack trace for all messages at or above panic level.")
	fs.StringVar(&o.Format, flagFormat, o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths, "Error output paths of log.")
	fs.BoolVar(
		&o.Development,
		flagDevelopment,
		o.Development,
		"Development puts the logger in development mode, which changes "+
			"the behavior of DPanicLevel and takes stacktraces more liberally.",
	)
	fs.StringVar(&o.Name, flagName, o.Name, "The name of the logger.")
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
