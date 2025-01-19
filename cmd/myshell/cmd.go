package main

import "strings"

const (
	CmdExit = "exit"
	CmdEcho = "echo"
	CmdType = "type"
	CmdPwd = "pwd"
)

type Cmd struct {
	Key  string
	Args []string
}

func toCmd(in string) (Cmd, error) {
	data := strings.Split(in, " ")
	if len(data) == 0 {
		return Cmd{Key: in}, nil
	}

	return Cmd{
		Key:  data[0],
		Args: data[1:],
	}, nil
}

func IsBuiltIn(name string) bool {
	return name == CmdExit ||
		name == CmdEcho ||
		name == CmdType
}
