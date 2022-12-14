package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Field = zapcore.Field
	Level = zapcore.Level
)

var Any = zap.Any
