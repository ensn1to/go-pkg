package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// alias for zap structure
type (
	Field = zapcore.Field
	Level = zapcore.Level
)

// alias for zap function
var Any = zap.Any
