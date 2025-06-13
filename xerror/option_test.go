package xerror

import (
	"testing"
	"time"
)

func TestWithCaller(t *testing.T) {
	// Create a base error
	baseErr := &Error{
		Code: "TEST_CODE",
	}
	
	// Apply WithCaller option
	optionFunc := WithCaller(1)
	resultErr := optionFunc(baseErr)
	
	// Check that the caller was set
	if resultErr.Caller == "" {
		t.Error("Expected non-empty caller")
	}
	
	// Check that the function returns the same error instance
	if resultErr != baseErr {
		t.Error("Expected WithCaller to return the same error instance")
	}
}

func TestWithDetail(t *testing.T) {
	// Create a base error
	baseErr := &Error{
		Code: "TEST_CODE",
	}
	
	// Apply WithDetail option
	detail := "Detailed error information"
	optionFunc := WithDetail(detail)
	resultErr := optionFunc(baseErr)
	
	// Check that the detail was set correctly
	if resultErr.Detail != detail {
		t.Errorf("Expected detail %q, got %q", detail, resultErr.Detail)
	}
	
	// Check that the function returns the same error instance
	if resultErr != baseErr {
		t.Error("Expected WithDetail to return the same error instance")
	}
}

func TestWithMessage(t *testing.T) {
	// Create a base error
	baseErr := &Error{
		Code: "TEST_CODE",
	}
	
	// Apply WithMessage option
	message := "Error message"
	optionFunc := WithMessage(message)
	resultErr := optionFunc(baseErr)
	
	// Check that the message was set correctly
	if resultErr.Message != message {
		t.Errorf("Expected message %q, got %q", message, resultErr.Message)
	}
	
	// Check that the function returns the same error instance
	if resultErr != baseErr {
		t.Error("Expected WithMessage to return the same error instance")
	}
}

func TestWithTimestamp(t *testing.T) {
	// Create a base error
	baseErr := &Error{
		Code: "TEST_CODE",
	}
	
	// Apply WithTimestamp option with a fixed time
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	optionFunc := WithTimestamp(fixedTime)
	resultErr := optionFunc(baseErr)
	
	// Check that the timestamp was set correctly
	expectedTimestamp := Timestamp(fixedTime.UTC())
	actualTimestamp := resultErr.Timestamp
	
	// Convert both to time.Time for comparison
	expectedTime := time.Time(expectedTimestamp)
	actualTime := time.Time(actualTimestamp)
	
	if !expectedTime.Equal(actualTime) {
		t.Errorf("Expected timestamp %v, got %v", expectedTime, actualTime)
	}
	
	// Check that the function returns the same error instance
	if resultErr != baseErr {
		t.Error("Expected WithTimestamp to return the same error instance")
	}
}

func TestWithJSONInfo(t *testing.T) {
	// Create a base error
	baseErr := &Error{
		Code: "TEST_CODE",
	}
	
	// Apply WithJSONInfo option
	key := "test_key"
	value := "test_value"
	optionFunc := WithJSONInfo(key, value)
	resultErr := optionFunc(baseErr)
	
	// Check that the info map was created
	if resultErr.Info == nil {
		t.Error("Expected non-nil Info map")
	}
	
	// Check that the key-value pair was added correctly
	if resultErr.Info[key] != value {
		t.Errorf("Expected Info[%q] to be %q, got %v", key, value, resultErr.Info[key])
	}
	
	// Check that the function returns the same error instance
	if resultErr != baseErr {
		t.Error("Expected WithJSONInfo to return the same error instance")
	}
	
	// Test adding a second key-value pair
	key2 := "test_key2"
	value2 := 42
	optionFunc2 := WithJSONInfo(key2, value2)
	resultErr2 := optionFunc2(resultErr)
	
	// Check that both key-value pairs exist
	if resultErr2.Info[key] != value {
		t.Errorf("Expected Info[%q] to be %q, got %v", key, value, resultErr2.Info[key])
	}
	
	if resultErr2.Info[key2] != value2 {
		t.Errorf("Expected Info[%q] to be %d, got %v", key2, value2, resultErr2.Info[key2])
	}
}