package xutil

import (
	"context"
	"errors"
	"github.com/kurzgesagtz/xgo/xerror"
	"testing"
	"time"
)

func TestEnumToValue(t *testing.T) {
	// Define a test enum map
	testEnum := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	tests := []struct {
		name      string
		enum      map[string]int
		input     string
		expected  int
		expectErr bool
	}{
		{
			name:      "valid enum key",
			enum:      testEnum,
			input:     "one",
			expected:  1,
			expectErr: false,
		},
		{
			name:      "another valid enum key",
			enum:      testEnum,
			input:     "three",
			expected:  3,
			expectErr: false,
		},
		{
			name:      "invalid enum key",
			enum:      testEnum,
			input:     "four",
			expected:  0, // zero value for int
			expectErr: true,
		},
		{
			name:      "empty enum key",
			enum:      testEnum,
			input:     "",
			expected:  0,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := EnumToValue(tc.enum, tc.input)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if !xerror.IsErrorCode(err, xerror.ErrCodeInvalidEnum) {
					t.Errorf("expected error code %s, got %v", xerror.ErrCodeInvalidEnum, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tc.expected {
					t.Errorf("expected %v, got %v", tc.expected, result)
				}
			}
		})
	}
}

func TestValueToEnum(t *testing.T) {
	// Define a test enum map
	testEnum := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	tests := []struct {
		name      string
		enum      map[string]int
		input     int
		expected  string
		expectErr bool
	}{
		{
			name:      "valid enum value",
			enum:      testEnum,
			input:     1,
			expected:  "one",
			expectErr: false,
		},
		{
			name:      "another valid enum value",
			enum:      testEnum,
			input:     3,
			expected:  "three",
			expectErr: false,
		},
		{
			name:      "invalid enum value",
			enum:      testEnum,
			input:     4,
			expected:  "", // zero value for string
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ValueToEnum(tc.enum, tc.input)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if !xerror.IsErrorCode(err, xerror.ErrCodeInvalidEnum) {
					t.Errorf("expected error code %s, got %v", xerror.ErrCodeInvalidEnum, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tc.expected {
					t.Errorf("expected %v, got %v", tc.expected, result)
				}
			}
		})
	}
}

func TestMapToSlice(t *testing.T) {
	// Define a test mapper function
	doubleMapper := func(i int) (int, error) {
		return i * 2, nil
	}

	errorMapper := func(i int) (int, error) {
		if i == 3 {
			return 0, errors.New("test error")
		}
		return i * 2, nil
	}

	tests := []struct {
		name      string
		mapper    func(int) (int, error)
		input     []int
		expected  []int
		expectErr bool
	}{
		{
			name:      "valid mapping",
			mapper:    doubleMapper,
			input:     []int{1, 2, 3, 4, 5},
			expected:  []int{2, 4, 6, 8, 10},
			expectErr: false,
		},
		{
			name:      "empty input",
			mapper:    doubleMapper,
			input:     []int{},
			expected:  []int{},
			expectErr: false,
		},
		{
			name:      "mapper returns error",
			mapper:    errorMapper,
			input:     []int{1, 2, 3, 4, 5},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MapToSlice(tc.mapper, tc.input)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(result) != len(tc.expected) {
					t.Errorf("expected result length %d, got %d", len(tc.expected), len(result))
				} else {
					for i, v := range result {
						if v != tc.expected[i] {
							t.Errorf("at index %d: expected %v, got %v", i, tc.expected[i], v)
						}
					}
				}
			}
		})
	}
}

