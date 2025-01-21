package main

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
	start := 0
	values := make([]string, 0)

	for i, ch := range in {
		switch ch {
		case ' ':
			{
				if !quoted {
					values = append(values, in[start:i])
					start = i
				}
			}
		case '\'':
			{
				// already quoted, end arg
				if quoted {
					values = append(values, in[start+1:i])
					quoted = false
				} else {
					start = i
					quoted = true
				}
			}
		}

	}

	return values
}

func IsBuiltIn(name string) bool {
	return name == CmdExit ||
		name == CmdEcho ||
		name == CmdType ||
		name == CmdPwd ||
		name == CmdCd
}
