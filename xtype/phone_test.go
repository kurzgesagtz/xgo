package xtype

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/phonenumbers"
)

func TestNewPhone(t *testing.T) {
	// Save original default region
	originalRegion := PhoneNumberDefaultRegion
	// Set a known region for testing
	PhoneNumberDefaultRegion = "TH"

	tests := []struct {
		name      string
		phone     string
		code      string
		expectErr bool
	}{
		{
			name:      "valid phone with country code",
			phone:     "+66812345678",
			code:      "TH",
			expectErr: false,
		},
		{
			name:      "valid phone without country code",
			phone:     "0812345678",
			code:      "TH",
			expectErr: false,
		},
		{
			name:      "invalid phone",
			phone:     "invalid",
			code:      "TH",
			expectErr: true,
		},
		{
			name:      "empty phone",
			phone:     "",
			code:      "TH",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p, err := NewPhone(tc.phone, tc.code)

			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if p == nil {
					t.Errorf("expected non-nil Phone, got nil")
				} else {
					// Verify that the phone number is valid
					if !phonenumbers.IsValidNumber(p.phone) {
						t.Errorf("expected valid phone number, got invalid")
					}
				}
			}
		})
	}

	// Restore original default region
	PhoneNumberDefaultRegion = originalRegion
}

func TestPhone_Scan(t *testing.T) {
	// Save original default region
	originalRegion := PhoneNumberDefaultRegion
	// Set a known region for testing
	PhoneNumberDefaultRegion = "TH"

	validPhone := "+66812345678"

	tests := []struct {
		name      string
		value     interface{}
		expectErr bool
	}{
		{
			name:      "string value",
			value:     validPhone,
			expectErr: false,
		},
		{
			name:      "byte slice",
			value:     []byte(validPhone),
			expectErr: false,
		},
		{
			name:      "invalid phone",
			value:     "invalid",
			expectErr: true,
		},
		{
			name:      "nil value",
			value:     nil,
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var p Phone
			err := p.Scan(tc.value)

			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if tc.value == nil {
					// For nil input, phone should be nil
					if p.phone != nil {
						t.Errorf("expected nil phone, got %v", p.phone)
					}
				} else if !tc.expectErr {
					// For valid input, phone should be non-nil and valid
					if p.phone == nil {
						t.Errorf("expected non-nil phone, got nil")
					} else if !phonenumbers.IsValidNumber(p.phone) {
						t.Errorf("expected valid phone number, got invalid")
					}
				}
			}
		})
	}

	// Restore original default region
	PhoneNumberDefaultRegion = originalRegion
}

func TestPhone_Value(t *testing.T) {
	// Save original default region
	originalRegion := PhoneNumberDefaultRegion
	// Set a known region for testing
	PhoneNumberDefaultRegion = "TH"

	validPhone, _ := NewPhone("+66812345678", "TH")

	tests := []struct {
		name      string
		phone     *Phone
		expectNil bool
		expected  string
	}{
		{
			name:      "valid phone",
			phone:     validPhone,
			expectNil: false,
			expected:  "+66812345678",
		},
		{
			name:      "nil phone",
			phone:     nil,
			expectNil: true,
			expected:  "",
		},
		{
			name:      "phone with nil internal value",
			phone:     &Phone{phone: nil},
			expectNil: true,
			expected:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.phone.Value()

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tc.expectNil {
				if val != nil {
					t.Errorf("expected nil, got %v", val)
				}
			} else {
				if val == nil {
					t.Errorf("expected non-nil value, got nil")
				}

				// Check that the value is the expected string
				strVal, ok := val.(string)
				if !ok {
					t.Errorf("expected string, got %T", val)
				} else if strVal != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, strVal)
				}
			}
		})
	}

	// Restore original default region
	PhoneNumberDefaultRegion = originalRegion
}

func TestPhone_String(t *testing.T) {
	// Save original default region
	originalRegion := PhoneNumberDefaultRegion
	// Set a known region for testing
	PhoneNumberDefaultRegion = "TH"

	validPhone, _ := NewPhone("+66812345678", "TH")

	result := validPhone.String()
	expected := "+66812345678"

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}

	// Restore original default region
	PhoneNumberDefaultRegion = originalRegion
}

func TestPhone_FormatMethods(t *testing.T) {
	// Save original default region
	originalRegion := PhoneNumberDefaultRegion
	// Set a known region for testing
	PhoneNumberDefaultRegion = "TH"

	validPhone, _ := NewPhone("+66812345678", "TH")

	tests := []struct {
		name     string
		method   func() string
		expected string
	}{
		{
			name:     "FormatE164",
			method:   validPhone.FormatE164,
			expected: "+66812345678",
		},
		{
			name:     "FormatInternational",
			method:   validPhone.FormatInternational,
			expected: "+66 81 234 5678",
		},
		{
			name:     "FormatNational",
			method:   validPhone.FormatNational,
			expected: "081 234 5678",
		},
		{
			name:     "FormatRFC3966",
			method:   validPhone.FormatRFC3966,
			expected: "tel:+66-81-234-5678",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.method()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}

	// Restore original default region
	PhoneNumberDefaultRegion = originalRegion
}

func TestPhone_MarshalJSON(t *testing.T) {
	// Save original default region
	originalRegion := PhoneNumberDefaultRegion
	// Set a known region for testing
	PhoneNumberDefaultRegion = "TH"

	validPhone, _ := NewPhone("+66812345678", "TH")

	data, err := validPhone.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := `"+66812345678"`
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}

	// Restore original default region
	PhoneNumberDefaultRegion = originalRegion
}

func TestPhone_UnmarshalJSON(t *testing.T) {
	// Save original default region
	originalRegion := PhoneNumberDefaultRegion
	// Set a known region for testing
	PhoneNumberDefaultRegion = "TH"

	tests := []struct {
		name      string
		input     []byte
		expectErr bool
		expected  string
	}{
		{
			name:      "valid phone",
			input:     []byte(`"+66812345678"`),
			expectErr: false,
			expected:  "+66812345678",
		},
		{
			name:      "invalid phone",
			input:     []byte(`"invalid"`),
			expectErr: true,
			expected:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var p Phone
			err := p.UnmarshalJSON(tc.input)

			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				result := p.String()
				if result != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, result)
				}
			}
		})
	}

	// Restore original default region
	PhoneNumberDefaultRegion = originalRegion
}

func TestPhone_JSON(t *testing.T) {
	// Save original default region
	originalRegion := PhoneNumberDefaultRegion
	// Set a known region for testing
	PhoneNumberDefaultRegion = "TH"

	type TestStruct struct {
		Name  string `json:"name"`
		Phone *Phone `json:"phone"`
	}

	// Create a valid phone
	validPhone, _ := NewPhone("+66812345678", "TH")

	// Create a test struct
	original := TestStruct{
		Name:  "John Doe",
		Phone: validPhone,
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal back
	var decoded TestStruct
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Check values
	if decoded.Name != original.Name {
		t.Errorf("Name: expected %q, got %q", original.Name, decoded.Name)
	}

	if decoded.Phone.String() != original.Phone.String() {
		t.Errorf("Phone: expected %q, got %q", original.Phone.String(), decoded.Phone.String())
	}

	// Restore original default region
	PhoneNumberDefaultRegion = originalRegion
}
