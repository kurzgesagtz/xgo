package xtype

import (
	"github.com/gotidy/ptr"
	"testing"
	"time"
)

func TestDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		expectErr bool
		expected  *time.Time
	}{
		{
			name:      "valid RFC3339 date",
			input:     []byte(`"2023-11-25T15:04:05Z"`),
			expectErr: false,
			expected:  ptr.Time(time.Date(2023, 11, 25, 15, 4, 5, 0, time.UTC)),
		},
		{
			name:      "valid simple date",
			input:     []byte(`"2023-11-25"`),
			expectErr: false,
			expected:  ptr.Time(time.Date(2023, 11, 25, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:      "invalid date format",
			input:     []byte(`"11-25-2023"`),
			expectErr: true,
			expected:  nil,
		},
		{
			name:      "empty input",
			input:     []byte(`""`),
			expectErr: true,
			expected:  nil,
		},
		{
			name:      "non string input",
			input:     []byte(`12345`),
			expectErr: true,
			expected:  nil,
		},
		{
			name:      "nil input",
			input:     nil,
			expectErr: true,
			expected:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var d Date
			err := d.UnmarshalJSON(tc.input)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected xerror, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected xerror: %v", err)
				}
				if tc.expected != nil && !time.Time(d).Equal(*tc.expected) {
					t.Errorf("expected %v, got %v", *tc.expected, time.Time(d))
				}
			}
		})
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     Date
		expected  string
		expectErr bool
	}{
		{
			name:     "valid RFC3339 date",
			input:    Date(time.Date(2023, 11, 25, 15, 4, 5, 0, time.UTC)),
			expected: `"2023-11-25T15:04:05Z"`,
		},
		{
			name:     "zero time",
			input:    Date(time.Time{}),
			expected: `"0001-01-01T00:00:00Z"`,
		},
		{
			name:     "far future date",
			input:    Date(time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)),
			expected: `"3000-01-01T00:00:00Z"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := tc.input.MarshalJSON()
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected xerror, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected xerror: %v", err)
				}
				if string(output) != tc.expected {
					t.Errorf("expected %s, got %s", tc.expected, string(output))
				}
			}
		})
	}
}
