package commands

import (
	"fmt"
	"testing"
)

var cliCommands = map[string]CliCommand{
	"success": {
		Name:        "success",
		Description: "Always returns nil.",
		Callback:    func(config *Config) error { return nil },
	},
	"failure": {
		Name:        "failure",
		Description: "Always returns an error.",
		Callback:    func(config *Config) error { return fmt.Errorf("failed") },
	},
}

type expectedOutput struct {
	command CliCommand
	output  string
}

func executeCommand(command string, config *Config) expectedOutput {
	if c, ok := cliCommands[command]; ok {
		response := c.Callback(config)
		output := expectedOutput{
			command: c,
			output:  "",
		}

		if response != nil {
			output.output = response.Error()
		}
		return output
	}
	return expectedOutput{
		command: CliCommand{},
		output:  "",
	}
}

func TestCommands(t *testing.T) {
	config := Config{
		Commands: cliCommands,
	}

	cases := []struct {
		name     string
		input    string
		expected expectedOutput
	}{
		{
			name:  "an unknown command returns nil",
			input: "unknown",
			expected: expectedOutput{
				command: CliCommand{},
				output:  "",
			},
		},
		{
			name:  "a successful command returns nil",
			input: "success",
			expected: expectedOutput{
				command: cliCommands["success"],
				output:  "",
			},
		},
		{
			name:  "a failed command returns an error",
			input: "failure",
			expected: expectedOutput{
				command: cliCommands["failure"],
				output:  "failed",
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			output := executeCommand(testCase.input, &config)
			actualCommand := output.command
			expectedCommand := testCase.expected.command

			if expectedCommand.Name != actualCommand.Name {
				t.Fatalf(
					"expected command name to be: %s, got: %s",
					expectedCommand.Name,
					actualCommand.Name,
				)
			}

			if testCase.expected.output != output.output {
				t.Fatalf("expected output: %v, got: %v", testCase.expected.output, output.output)
			}
		})
	}
}