func TestMapToSliceWithOption(t *testing.T) {
	// Define a test mapper function with option
	multiplyMapper := func(i int, factors ...int) (int, error) {
		if len(factors) == 0 {
			return i, nil
		}
		result := i
		for _, factor := range factors {
			result *= factor
		}
		return result, nil
	}

	errorMapper := func(i int, factors ...int) (int, error) {
		if i == 3 {
			return 0, errors.New("test error")
		}
		result := i
		for _, factor := range factors {
			result *= factor
		}
		return result, nil
	}

	tests := []struct {
		name      string
		mapper    func(int, ...int) (int, error)
		input     []int
		options   []int
		expected  []int
		expectErr bool
	}{
		{
			name:      "mapping with single option",
			mapper:    multiplyMapper,
			input:     []int{1, 2, 3, 4, 5},
			options:   []int{2},
			expected:  []int{2, 4, 6, 8, 10},
			expectErr: false,
		},
		{
			name:      "mapping with multiple options",
			mapper:    multiplyMapper,
			input:     []int{1, 2, 3, 4, 5},
			options:   []int{2, 3},
			expected:  []int{6, 12, 18, 24, 30},
			expectErr: false,
		},
		{
			name:      "mapping with no options",
			mapper:    multiplyMapper,
			input:     []int{1, 2, 3, 4, 5},
			options:   []int{},
			expected:  []int{1, 2, 3, 4, 5},
			expectErr: false,
		},
		{
			name:      "empty input",
			mapper:    multiplyMapper,
			input:     []int{},
			options:   []int{2},
			expected:  []int{},
			expectErr: false,
		},
		{
			name:      "mapper returns error",
			mapper:    errorMapper,
			input:     []int{1, 2, 3, 4, 5},
			options:   []int{2},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MapToSliceWithOption(tc.mapper, tc.input, tc.options...)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(result) != len(tc.expected) {
					t.Errorf("expected result length %d, got %d", len(tc.expected), len(result))
				} else {
					for i, v := range result {
						if v != tc.expected[i] {
							t.Errorf("at index %d: expected %v, got %v", i, tc.expected[i], v)
						}
					}
				}
			}
		})
	}
}

func TestMapToSliceAsync(t *testing.T) {
	// Define a test processor function
	doubleProcessor := func(ctx context.Context, i int) (int, error) {
		return i * 2, nil
	}

	errorProcessor := func(ctx context.Context, i int) (int, error) {
		if i == 3 {
			return 0, errors.New("test error")
		}
		return i * 2, nil
	}

	delayedProcessor := func(ctx context.Context, i int) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(10 * time.Millisecond):
			return i * 2, nil
		}
	}

	tests := []struct {
		name        string
		processor   func(context.Context, int) (int, error)
		input       []int
		workers     int
		expected    []int
		expectErr   bool
		cancelCtx   bool
		description string
	}{
		{
			name:        "valid processing",
			processor:   doubleProcessor,
			input:       []int{1, 2, 3, 4, 5},
			workers:     2,
			expected:    []int{2, 4, 6, 8, 10},
			expectErr:   false,
			cancelCtx:   false,
			description: "Process a slice with multiple workers",
		},
		{
			name:        "empty input",
			processor:   doubleProcessor,
			input:       []int{},
			workers:     2,
			expected:    []int{},
			expectErr:   false,
			cancelCtx:   false,
			description: "Process an empty slice",
		},
		{
			name:        "processor returns error",
			processor:   errorProcessor,
			input:       []int{1, 2, 3, 4, 5},
			workers:     2,
			expected:    nil,
			expectErr:   true,
			cancelCtx:   false,
			description: "Handle error from processor",
		},
		{
			name:        "more workers than inputs",
			processor:   doubleProcessor,
			input:       []int{1, 2},
			workers:     5,
			expected:    []int{2, 4},
			expectErr:   false,
			cancelCtx:   false,
			description: "Handle case where workers > inputs",
		},
		{
			name:        "cancelled context",
			processor:   delayedProcessor,
			input:       []int{1, 2, 3, 4, 5},
			workers:     2,
			expected:    nil,
			expectErr:   true,
			cancelCtx:   true,
			description: "Handle cancelled context",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if tc.cancelCtx {
				// Cancel the context after a short delay
				go func() {
					time.Sleep(5 * time.Millisecond)
					cancel()
				}()
			}

			result, err := MapToSliceAsync(ctx, tc.workers, tc.processor, tc.input)
			
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(result) != len(tc.expected) {
					t.Errorf("expected result length %d, got %d", len(tc.expected), len(result))
				} else {
					for i, v := range result {
						if v != tc.expected[i] {
							t.Errorf("at index %d: expected %v, got %v", i, tc.expected[i], v)
						}
					}
				}
			}
		})
	}
}