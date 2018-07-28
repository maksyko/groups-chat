package app

import (
	"go.uber.org/zap"
)

var zapLogger, _ = zap.NewDevelopment()

// Logger is the global logger used by application
var Logger = zapLogger.Sugar()
