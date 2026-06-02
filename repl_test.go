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
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "  single  ",
			expected: []string{"single"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q): expected %d words, got %d", c.input, len(c.expected), len(actual))
			continue
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("cleanInput(%q): word %d expected %q, got %q", c.input, i, c.expected[i], actual[i])
			}
		}
	}
}
