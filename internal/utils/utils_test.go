package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLastFourDigits(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"standard card", "1234567812345678", "5678"},
		{"short card", "123456789012345", "2345"},
		{"very short", "123", "123"},
		{"exactly four", "1234", "1234"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLastFourDigits(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatExpiryDate(t *testing.T) {
	tests := []struct {
		name     string
		month    int
		year     int
		expected string
	}{
		{"2-digit year", 12, 25, "12/25"},
		{"4-digit year", 12, 2025, "12/25"},
		{"single digit month", 5, 2026, "05/26"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatExpiryDate(tt.month, tt.year)
			assert.Equal(t, tt.expected, result)
		})
	}
}
