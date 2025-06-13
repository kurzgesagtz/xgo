package xlog

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kurzgesagtz/xgo/xerror"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewLogEvent(t *testing.T) {
	// Test creating a new log event with different levels
	levels := []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
		zapcore.PanicLevel,
		zapcore.FatalLevel,
	}

	for _, level := range levels {
		event := newLogEvent(level)

		if event == nil {
			t.Errorf("Expected non-nil LogEvent for level %v", level)
			continue
		}

		if event.level != level {
			t.Errorf("Expected level %v, got %v", level, event.level)
		}

		if event.appName != appName {
			t.Errorf("Expected appName %s, got %s", appName, event.appName)
		}

		if event.pretty {
			t.Error("Expected pretty to be false by default")
		}

		if len(event.fields) == 0 {
			t.Error("Expected non-empty fields")
		}

		// Check that app_name field is set
		found := false
		for _, field := range event.fields {
			if field.Key == "app_name" {
				found = true
				break
			}
		}

		if !found {
			t.Error("Expected app_name field to be set")
		}
	}
}

func TestLogEvent_Err(t *testing.T) {
	// Test with nil error
	event := newLogEvent(zapcore.InfoLevel)
	result := event.Err(nil)

	if result != event {
		t.Error("Expected Err to return the same LogEvent instance")
	}

	if event.err != nil {
		t.Error("Expected err to be nil")
	}

	// Test with standard error
	stdErr := errors.New("standard error")
	event = newLogEvent(zapcore.InfoLevel)
	result = event.Err(stdErr)

	if result != event {
		t.Error("Expected Err to return the same LogEvent instance")
	}

	if event.err != stdErr {
		t.Errorf("Expected err to be %v, got %v", stdErr, event.err)
	}

	// Check that error field is set
	found := false
	for _, field := range event.fields {
		if field.Key == "error" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected error field to be set")
	}

	// Test with xerror.Error
	xErr := xerror.NewError("TEST_CODE", xerror.WithMessage("test message"))
	event = newLogEvent(zapcore.InfoLevel)
	result = event.Err(xErr)

	if result != event {
		t.Error("Expected Err to return the same LogEvent instance")
	}

	if event.err != xErr {
		t.Errorf("Expected err to be %v, got %v", xErr, event.err)
	}

	// Check that error_caller field is set
	found = false
	for _, field := range event.fields {
		if field.Key == "error_caller" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected error_caller field to be set")
	}
}

func TestLogEvent_Context(t *testing.T) {
	// Test with standard context
	ctx := context.Background()
	event := newLogEvent(zapcore.InfoLevel)
	result := event.Context(ctx)

	if result != event {
		t.Error("Expected Context to return the same LogEvent instance")
	}

	// Test with gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var err error
	c.Request, err = http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	c.Request.Header.Set("User-Agent", "test-agent")

	event = newLogEvent(zapcore.InfoLevel)
	result = event.Context(c)

	if result != event {
		t.Error("Expected Context to return the same LogEvent instance")
	}

	// Check that IP address, user agent, method, and path fields are set
	fields := map[string]bool{
		"ip_address": false,
		"user_agent": false,
		"method":     false,
		"path":       false,
	}

	for _, field := range event.fields {
		if _, ok := fields[field.Key]; ok {
			fields[field.Key] = true
		}
	}

	for field, found := range fields {
		if !found {
			t.Errorf("Expected %s field to be set", field)
		}
	}
}

func TestLogEvent_AddCallerSkip(t *testing.T) {
	event := newLogEvent(zapcore.InfoLevel)
	result := event.AddCallerSkip(2)

	if result != event {
		t.Error("Expected AddCallerSkip to return the same LogEvent instance")
	}

	if event.callerSkip != 2 {
		t.Errorf("Expected callerSkip to be 2, got %d", event.callerSkip)
	}

	// Test adding more caller skips
	result = event.AddCallerSkip(3)

	if event.callerSkip != 5 {
		t.Errorf("Expected callerSkip to be 5, got %d", event.callerSkip)
	}
}

func TestLogEvent_Field(t *testing.T) {
	event := newLogEvent(zapcore.InfoLevel)
	result := event.Field("test_key", "test_value")

	if result != event {
		t.Error("Expected Field to return the same LogEvent instance")
	}

	// Check that the field was added to data
	found := false
	for _, field := range event.data {
		if field.Key == "test_key" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected test_key field to be added to data")
	}

	// Check that the field was added to raw.data
	if event.raw.data["test_key"] != "test_value" {
		t.Errorf("Expected raw.data[test_key] to be test_value, got %v", event.raw.data["test_key"])
	}
}

func TestLogEvent_Pretty(t *testing.T) {
	event := newLogEvent(zapcore.InfoLevel)
	result := event.Pretty()

	if result != event {
		t.Error("Expected Pretty to return the same LogEvent instance")
	}

	if !event.pretty {
		t.Error("Expected pretty to be true")
	}
}

// TestLogEvent_Msg is difficult to test thoroughly without mocking the logger
// This is a simple test to ensure it doesn't panic
func TestLogEvent_Msg(t *testing.T) {
	// Test with different modes and pretty settings
	tests := []struct {
		name   string
		mode   string
		pretty bool
	}{
		{
			name:   "development mode, not pretty",
			mode:   development,
			pretty: false,
		},
		{
			name:   "development mode, pretty",
			mode:   development,
			pretty: true,
		},
		{
			name:   "production mode, not pretty",
			mode:   production,
			pretty: false,
		},
		{
			name:   "pretty mode, not pretty",
			mode:   pretty,
			pretty: false,
		},
	}

	// Save original mode and restore it after the test
	originalMode := mode
	defer func() { mode = originalMode }()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set mode for this test
			mode = tc.mode

			event := newLogEvent(zapcore.InfoLevel)
			if tc.pretty {
				event.Pretty()
			}

			// Add some fields
			event.Field("test_key", "test_value")

			// This should not panic
			event.Msg("test message")
		})
	}
}
