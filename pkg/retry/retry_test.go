package retry

import (
	"errors"
	"testing"
	"time"
)

func TestWithBackoff(t *testing.T) {
	tests := []struct {
		name           string
		maxRetries     int
		initialBackoff time.Duration
		operation      func() (string, error)
		expectedResult string
		expectError    bool
	}{
		{
			name:           "SuccessFirstTry",
			maxRetries:     3,
			initialBackoff: 10 * time.Millisecond,
			operation: func() (string, error) {
				return "success", nil
			},
			expectedResult: "success",
			expectError:    false,
		},
		{
			name:           "FailAndRetry",
			maxRetries:     3,
			initialBackoff: 10 * time.Millisecond,
			operation: func() func() (string, error) {
				staticCounter := 0 // Declare persistent counter outside the returned function
				return func() (string, error) {
					staticCounter++
					if staticCounter == 3 {
						return "success", nil
					}
					return "", errors.New("temporary error")
				}
			}(),
			expectedResult: "success",
			expectError:    false,
		},
		{
			name:           "ExceedMaxRetries",
			maxRetries:     3,
			initialBackoff: 10 * time.Millisecond,
			operation: func() (string, error) {
				return "", errors.New("always fails")
			},
			expectedResult: "",
			expectError:    true,
		},
		{
			name:           "ZeroMaxRetries",
			maxRetries:     0,
			initialBackoff: 10 * time.Millisecond,
			operation: func() (string, error) {
				return "", errors.New("fails")
			},
			expectedResult: "",
			expectError:    true,
		},
		{
			name:           "NegativeMaxRetries",
			maxRetries:     -1,
			initialBackoff: 10 * time.Millisecond,
			operation: func() (string, error) {
				return "", errors.New("fails")
			},
			expectedResult: "",
			expectError:    true,
		},
		{
			name:           "OperationAlwaysSucceeds",
			maxRetries:     3,
			initialBackoff: 10 * time.Millisecond,
			operation: func() (string, error) {
				return "always succeeds", nil
			},
			expectedResult: "always succeeds",
			expectError:    false,
		},
		{
			name:           "ImmediateSuccessRetryLimit",
			maxRetries:     3,
			initialBackoff: 10 * time.Millisecond,
			operation: func() (string, error) {
				return "immediate success", nil
			},
			expectedResult: "immediate success",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := WithBackoff(tt.maxRetries, tt.initialBackoff, tt.operation)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error, got %v", err)
				}
				if result != tt.expectedResult {
					t.Errorf("expected result %v, got %v", tt.expectedResult, result)
				}
			}
		})
	}
}
