package main

import (
	"strings"
)

const (
	CmdExit = "exit"
	CmdEcho = "echo"
	CmdType = "type"
	CmdPwd  = "pwd"
	CmdCd   = "cd"
)

type Cmd struct {
	Key  string
	Args []string
}

func toCmd(in string) (Cmd, error) {
	args := parseCommand(in)

	if len(args) == 0 {
		return Cmd{Key: in}, nil
	}

	return Cmd{
		Key:  args[0],
		Args: args[1:],
	}, nil
}

func parseCommand(in string) []string {
	quoted := false
	var sb strings.Builder
	tokens := make([]string, 0)

	appendToken := func() {
		if sb.Len() == 0 {
			return
		}
		tokens = append(tokens, sb.String())
		sb.Reset()
	}

	for _, ch := range in {
		switch ch {
		case ' ':
			{
				if quoted {
					sb.WriteRune(ch)
					continue
				}
				appendToken()
			}
		case '\'':
			{
				if quoted {
					// already quoted, parse argument
					appendToken()
					quoted = false
				} else {
					quoted = true
				}
			}
		default:
			sb.WriteRune(ch)
		}
	}

	// append remaining token
	appendToken()

	return tokens
}

func IsBuiltIn(name string) bool {
	return name == CmdExit ||
		name == CmdEcho ||
		name == CmdType ||
		name == CmdPwd ||
		name == CmdCd
}
