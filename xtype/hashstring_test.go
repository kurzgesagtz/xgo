package xtype

import (
	"context"
	"reflect"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/schema"
)

func TestNewHashString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		isHashed bool
	}{
		{
			name:     "plain text",
			input:    "password123",
			isHashed: false,
		},
		{
			name:     "empty string",
			input:    "",
			isHashed: false,
		},
		{
			name: "hashed string",
			input: func() string {
				hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), DefaultBcryptCost)
				return string(hash)
			}(),
			isHashed: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hs := NewHashString(tc.input)

			// Check if the string is preserved
			if hs.str != tc.input {
				t.Errorf("expected string %q, got %q", tc.input, hs.str)
			}

			// Check if hash status is correctly detected
			if hs.hash != tc.isHashed {
				t.Errorf("expected hash status %v, got %v", tc.isHashed, hs.hash)
			}

			// If it's a hashed string, check if cost is correctly detected
			if tc.isHashed {
				expectedCost, err := bcrypt.Cost([]byte(tc.input))
				if err != nil {
					t.Fatalf("failed to get cost: %v", err)
				}
				if hs.cost != expectedCost {
					t.Errorf("expected cost %d, got %d", expectedCost, hs.cost)
				}
			}
		})
	}
}

func TestHashString_Scan(t *testing.T) {
	// Create a hashed password for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), DefaultBcryptCost)
	hashedPasswordStr := string(hashedPassword)

	tests := []struct {
		name      string
		value     interface{}
		expectErr bool
		isHashed  bool
	}{
		{
			name:      "string value (hashed)",
			value:     hashedPasswordStr,
			expectErr: false,
			isHashed:  true,
		},
		{
			name:      "byte slice (hashed)",
			value:     hashedPassword,
			expectErr: false,
			isHashed:  true,
		},
		{
			name:      "string value (plain)",
			value:     "password123",
			expectErr: true, // The Scan method will return an error for plain text passwords
			isHashed:  true, // The Scan method always sets hash: true for non-nil values
		},
		{
			name:      "nil value",
			value:     nil,
			expectErr: false,
			isHashed:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var hs HashString
			ctx := context.Background()
			field := &schema.Field{}
			dst := reflect.ValueOf(&hs).Elem()

			err := hs.Scan(ctx, field, dst, tc.value)

			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if tc.value == nil {
					if hs.str != "" {
						t.Errorf("expected empty string, got %q", hs.str)
					}
					if hs.hash != false {
						t.Errorf("expected hash status false, got true")
					}
				} else {
					// For non-nil values, check if hash status is correctly detected
					if hs.hash != tc.isHashed {
						t.Errorf("expected hash status %v, got %v", tc.isHashed, hs.hash)
					}

					// Check if the string is preserved
					var expectedStr string
					switch v := tc.value.(type) {
					case string:
						expectedStr = v
					case []byte:
						expectedStr = string(v)
					}

					if hs.str != expectedStr {
						t.Errorf("expected string %q, got %q", expectedStr, hs.str)
					}
				}
			}
		})
	}
}

func TestHashString_Value(t *testing.T) {
	// Create a hashed password for testing
	plainPassword := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), DefaultBcryptCost)
	hashedPasswordStr := string(hashedPassword)

	tests := []struct {
		name      string
		hs        HashString
		expectNil bool
		expectHashed bool
	}{
		{
			name: "plain text",
			hs: HashString{
				hash: false,
				str:  plainPassword,
			},
			expectNil: false,
			expectHashed: true,
		},
		{
			name: "already hashed",
			hs: HashString{
				hash: true,
				str:  hashedPasswordStr,
			},
			expectNil: false,
			expectHashed: true,
		},
		{
			name: "empty string",
			hs: HashString{
				hash: false,
				str:  "",
			},
			expectNil: true,
			expectHashed: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			field := &schema.Field{}
			dst := reflect.ValueOf(&tc.hs).Elem()

			val, err := tc.hs.Value(ctx, field, dst, tc.hs)

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

				// If we expect a hashed value, verify it's a valid bcrypt hash
				if tc.expectHashed {
					strVal, ok := val.(string)
					if !ok {
						t.Errorf("expected string, got %T", val)
					} else {
						// For plain text input, it should now be hashed
						if !tc.hs.hash {
							// Should be different from the original
							if strVal == tc.hs.str {
								t.Errorf("expected hashed value, got original string")
							}

							// Should be a valid bcrypt hash
							err := bcrypt.CompareHashAndPassword([]byte(strVal), []byte(tc.hs.str))
							if err != nil {
								t.Errorf("expected valid bcrypt hash, got error: %v", err)
							}
						} else {
							// For already hashed input, it should be the same
							if strVal != tc.hs.str {
								t.Errorf("expected same hash, got different value")
							}
						}
					}
				}
			}
		})
	}
}

func TestHashString_Equal(t *testing.T) {
	plainPassword := "password123"
	wrongPassword := "wrongpassword"

	// Create a hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), DefaultBcryptCost)
	hashedPasswordStr := string(hashedPassword)

	tests := []struct {
		name     string
		hs       *HashString
		input    string
		expected bool
	}{
		{
			name: "plain text - correct password",
			hs: &HashString{
				hash: false,
				str:  plainPassword,
			},
			input:    plainPassword,
			expected: true,
		},
		{
			name: "plain text - wrong password",
			hs: &HashString{
				hash: false,
				str:  plainPassword,
			},
			input:    wrongPassword,
			expected: false,
		},
		{
			name: "hashed - correct password",
			hs: &HashString{
				hash: true,
				str:  hashedPasswordStr,
			},
			input:    plainPassword,
			expected: true,
		},
		{
			name: "hashed - wrong password",
			hs: &HashString{
				hash: true,
				str:  hashedPasswordStr,
			},
			input:    wrongPassword,
			expected: false,
		},
		{
			name:     "nil HashString",
			hs:       nil,
			input:    plainPassword,
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.hs.Equal(tc.input)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestHashString_String(t *testing.T) {
	plainPassword := "password123"

	tests := []struct {
		name     string
		hs       *HashString
		expected string
	}{
		{
			name: "plain text",
			hs: &HashString{
				hash: false,
				str:  plainPassword,
			},
			expected: plainPassword,
		},
		{
			name: "empty string",
			hs: &HashString{
				hash: false,
				str:  "",
			},
			expected: "",
		},
		{
			name:     "nil HashString",
			hs:       nil,
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.hs.String()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}
