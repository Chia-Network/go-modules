package amazon

import "testing"

func TestGetTotalNumberParts(t *testing.T) {
	tests := []struct {
		FileSize int64
		PartSize int64
		Expected int64
	}{
		{5, 2, 3}, // Test quotient where division produced a non-whole number
		{2, 1, 2}, // Test quotient where division produced a whole number
		// Test with a couple much larger numbers
		{104857600, 5242880, 20},       // 100MB file, 5MB chunks, 20 total chunks
		{132070244351, 8388200, 15745}, // ~123GB file, ~8MB chunks, 15745 total chunks
	}

	for _, test := range tests {
		result := getTotalNumberParts(test.FileSize, test.PartSize)
		if result != test.Expected {
			t.Errorf("operation failed for %d / %d. Expected %d, got %d", test.FileSize, test.PartSize, test.Expected, result)
		}
	}
}
