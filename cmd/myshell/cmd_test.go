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
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			p := parseCommand(tC.input)
			if len(p) != len(tC.wantArgs) {
				t.Fatalf("parsed args length differs, want: %d, got: %d (args=%s)", len(tC.wantArgs), len(p), strings.Join(p, ","))
			}
			for i, arg := range tC.wantArgs {
				if p[i] != arg {
					t.Errorf("args value differs at index = %d, value wanted: %s, got: %s", i, arg, p[i])
				}
			}
		})
	}
}

func TestIsBuiltIn(t *testing.T) {
	testCases := []struct {
		desc    string
		cmd     string
		builtin bool
	}{
		{
			desc:    "Command exit is builtin",
			cmd:     "exit",
			builtin: true,
		},
		{
			desc:    "Command echo is builtin",
			cmd:     "echo",
			builtin: true,
		},
		{
			desc:    "Command type is builtin",
			cmd:     "type",
			builtin: true,
		},
		{
			desc:    "Command pwd is builtin",
			cmd:     "pwd",
			builtin: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := IsBuiltIn(tC.cmd); got != tC.builtin {
				t.Errorf("IsBuiltIn(%s), want = %t, got: %t", tC.cmd, tC.builtin, got)
			}
		})
	}
}
