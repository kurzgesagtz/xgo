package xlog

import (
	"os"
	"testing"
)

func TestLogLevels(t *testing.T) {
	// Test that all log level functions return a non-nil LogEvent
	tests := []struct {
		name     string
		logFunc  func() *LogEvent
	}{
		{
			name:    "Debug",
			logFunc: Debug,
		},
		{
			name:    "Info",
			logFunc: Info,
		},
		{
			name:    "Warn",
			logFunc: Warn,
		},
		{
			name:    "Error",
			logFunc: Error,
		},
		{
			name:    "Panic",
			logFunc: Panic,
		},
		{
			name:    "Fatal",
			logFunc: Fatal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			event := tc.logFunc()
			if event == nil {
				t.Errorf("Expected non-nil LogEvent from %s()", tc.name)
			}
		})
	}
}

func TestLogModeEnvironmentVariable(t *testing.T) {
	// This test is limited because we can't directly call init()
	// We can only verify the current mode based on environment variables

	// Save original mode and restore it after the test
	originalMode := mode
	defer func() { mode = originalMode }()

	// Test the current mode based on environment variables
	logMode, exists := os.LookupEnv(logProductionKey)
	if !exists {
		// If environment variable is not set, mode should be development
		if mode != development {
			t.Errorf("Expected mode to be %s when %s is not set, got %s", 
				development, logProductionKey, mode)
		}
	} else if logMode == production {
		// If environment variable is set to production, mode should be production
		if mode != production {
			t.Errorf("Expected mode to be %s when %s=%s, got %s", 
				production, logProductionKey, production, mode)
		}
	} else if logMode == pretty {
		// If environment variable is set to pretty, mode should be pretty
		if mode != pretty {
			t.Errorf("Expected mode to be %s when %s=%s, got %s", 
				pretty, logProductionKey, pretty, mode)
		}
	}
}

func TestAppNameEnvironmentVariable(t *testing.T) {
	// This test is limited because we can't directly call init()
	// We can only verify the current appName based on environment variables

	// Save original app name and restore it after the test
	originalAppName := appName
	defer func() { appName = originalAppName }()

	// Test the current appName based on environment variables
	envAppName, exists := os.LookupEnv(logAppNameKey)
	if !exists {
		// If environment variable is not set, appName should be "local"
		if appName != "local" {
			t.Errorf("Expected appName to be %s when %s is not set, got %s", 
				"local", logAppNameKey, appName)
		}
	} else {
		// If environment variable is set, appName should match it
		if appName != envAppName {
			t.Errorf("Expected appName to be %s when %s=%s, got %s", 
				envAppName, logAppNameKey, envAppName, appName)
		}
	}
}

func TestPrettyPrint(t *testing.T) {
	// This is a simple test to ensure PrettyPrint doesn't panic
	// It's difficult to test the actual output without capturing stdout

	// Test with a simple object
	obj := map[string]string{"key": "value"}

	// This should not panic
	PrettyPrint(obj)

	// Test with multiple objects
	obj2 := []int{1, 2, 3}

	// This should not panic
	PrettyPrint(obj, obj2)
}
