package xlog

import (
	"bytes"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrettyLevelToColor(t *testing.T) {
	// Test that all log levels have a corresponding color
	levels := []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
		zapcore.DPanicLevel,
		zapcore.PanicLevel,
		zapcore.FatalLevel,
	}

	for _, level := range levels {
		color, ok := _prettyLevelToColor[level]
		if !ok {
			t.Errorf("Expected color for level %v, but none found", level)
			continue
		}

		if color < Black || color > White {
			t.Errorf("Expected color to be between %d and %d, got %d", Black, White, color)
		}
	}
}

func TestColor_Add(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		input string
		want  string
	}{
		{
			name:  "red color",
			color: Red,
			input: "error message",
			want:  "\x1b[31merror message\x1b[0m",
		},
		{
			name:  "blue color",
			color: Blue,
			input: "info message",
			want:  "\x1b[34minfo message\x1b[0m",
		},
		{
			name:  "yellow color",
			color: Yellow,
			input: "warning message",
			want:  "\x1b[33mwarning message\x1b[0m",
		},
		{
			name:  "empty string",
			color: Green,
			input: "",
			want:  "\x1b[32m\x1b[0m",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.color.Add(tc.input)
			if got != tc.want {
				t.Errorf("Color.Add() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestPrintJSON(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Restore stdout when the test finishes
	defer func() {
		os.Stdout = oldStdout
	}()

	// Test with a map
	mapObj := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	// Call printJSON with nil color (no coloring)
	err = printJSON(nil, mapObj)
	if err != nil {
		t.Errorf("printJSON() error = %v", err)
	}

	// Close the writer to get the output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that the output contains the expected keys and values
	if !strings.Contains(output, "key1") || !strings.Contains(output, "value1") {
		t.Errorf("printJSON() output doesn't contain expected key1/value1, got: %s", output)
	}
	if !strings.Contains(output, "key2") || !strings.Contains(output, "42") {
		t.Errorf("printJSON() output doesn't contain expected key2/42, got: %s", output)
	}

	// Test with a slice
	// Reset the pipe
	r, w, err = os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	sliceObj := []interface{}{1, "two", true}

	// Call printJSON with a color
	color := Red
	err = printJSON(&color, sliceObj)
	if err != nil {
		t.Errorf("printJSON() error = %v", err)
	}

	// Close the writer to get the output
	w.Close()
	buf.Reset()
	io.Copy(&buf, r)
	output = buf.String()

	// Check that the output contains the expected values
	// Note: The output will be colored, but we can still check for the values
	if !strings.Contains(output, "1") || !strings.Contains(output, "two") || !strings.Contains(output, "true") {
		t.Errorf("printJSON() output doesn't contain expected values, got: %s", output)
	}

	// Test with invalid JSON
	// Reset the pipe
	r, w, err = os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Create a circular reference that can't be marshaled to JSON
	type circular struct {
		Self *circular
	}
	c := &circular{}
	c.Self = c

	// This should return an error
	err = printJSON(nil, c)
	if err == nil {
		t.Error("printJSON() expected error for circular reference, got nil")
	}

	// Test with a string that's not valid JSON
	// Reset the pipe
	r, w, err = os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	// This should return an error because it's not valid JSON
	err = printJSON(nil, "not-json")
	if err == nil {
		t.Error("printJSON() expected error for invalid JSON string, got nil")
	}
}