package xlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

const logProductionKey = "LOG_MODE" // "production" or "development" or "pretty"
const logAppNameKey = "APP_NAME"    // default is localhost

const (
	development string = "development"
	production  string = "production"
	pretty      string = "pretty"

	ObjDebugPrettyPrint string = "$__pretty_print"
)

var logger *zap.Logger
var mode = development
var appName = "local"

func init() {
	if aName, ok := os.LookupEnv(logAppNameKey); ok {
		appName = aName
	}
	pe := zap.NewDevelopmentConfig()
	if mdos, ok := os.LookupEnv(logProductionKey); ok {
		md := strings.ToLower(mdos)
		if md == production {
			pe = zap.NewProductionConfig()
			mode = production
		} else if md == pretty {
			mode = pretty
		}
	}
	pe.EncoderConfig.TimeKey = "timestamp"
	pe.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		zapcore.RFC3339NanoTimeEncoder(t.UTC(), enc)
	}

	if pe.Development {
		pe.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	pe.OutputPaths = []string{"stdout"}
	opt := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.PanicLevel),
	}

	logger, _ = pe.Build(opt...)
}

func Debug() *LogEvent {
	return newLogEvent(zapcore.DebugLevel)
}

func Info() *LogEvent {
	return newLogEvent(zapcore.InfoLevel)
}

func Warn() *LogEvent {
	return newLogEvent(zapcore.WarnLevel)
}

func Error() *LogEvent {
	return newLogEvent(zapcore.ErrorLevel)
}

func Panic() *LogEvent {
	return newLogEvent(zapcore.PanicLevel)
}

func Fatal() *LogEvent {
	return newLogEvent(zapcore.FatalLevel)
}

func PrettyPrint(obj ...any) {
	for _, o := range obj {
		Debug().AddCallerSkip(1).Pretty().Field(ObjDebugPrettyPrint, o).Msg("pretty print")
	}
}
