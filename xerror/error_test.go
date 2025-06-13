package xerror

import (
	"errors"
	"os"
	"testing"
)

func TestNewError(t *testing.T) {
	// Save original app name and restore it after the test
	originalAppName := appName
	defer func() { appName = originalAppName }()

	// Test with default app name
	appName = "test-app"
	err := NewError("TEST_CODE")

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	if err.Code != "TEST_CODE" {
		t.Errorf("Expected code TEST_CODE, got %s", err.Code)
	}

	if err.AppName != "test-app" {
		t.Errorf("Expected app name test-app, got %s", err.AppName)
	}

	if err.Caller == "" {
		t.Error("Expected non-empty caller")
	}

	if err.Callers == nil {
		t.Error("Expected non-nil callers")
	}

	// Test with options
	err = NewError("TEST_CODE", WithMessage("Test message"))

	if err.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got %s", err.Message)
	}
}

func TestNewErrorWithEnv(t *testing.T) {
	// Save original app name and restore it after the test
	originalAppName := appName
	defer func() { appName = originalAppName }()

	// Set environment variable and update appName directly
	os.Setenv(xErrorAppNameKey, "env-app")
	appName = "env-app"

	err := NewError("TEST_CODE")

	if err.AppName != "env-app" {
		t.Errorf("Expected app name env-app, got %s", err.AppName)
	}

	// Clean up
	os.Unsetenv(xErrorAppNameKey)
}

func TestIsErrorCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		code     string
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			code:     "TEST_CODE",
			expected: false,
		},
		{
			name:     "matching code",
			err:      NewError("TEST_CODE"),
			code:     "TEST_CODE",
			expected: true,
		},
		{
			name:     "non-matching code",
			err:      NewError("TEST_CODE"),
			code:     "OTHER_CODE",
			expected: false,
		},
		{
			name:     "non-xerror error with internal error code",
			err:      errors.New("standard error"),
			code:     ErrCodeInternalError,
			expected: true,
		},
		{
			name:     "non-xerror error with other code",
			err:      errors.New("standard error"),
			code:     "OTHER_CODE",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsErrorCode(tc.err, tc.code)
			if result != tc.expected {
				t.Errorf("Expected IsErrorCode to return %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestGetErrorMessage(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: "",
		},
		{
			name:     "xerror with message",
			err:      NewError("TEST_CODE", WithMessage("Test message")),
			expected: "Test message",
		},
		{
			name:     "standard error",
			err:      errors.New("standard error"),
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GetErrorMessage(tc.err)
			if result != tc.expected {
				t.Errorf("Expected message %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestCaller(t *testing.T) {
	result := caller(1)
	if result == "" {
		t.Error("Expected non-empty caller string")
	}
}

func TestCallers(t *testing.T) {
	result := callers(1)
	if result == nil {
		t.Error("Expected non-nil callers stack")
	}

	if len(*result) == 0 {
		t.Error("Expected non-empty callers stack")
	}
}
