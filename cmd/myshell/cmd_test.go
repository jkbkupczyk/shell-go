package main

import "testing"

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
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := IsBuiltIn(tC.cmd); got != tC.builtin {
				t.Errorf("IsBuiltIn(%s), want = %t, got: %t", tC.cmd, tC.builtin, got)
			}
		})
	}
}
