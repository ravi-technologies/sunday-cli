package cli

import "testing"

// TestTruncate_Short verifies that the truncate function returns the original
// string unchanged when it is shorter than the maximum length.
func TestTruncate_Short(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		max      int
		expected string
	}{
		{
			name:     "short string with large max",
			input:    "hello",
			max:      100,
			expected: "hello",
		},
		{
			name:     "empty string",
			input:    "",
			max:      10,
			expected: "",
		},
		{
			name:     "single character",
			input:    "a",
			max:      5,
			expected: "a",
		},
		{
			name:     "string much shorter than max",
			input:    "hi",
			max:      50,
			expected: "hi",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := truncate(tc.input, tc.max)
			if result != tc.expected {
				t.Errorf("truncate(%q, %d) = %q, want %q", tc.input, tc.max, result, tc.expected)
			}
		})
	}
}

// TestTruncate_Long verifies that the truncate function adds an ellipsis
// when the string exceeds the maximum length.
func TestTruncate_Long(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		max      int
		expected string
	}{
		{
			name:     "long string truncated",
			input:    "hello world, this is a long string",
			max:      10,
			expected: "hello w...",
		},
		{
			name:     "exactly needs truncation",
			input:    "abcdefghij",
			max:      8,
			expected: "abcde...",
		},
		{
			name:     "unicode string truncation",
			input:    "hello world",
			max:      8,
			expected: "hello...",
		},
		{
			name:     "truncate to minimum useful length",
			input:    "abcdefghij",
			max:      5,
			expected: "ab...",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := truncate(tc.input, tc.max)
			if result != tc.expected {
				t.Errorf("truncate(%q, %d) = %q, want %q", tc.input, tc.max, result, tc.expected)
			}

			// Verify the result doesn't exceed max length
			if len(result) > tc.max {
				t.Errorf("truncate(%q, %d) = %q (len=%d), exceeds max %d",
					tc.input, tc.max, result, len(result), tc.max)
			}

			// Verify ellipsis is present when truncated
			if len(tc.input) > tc.max && len(result) > 3 {
				if result[len(result)-3:] != "..." {
					t.Errorf("truncate(%q, %d) = %q, expected ellipsis suffix", tc.input, tc.max, result)
				}
			}
		})
	}
}

// TestTruncate_Exact verifies the boundary condition when the string length
// is exactly equal to the maximum length.
func TestTruncate_Exact(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		max      int
		expected string
	}{
		{
			name:     "exact length match",
			input:    "hello",
			max:      5,
			expected: "hello", // No truncation needed
		},
		{
			name:     "exact length with 10 chars",
			input:    "0123456789",
			max:      10,
			expected: "0123456789", // No truncation needed
		},
		{
			name:     "one char over max",
			input:    "abcdef",
			max:      5,
			expected: "ab...", // One over, needs truncation
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := truncate(tc.input, tc.max)
			if result != tc.expected {
				t.Errorf("truncate(%q, %d) = %q, want %q", tc.input, tc.max, result, tc.expected)
			}
		})
	}
}

// TestTruncate_MinMax verifies that the truncate function handles
// edge cases where max is less than 4 (minimum useful truncation length).
// The function should enforce a minimum of 4 to allow "X..." format.
func TestTruncate_MinMax(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		max         int
		expectMin   int // The effective minimum that should be applied
		description string
	}{
		{
			name:        "max less than 4 with long string",
			input:       "hello world",
			max:         3,
			expectMin:   4,
			description: "max=3 should be treated as max=4",
		},
		{
			name:        "max equals 1 with long string",
			input:       "hello world",
			max:         1,
			expectMin:   4,
			description: "max=1 should be treated as max=4",
		},
		{
			name:        "max equals 0 with long string",
			input:       "hello world",
			max:         0,
			expectMin:   4,
			description: "max=0 should be treated as max=4",
		},
		{
			name:        "negative max with long string",
			input:       "hello world",
			max:         -5,
			expectMin:   4,
			description: "negative max should be treated as max=4",
		},
		{
			name:        "max exactly 4 with long string",
			input:       "hello world",
			max:         4,
			expectMin:   4,
			description: "max=4 should remain max=4",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := truncate(tc.input, tc.max)

			// When input is longer than effective max, result should be exactly effectiveMax
			if len(tc.input) > tc.expectMin {
				if len(result) != tc.expectMin {
					t.Errorf("truncate(%q, %d) = %q (len=%d), expected len=%d. %s",
						tc.input, tc.max, result, len(result), tc.expectMin, tc.description)
				}

				// Should end with ellipsis
				if result[len(result)-3:] != "..." {
					t.Errorf("truncate(%q, %d) = %q, expected ellipsis suffix. %s",
						tc.input, tc.max, result, tc.description)
				}
			}
		})
	}

	// Test with short string that doesn't need truncation even with small max
	t.Run("short string with small max", func(t *testing.T) {
		result := truncate("hi", 2)
		// After min enforcement (max becomes 4), "hi" (len=2) is <= 4, so no truncation
		if result != "hi" {
			t.Errorf("truncate(\"hi\", 2) = %q, want \"hi\" (no truncation needed after min enforcement)", result)
		}
	})
}

// TestTruncate_Comprehensive runs additional edge case tests
// to ensure robust behavior of the truncate function.
func TestTruncate_Comprehensive(t *testing.T) {
	// Test that result length never exceeds max (after min enforcement)
	t.Run("result never exceeds effective max", func(t *testing.T) {
		testInputs := []string{
			"",
			"a",
			"ab",
			"abc",
			"abcd",
			"abcde",
			"hello world this is a test string",
			"123456789012345678901234567890",
		}

		for _, input := range testInputs {
			for max := 0; max <= 20; max++ {
				result := truncate(input, max)

				// Effective max is at least 4
				effectiveMax := max
				if effectiveMax < 4 {
					effectiveMax = 4
				}

				if len(result) > effectiveMax {
					t.Errorf("truncate(%q, %d) = %q (len=%d), exceeds effective max %d",
						input, max, result, len(result), effectiveMax)
				}
			}
		}
	})

	// Test that short strings pass through unchanged
	t.Run("short strings unchanged", func(t *testing.T) {
		testCases := []struct {
			input string
			max   int
		}{
			{"", 10},
			{"a", 10},
			{"hello", 10},
			{"exact", 5},
		}

		for _, tc := range testCases {
			result := truncate(tc.input, tc.max)
			if len(tc.input) <= tc.max && result != tc.input {
				t.Errorf("truncate(%q, %d) = %q, expected original string unchanged",
					tc.input, tc.max, result)
			}
		}
	})
}
