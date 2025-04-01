package retry

import (
	"errors"
	"time"
)

// WithBackoff attempts a function up to maxRetries times, with exponential backoff
func WithBackoff[T any](maxRetries int, initialBackoff time.Duration, operation func() (T, error)) (T, error) {
	var result T
	var err error

	backoff := initialBackoff // Use configurable initial backoff duration
	for i := 0; i < maxRetries; i++ {
		result, err = operation()
		if err == nil {
			return result, nil // Success, return the result
		}

		if i < maxRetries-1 { // Wait before retrying (if more retries remain)
			time.Sleep(backoff)
			backoff *= 2 // Exponential backoff
		}
	}

	var zeroValue T // Zero value for T (e.g., nil for pointers, empty struct for structs)
	return zeroValue, errors.New("operation failed after maximum retries")
}
