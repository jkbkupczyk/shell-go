package main

import (
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
