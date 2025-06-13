package xlog

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kurzgesagtz/xgo/xerror"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type rawEvent struct {
	fields map[string]any
	data   map[string]any
}

type LogEvent struct {
	level      zapcore.Level
	appName    string
	err        error
	fields     []zapcore.Field
	data       []zapcore.Field
	raw        rawEvent
	pretty     bool
	callerSkip int
}

func newLogEvent(level zapcore.Level) *LogEvent {
	ev := &LogEvent{
		level:   level,
		data:    make([]zapcore.Field, 0),
		fields:  make([]zapcore.Field, 0),
		appName: appName,
		pretty:  false,
		raw: rawEvent{
			fields: make(map[string]any),
			data:   make(map[string]any),
		},
	}
	ev.fields = append(ev.fields, zap.Any("app_name", appName))
	return ev
}

func (l *LogEvent) addField(key string, value any) {
	l.fields = append(l.fields, zap.Any(key, value))
	l.raw.fields[key] = value
}

func (l *LogEvent) Err(err error) *LogEvent {
	if err == nil {
		return l
	}
	l.err = err
	l.fields = append(l.fields, zap.Error(err))
	l.raw.fields["error"] = err.Error()
	var xErr *xerror.Error
	if errors.As(err, &xErr) {
		if xErr.Caller != "" {
			l.addField("error_caller", xErr.Caller)
		}
		if xErr.AppName != "local" && xErr.AppName != "" {
			l.addField("error_app_name", xErr.AppName)
		}
	}
	return l
}

func (l *LogEvent) Context(ctx context.Context) *LogEvent {
	if gCtx, ok := ctx.(*gin.Context); ok {
		l.addField("ip_address", gCtx.ClientIP())
		l.addField("user_agent", gCtx.Request.UserAgent())
		l.addField("method", gCtx.Request.Method)
		l.addField("path", gCtx.Request.URL.Path)

		ctx = gCtx.Request.Context()
	}
	if span := trace.SpanContextFromContext(ctx); span.IsValid() {
		l.addField("span_id", span.SpanID().String())
		l.addField("trace_id", span.TraceID().String())
		l.addField("trace_sample", span.IsSampled())
	}
	return l
}

func (l *LogEvent) AddCallerSkip(n int) *LogEvent {
	l.callerSkip += n
	return l
}

func (l *LogEvent) Field(key string, val any) *LogEvent {
	l.data = append(l.data, zap.Any(key, val))
	l.raw.data[key] = val
	return l
}

func (l *LogEvent) Pretty() *LogEvent {
	l.pretty = true
	return l
}

func (l *LogEvent) Msg(msg string) {
	if mode == pretty || (mode == development && l.pretty) {
		logger.WithOptions(zap.AddCallerSkip(l.callerSkip)).Log(l.level, msg)

		raw := make(map[string]interface{})
		if len(l.raw.fields) > 0 {
			raw = l.raw.fields
		}
		if len(l.raw.data) > 0 {
			raw["data"] = l.raw.data
		}
		if data, ok := raw["data"]; ok {
			if obj, ok := data.(map[string]interface{}); ok {
				if _obj, ok := obj[ObjDebugPrettyPrint]; ok {
					clr := _prettyLevelToColor[l.level]
					if err := printJSON(&clr, _obj); err != nil {
						fmt.Println(clr.Add(fmt.Sprintf("%v", _obj)))
					}
					return
				}
			}
		}
		if len(raw) > 0 {
			clr := _prettyLevelToColor[l.level]
			if err := printJSON(&clr, raw); err != nil {
				panic(err)
			}
		}
	} else {
		logger.WithOptions(zap.AddCallerSkip(l.callerSkip)).
			Log(l.level, msg, append(l.fields, zap.Any("data", l.data))...)
	}
}
