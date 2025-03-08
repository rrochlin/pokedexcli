package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "I hate TDD it sucks",
			expected: []string{"i", "hate", "tdd", "it", "sucks"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(c.expected) != len(actual) {
			t.Errorf("mismatched object length\nActual:%v, L:%v\nExpected:%v, L:%v", actual, len(actual), c.expected, len(c.expected))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("mismatched words in postion %v\nActual:%v\nExpected:%v", i, actual, c.expected)
			}

		}
	}
}
