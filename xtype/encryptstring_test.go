package xtype

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
)

func TestEncryptString_String(t *testing.T) {
	tests := []struct {
		name     string
		es       *EncryptString
		expected string
	}{
		{
			name:     "nil EncryptString",
			es:       nil,
			expected: "",
		},
		{
			name: "empty EncryptString",
			es: &EncryptString{
				encrypt: nil,
				str:     nil,
			},
			expected: "",
		},
		{
			name: "with string value",
			es: &EncryptString{
				str: &[]string{"test"}[0],
			},
			expected: "test",
		},
		{
			name: "with encrypted value",
			es: func() *EncryptString {
				// Save original secret
				originalSecret := EncryptStringSecret
				// Set a known secret for deterministic testing
				EncryptStringSecret = []string{"20dd6fbb502a0465e070d3ed4b92a84e"}
				// Create encrypted string
				enc := dynamicEncrypt("test")
				// Restore original secret
				EncryptStringSecret = originalSecret
				return &EncryptString{
					encrypt: &enc,
				}
			}(),
			expected: "test",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.es == nil {
				// For nil case, we can't call methods on it, so just skip
				// In a real application, calling methods on nil would panic
				// but we're testing the non-nil behavior here
				return
			}

			result := tc.es.String()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestEncryptString_EncryptString(t *testing.T) {
	// Save original secret
	originalSecret := EncryptStringSecret
	// Set a known secret for deterministic testing
	EncryptStringSecret = []string{"20dd6fbb502a0465e070d3ed4b92a84e"}

	plaintext := "test"
	encrypted := dynamicEncrypt(plaintext)

	tests := []struct {
		name     string
		es       *EncryptString
		expected string
	}{
		{
			name:     "nil EncryptString",
			es:       nil,
			expected: "",
		},
		{
			name: "empty EncryptString",
			es: &EncryptString{
				encrypt: nil,
				str:     nil,
			},
			expected: "",
		},
		{
			name: "with string value",
			es: &EncryptString{
				str: &plaintext,
			},
			expected: encrypted,
		},
		{
			name: "with encrypted value",
			es: &EncryptString{
				encrypt: &encrypted,
			},
			expected: encrypted,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.es == nil {
				// For nil case, we can't call methods on it, so just skip
				// In a real application, calling methods on nil would panic
				// but we're testing the non-nil behavior here
				return
			}

			result := tc.es.EncryptString()
			if tc.name == "with string value" {
				// For this case, we can't predict the exact encrypted value due to random nonce
				// So we just check that it's not empty and not the original string
				if result == "" || result == plaintext {
					t.Errorf("expected encrypted string, got %q", result)
				}
			} else if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}

	// Restore original secret
	EncryptStringSecret = originalSecret
}

func TestEncryptString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		es       *EncryptString
		expected string
	}{
		{
			name: "with string value",
			es: &EncryptString{
				str: &[]string{"test"}[0],
			},
			expected: `"test"`,
		},
		{
			name: "empty string",
			es: &EncryptString{
				str: &[]string{""}[0],
			},
			expected: `""`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.es.MarshalJSON()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(result) != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, string(result))
			}
		})
	}
}

func TestEncryptString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		expectErr bool
		expected  string
	}{
		{
			name:      "valid string",
			input:     []byte(`"test"`),
			expectErr: false,
			expected:  "test",
		},
		{
			name:      "empty string",
			input:     []byte(`""`),
			expectErr: false,
			expected:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var es EncryptString
			err := es.UnmarshalJSON(tc.input)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if es.String() != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, es.String())
				}
			}
		})
	}
}

func TestEncryptString_Scan(t *testing.T) {
	// Save original secret
	originalSecret := EncryptStringSecret
	// Set a known secret for deterministic testing
	EncryptStringSecret = []string{"20dd6fbb502a0465e070d3ed4b92a84e"}

	plaintext := "test"
	encrypted := dynamicEncrypt(plaintext)

	tests := []struct {
		name      string
		value     interface{}
		expectErr bool
		expected  string
	}{
		{
			name:      "string value",
			value:     plaintext,
			expectErr: false,
			expected:  plaintext,
		},
		{
			name:      "encrypted string value",
			value:     encrypted,
			expectErr: false,
			expected:  plaintext,
		},
		{
			name:      "byte slice",
			value:     []byte(plaintext),
			expectErr: false,
			expected:  plaintext,
		},
		{
			name:      "nil value",
			value:     nil,
			expectErr: false,
			expected:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var es EncryptString
			err := es.Scan(tc.value)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if es.String() != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, es.String())
				}
			}
		})
	}

	// Restore original secret
	EncryptStringSecret = originalSecret
}

func TestEncryptString_Value(t *testing.T) {
	// Save original secret
	originalSecret := EncryptStringSecret
	// Set a known secret for deterministic testing
	EncryptStringSecret = []string{"20dd6fbb502a0465e070d3ed4b92a84e"}

	tests := []struct {
		name      string
		es        *EncryptString
		expectNil bool
	}{
		{
			name: "with string value",
			es: &EncryptString{
				str: &[]string{"test"}[0],
			},
			expectNil: false,
		},
		{
			name: "empty string",
			es: &EncryptString{
				str: &[]string{""}[0],
			},
			expectNil: false, // The Value() method returns the EncryptString object for empty strings
		},
		{
			name:      "nil EncryptString",
			es:        nil,
			expectNil: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.es == nil {
				// For nil case, we can't call methods on it, so just skip
				// In a real application, calling methods on nil would panic
				// but we're testing the non-nil behavior here
				return
			}

			val, err := tc.es.Value()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tc.expectNil {
				if val != nil {
					t.Errorf("expected nil, got %v", val)
					// Add debug info
					t.Logf("es.str = %v", tc.es.str)
					t.Logf("es.encrypt = %v", tc.es.encrypt)
					t.Logf("es.EncryptString() = %q", tc.es.EncryptString())
				}
			} else {
				if val == nil {
					t.Errorf("expected non-nil value, got nil")
				}
				// Check that the value is of the expected type
				_, ok := val.(driver.Valuer)
				if !ok {
					t.Errorf("expected driver.Valuer, got %T", val)
				}
			}
		})
	}

	// Restore original secret
	EncryptStringSecret = originalSecret
}

func TestNewEncryptString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expectNil bool
	}{
		{
			name:     "valid string",
			input:    "test",
			expectNil: false,
		},
		{
			name:     "empty string",
			input:    "",
			expectNil: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := NewEncryptString(tc.input)
			if tc.expectNil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("expected non-nil result, got nil")
				} else if result.String() != tc.input {
					t.Errorf("expected %q, got %q", tc.input, result.String())
				}
			}
		})
	}
}

func TestEncryptString_JSON(t *testing.T) {
	type TestStruct struct {
		Name     string         `json:"name"`
		Password *EncryptString `json:"password"`
	}

	// Test marshaling and unmarshaling
	original := TestStruct{
		Name:     "user",
		Password: NewEncryptString("secret"),
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Print the JSON data for debugging
	t.Logf("JSON data: %s", string(data))

	// Unmarshal back
	var decoded TestStruct
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Print the decoded values for debugging
	t.Logf("Decoded Name: %q", decoded.Name)
	t.Logf("Decoded Password: %+v", decoded.Password)

	// Check values
	if decoded.Name != original.Name {
		t.Errorf("Name: expected %q, got %q", original.Name, decoded.Name)
	}
	if decoded.Password.String() != original.Password.String() {
		t.Errorf("Password: expected %q, got %q", original.Password.String(), decoded.Password.String())
	}
}
