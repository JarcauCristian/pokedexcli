package repl

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "trimming the input works",
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			name:     "spaces inside names splits the word in two",
			input:    " he llo world ",
			expected: []string{"he", "llo", "world"},
		},
		{
			name:     "only lowercase letters, numbers, underscores and dashes are allowed in entries",
			input:    " pokedex, is fun",
			expected: []string{"pokedex", "is", "fun"},
		},
		{
			name:     "dashes are preserved inside an entry",
			input:    "dashes-in-the_input are accepted",
			expected: []string{"dashes-in-the_input", "are", "accepted"},
		},
		{
			name:     "numbers are preserved inside an entry",
			input:    "there are 10000 pokemons",
			expected: []string{"there", "are", "10000", "pokemons"},
		},
		{
			name:     "only spaces returns an empty string",
			input:    "             ",
			expected: []string{""},
		},
		{
			name:     "an empty string returns an empty string",
			input:    "",
			expected: []string{""},
		},
		{
			name:     "upper case text is first lowercased",
			input:    "EVERYTHING MUST BE lower.",
			expected: []string{"everything", "must", "be", "lower"},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := cleanInput(testCase.input)

			if len(actual) != len(testCase.expected) {
				t.Fatalf(
					"the length of the output %d does not match the expected %d",
					len(actual),
					len(testCase.expected),
				)
			}

			for i := range actual {
				word := actual[i]
				expectedWord := testCase.expected[i]

				if word != expectedWord {
					t.Fatalf("expected word: %s, got word: %s", expectedWord, word)
				}
			}
		})
	}
}
