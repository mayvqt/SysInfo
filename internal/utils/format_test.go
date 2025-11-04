package utils

import (
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    uint64
		expected string
	}{
		{
			name:     "zero bytes",
			bytes:    0,
			expected: "0 B",
		},
		{
			name:     "bytes under 1KB",
			bytes:    512,
			expected: "512 B",
		},
		{
			name:     "exactly 1KB",
			bytes:    1024,
			expected: "1.00 KB",
		},
		{
			name:     "kilobytes",
			bytes:    1536,
			expected: "1.50 KB",
		},
		{
			name:     "exactly 1MB",
			bytes:    1024 * 1024,
			expected: "1.00 MB",
		},
		{
			name:     "megabytes",
			bytes:    1024 * 1024 * 100,
			expected: "100.00 MB",
		},
		{
			name:     "exactly 1GB",
			bytes:    1024 * 1024 * 1024,
			expected: "1.00 GB",
		},
		{
			name:     "gigabytes",
			bytes:    1024 * 1024 * 1024 * 5,
			expected: "5.00 GB",
		},
		{
			name:     "exactly 1TB",
			bytes:    1024 * 1024 * 1024 * 1024,
			expected: "1.00 TB",
		},
		{
			name:     "terabytes",
			bytes:    1024 * 1024 * 1024 * 1024 * 2,
			expected: "2.00 TB",
		},
		{
			name:     "fractional gigabytes",
			bytes:    1024 * 1024 * 1024 * 3 / 2,
			expected: "1.50 GB",
		},
		{
			name:     "large value",
			bytes:    1024 * 1024 * 1024 * 1024 * 1024,
			expected: "1.00 PB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("FormatBytes(%d) = %q; want %q", tt.bytes, result, tt.expected)
			}
		})
	}
}

func BenchmarkFormatBytes(b *testing.B) {
	testValues := []uint64{
		0,
		512,
		1024 * 1024,
		1024 * 1024 * 1024,
		1024 * 1024 * 1024 * 1024,
	}

	for _, val := range testValues {
		b.Run("", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				FormatBytes(val)
			}
		})
	}
}
