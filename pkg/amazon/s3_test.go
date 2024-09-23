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
		{132070244351, 8388200, 10000}, // ~123GB file, ~8MB chunks, 15745 total chunks - limit is 10,000 so expect 10k
		{100658400000, 8388200, 10000}, // Testing a perfect division into parts
	}

	for _, test := range tests {
		result := getTotalNumberParts(test.FileSize, test.PartSize)
		if result != test.Expected {
			t.Errorf("operation failed for %d / %d. Expected %d, got %d", test.FileSize, test.PartSize, test.Expected, result)
		}
	}
}

func TestGetPartSize(t *testing.T) {
	tests := []struct {
		FileSize int64
		PartSize int64
		Expected int64
	}{
		{5, 2, 2}, // Test quotient where division produced a non-whole number
		{2, 1, 1}, // Test quotient where division produced a whole number
		// Test with a couple much larger numbers
		{104857600, 5242880, 5242880},     // 100MB file, 5MB chunks, 20 total chunks
		{132070244351, 8388200, 13208345}, // ~123GB file, ~8MB chunks, 15745 total chunks - limit is 10,000 so expect 10k
		{132083450000, 8388200, 13208345}, // Should be an exact filesize to fit perfectly into the parts
		{132083450001, 8388200, 13209665}, // One too big to fit perfectly into the parts, so expect larger part sizes to accommodate
		{132096658344, 8388200, 13210986}, // -1 too big to fit perfectly into the parts, so expect larger part sizes to accommodate
	}

	for _, test := range tests {
		numParts := getTotalNumberParts(test.FileSize, test.PartSize)
		result := getPartSize(test.FileSize, numParts, test.PartSize)
		if result != test.Expected {
			t.Errorf("operation failed for %d / %d. Expected %d, got %d", test.FileSize, test.PartSize, test.Expected, result)
		}
	}
}
