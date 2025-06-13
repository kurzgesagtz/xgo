package xerror

import (
	"encoding/json"
	"testing"
	"time"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name: "with code and message",
			err: &Error{
				Code:    "TEST_CODE",
				Message: "Test message",
			},
			expected: "TEST_CODE: Test message",
		},
		{
			name: "with code only",
			err: &Error{
				Code: "TEST_CODE",
			},
			expected: "TEST_CODE: ",
		},
		{
			name: "with message only",
			err: &Error{
				Message: "Test message",
			},
			expected: ": Test message",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.err.Error()
			if result != tc.expected {
				t.Errorf("Expected error string %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestError_StackTrace(t *testing.T) {
	// Test with nil callers
	err := &Error{
		Code:    "TEST_CODE",
		Message: "Test message",
		Callers: nil,
	}

	st := err.StackTrace()
	if len(st) != 0 {
		t.Errorf("Expected empty stack trace for nil callers, got length %d", len(st))
	}

	// Test with non-nil callers
	cs := callers(1)
	err = &Error{
		Code:    "TEST_CODE",
		Message: "Test message",
		Callers: cs,
	}

	st = err.StackTrace()
	if len(st) != len(*cs) {
		t.Errorf("Expected stack trace length %d, got %d", len(*cs), len(st))
	}
}

func TestTimestamp_MarshalJSON(t *testing.T) {
	// Create a fixed time for testing
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	ts := Timestamp(fixedTime)

	// Marshal to JSON (using pointer to Timestamp since MarshalJSON is defined on pointer receiver)
	bytes, err := json.Marshal(&ts)
	if err != nil {
		t.Fatalf("Failed to marshal timestamp: %v", err)
	}

	// Convert bytes to string and compare with expected
	actualStr := string(bytes)
	expectedStr := "1672574400000" // 2023-01-01 12:00:00 UTC in milliseconds

	if actualStr != expectedStr {
		t.Errorf("Expected marshaled timestamp %s, got %s", expectedStr, actualStr)
	}
}

func TestTimestamp_UnmarshalJSON(t *testing.T) {
	// Test valid timestamp
	jsonData := []byte("1672574400000") // 2023-01-01 12:00:00 UTC in milliseconds
	var ts Timestamp

	err := json.Unmarshal(jsonData, &ts)
	if err != nil {
		t.Fatalf("Failed to unmarshal timestamp: %v", err)
	}

	// Convert to time.Time and check
	tm := time.Time(ts)
	expected := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	if !tm.Equal(expected) {
		t.Errorf("Expected time %v, got %v", expected, tm)
	}

	// Test invalid timestamp
	jsonData = []byte("\"not-a-number\"")
	err = json.Unmarshal(jsonData, &ts)
	if err == nil {
		t.Error("Expected error when unmarshaling invalid timestamp, got nil")
	}
}

func TestTimestampNow(t *testing.T) {
	before := time.Now().UTC()
	ts := TimestampNow()
	after := time.Now().UTC()

	// Convert Timestamp to time.Time
	tm := time.Time(ts)

	// Check that the timestamp is between before and after
	if tm.Before(before) || tm.After(after) {
		t.Errorf("TimestampNow() returned %v, which is not between %v and %v", tm, before, after)
	}
}

func TestErrorConstants(t *testing.T) {
	// Test that error constants are defined
	constants := []string{
		ErrCodeUnauthorized,
		ErrCodePermissionDenied,
		ErrCodeInvalidRequest,
		ErrCodeNotFound,
		ErrCodeNotImplement,
		ErrCodeInternalError,
		ErrCodePartnerInternalError,
		ErrCodePartnerBadResponseError,
		ErrCodeServiceInternalError,
		ErrCodeServerPanic,
		ErrCodeInvalidEnum,
		ErrCodeClientRequestCanceled,
		ErrCodeClientRequestDeadlineExceed,
		ErrCodeAlreadyExists,
		ErrCodeRateLimitExceeded,
	}

	for _, c := range constants {
		if c == "" {
			t.Error("Expected non-empty error code constant")
		}
	}
}
