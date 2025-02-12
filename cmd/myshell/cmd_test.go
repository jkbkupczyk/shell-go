package main

import (
	"bytes"
	"testing"
)

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
		{
			desc:    "Command cd is builtin",
			cmd:     "cd",
			builtin: true,
		},
		{
			desc:    "Command ls is not a builtin",
			cmd:     "ls",
			builtin: false,
		},
		{
			desc:    "Command cat is not a builtin",
			cmd:     "cat",
			builtin: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := isBuiltIn(tC.cmd); got != tC.builtin {
				t.Errorf("isBuiltIn(%s), want = %t, got: %t", tC.cmd, tC.builtin, got)
			}
		})
	}
}

// for more sophisticated tests see arg_parser_test.go
func TestToCmd(t *testing.T) {
	testCases := []struct {
		desc       string
		rawCommand string
		wantCmd    Cmd
	}{
		{
			desc:       "empty command",
			rawCommand: "",
			wantCmd:    Cmd{"", nil},
		},
		{
			desc:       "only command",
			rawCommand: "test",
			wantCmd:    Cmd{Key: "test", Args: nil},
		},
		{
			desc:       "command with empty args",
			rawCommand: "test ",
			wantCmd:    Cmd{Key: "test", Args: nil},
		},
		{
			desc:       "command with 1 arg",
			rawCommand: "test hi",
			wantCmd:    Cmd{Key: "test", Args: []string{"hi"}},
		},
		{
			desc:       "command with multiple args",
			rawCommand: "test 1 '2' \"3\"",
			wantCmd:    Cmd{Key: "test", Args: []string{"1", "2", "3"}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			cmd, err := toCmd(tC.rawCommand)
			if err != nil {
				t.Fatalf("toCmd returned error: %v", err)
			}
			if cmd.Key != tC.wantCmd.Key {
				t.Errorf("invalid command key, want: '%s', got: '%v'", tC.wantCmd.Key, cmd.Key)
			}
			if len(cmd.Args) != len(tC.wantCmd.Args) {
				t.Fatalf("invalid number of args, want: %d, got: %d", len(tC.wantCmd.Args), len(cmd.Args))
			}
			for i, wantArg := range tC.wantCmd.Args {
				if cmd.Args[i] != wantArg {
					t.Errorf("args[%d] want: %s, got: %s", i, wantArg, cmd.Args[i])
				}
			}
		})
	}
}

func TestCmdEcho(t *testing.T) {
	testCases := []struct {
		desc       string
		args       []string
		wantResult string
	}{
		{
			desc:       "NIL args",
			args:       nil,
			wantResult: "\r\n",
		},
		{
			desc:       "empty args",
			args:       []string{},
			wantResult: "\r\n",
		},
		{
			desc:       "one arg",
			args:       []string{"john"},
			wantResult: "john\r\n",
		},
		{
			desc:       "multiple args",
			args:       []string{"john", "richard", "doe"},
			wantResult: "john richard doe\r\n",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			buff := make([]byte, 0)
			w := bytes.NewBuffer(buff)

			cmdEcho(w, tC.args)

			result := w.String()
			if result != tC.wantResult {
				t.Errorf("invalid command result, wanted: '%s', got: '%s'", tC.wantResult, result)
			}
		})
	}
}
