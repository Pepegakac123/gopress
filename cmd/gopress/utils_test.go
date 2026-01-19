package main

import (
	"testing"
)

func TestSanitizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal path",
			input:    `C:\Users\test\images`,
			expected: `C:\Users\test\images`,
		},
		{
			name:     "Path with double quotes",
			input:    `"C:\Users\test\images"`,
			expected: `C:\Users\test\images`,
		},
		{
			name:     "Path with single quotes",
			input:    `'C:\Users\test\images'`,
			expected: `C:\Users\test\images`,
		},
		{
			name:     "Path with trailing space and quotes",
			input:    ` "C:\Users\test\images" `,
			expected: `C:\Users\test\images`,
		},
		{
			name:     "Empty path",
			input:    ``,
			expected: ``,
		},
		{
			name:     "Path with spaces inside",
			input:    `"C:\My Documents\Images"`,
			expected: `C:\My Documents\Images`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizePath(tt.input)
			if got != tt.expected {
				t.Errorf("sanitizePath(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
