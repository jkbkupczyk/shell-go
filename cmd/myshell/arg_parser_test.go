package main

import (
	"strings"
	"testing"
)

func TestParseCommand(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		wantArgs []string
	}{
		{
			desc:     "Empty command",
			input:    "",
			wantArgs: []string{},
		},
		{
			desc:     "Blank command",
			input:    "  ",
			wantArgs: []string{},
		},
		{
			desc:     "No arguments",
			input:    "echo",
			wantArgs: []string{"echo"},
		},
		{
			desc:     "No arguments (space)",
			input:    "echo ",
			wantArgs: []string{"echo"},
		},
		{
			desc:     "With 1 unescaped argument",
			input:    "echo hello",
			wantArgs: []string{"echo", "hello"},
		},
		{
			desc:     "With 1 unescaped argument with long spaces",
			input:    "echo                     hello",
			wantArgs: []string{"echo", "hello"},
		},
		{
			desc:     "With 1 escaped argument",
			input:    "echo 'hello'",
			wantArgs: []string{"echo", "hello"},
		},
		{
			desc:     "With 2 unescaped arguments",
			input:    "echo hello world",
			wantArgs: []string{"echo", "hello", "world"},
		},
		{
			desc:     "With 2 unescaped arguments with long spaces",
			input:    "echo hello                                    world",
			wantArgs: []string{"echo", "hello", "world"},
		},
		{
			desc:     "With 1 unescaped argument and one escaped by quote",
			input:    "echo hello 'world'",
			wantArgs: []string{"echo", "hello", "world"},
		},
		{
			desc:     "With 2 arguments escaped by quote",
			input:    "echo 'hello' 'world'",
			wantArgs: []string{"echo", "hello", "world"},
		},
		{
			desc:     "",
			input:    "echo 'shell     example' 'script''hello'",
			wantArgs: []string{"echo", "shell     example", "scripthello"},
		},
		{
			desc:     "",
			input:    "echo '1''2''3''4'",
			wantArgs: []string{"echo", "1234"},
		},
		{
			desc:     "Double quotes",
			input:    "echo \"quz  hello\"  \"bar\"",
			wantArgs: []string{"echo", "quz  hello", "bar"},
		},
		{
			desc:     "Double quotes",
			input:    "echo \"bar\"  \"shell's\"  \"foo\"",
			wantArgs: []string{"echo", "bar", "shell's", "foo"},
		},
		{
			desc:     "Double quotes",
			input:    "echo foo 'bar' \"baz\"",
			wantArgs: []string{"echo", "foo", "bar", "baz"},
		},
		{
			desc:     "Double quotes",
			input:    "echo \"hello world\"",
			wantArgs: []string{"echo", "hello world"},
		},
		{
			desc:     "Double quotes",
			input:    "echo \"\\$\" \"\\`\" \"\\\\\"",
			wantArgs: []string{"echo", "$", "`", "\\"},
		},
		{
			desc:     "Double quotes",
			input:    "echo \"arg with \\\"escaped quotes\\\"\"",
			wantArgs: []string{"echo", "arg with \"escaped quotes\""},
		},
		{
			desc:     "Backslash outside quotes",
			input:    "echo \"before\\   after\"",
			wantArgs: []string{"echo", "before\\   after"},
		},
		{
			desc:     "Backslash outside quotes",
			input:    "echo world\\ \\ \\ \\ \\ \\ script",
			wantArgs: []string{"echo", "world      script"},
		},
		{
			desc:     "Backslash within single quotes",
			input:    "echo 'shell\\\\\\nscript'",
			wantArgs: []string{"echo", "shell\\\\\\nscript"},
		},
		{
			desc:     "Backslash within single quotes",
			input:    "echo 'example\\\"testhello\\\"shell'",
			wantArgs: []string{"echo", "example\\\"testhello\\\"shell"},
		},
		{
			desc:     "Backslash within double quotes",
			input:    `echo "hello'script'\\n'world"`,
			wantArgs: []string{"echo", `hello'script'\n'world`},
		},
		{
			desc:     "Backslash within double quotes",
			input:    `echo "hello\"insidequotes"script\"`,
			wantArgs: []string{"echo", `hello"insidequotesscript"`},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p, err := parseCommand(tC.input)
			if err != nil {
				t.Fatalf("parseCommand returned error: %v", err)
			}
			argsStr := strings.Join(p, ",")
			if len(p) != len(tC.wantArgs) {
				t.Fatalf("parsed args length differs, want: %d, got: %d (args=%s)", len(tC.wantArgs), len(p), argsStr)
			}
			for i, arg := range tC.wantArgs {
				if p[i] != arg {
					t.Errorf("args value differs at index = %d, value wanted: %s, got: %s (args=%s)", i, arg, p[i], argsStr)
				}
			}
		})
	}
}
