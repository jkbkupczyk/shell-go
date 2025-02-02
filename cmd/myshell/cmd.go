package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

const (
	CmdExit = "exit"
	CmdEcho = "echo"
	CmdType = "type"
	CmdPwd  = "pwd"
	CmdCd   = "cd"
)

func isBuiltIn(name string) bool {
	return name == CmdExit ||
		name == CmdEcho ||
		name == CmdType ||
		name == CmdPwd ||
		name == CmdCd
}

type Cmd struct {
	Key  string
	Args []string
}

func toCmd(in string) (Cmd, error) {
	args, err := parseCommand(in)
	if err != nil {
		return Cmd{}, err
	}

	if len(args) == 0 {
		return Cmd{Key: in}, nil
	}

	return Cmd{
		// first element is program name
		Key:  args[0],
		Args: args[1:],
	}, nil
}

func listSuggestions(value string) []string {
	if value == "" {
		return []string{}
	}

	suggestions := make([]string, 0)
	availableCommands := append([]string{CmdExit, CmdEcho, CmdType, CmdPwd, CmdCd}, listPathCommands()...)

	for _, cmd := range availableCommands {
		if (value == cmd || strings.HasPrefix(cmd, value)) && !slices.Contains(suggestions, cmd) {
			suggestions = append(suggestions, cmd)
		}
	}

	slices.Sort(suggestions)

	return suggestions
}

func hasLongestCommonPrefix(suggestions []string) bool {
	if len(suggestions) == 0 {
		return false
	}
	if len(suggestions) == 1 {
		return true
	}

	// suggestions already includes prefix and are already sorted in ASC order
	first := suggestions[0]
	last := suggestions[len(suggestions)-1]

	return len(last) == len(first)
}

func cmdExit(stderr io.Writer, args []string) {
	if len(args) == 0 {
		os.Exit(0)
	}

	exitCode, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "invalid exit code value: %s\r\n", args[0])
	}

	os.Exit(exitCode)
}

func cmdEcho(stdout io.Writer, args []string) {
	fmt.Fprint(stdout, strings.Join(args, " "), "\r\n")
}

func cmdType(stdout io.Writer, args []string) {
	if len(args) == 0 {
		return
	}

	if isBuiltIn(args[0]) {
		fmt.Fprintf(stdout, "%s is a shell builtin\r\n", args[0])
		return
	}

	filePath := findFile(args[0])
	if filePath == "" {
		fmt.Fprintf(stdout, "%s: not found\r\n", args[0])
		return
	}

	fmt.Fprintf(stdout, "%s is %s\r\n", args[0], filePath)
}

func cmdPwd(stdout io.Writer, stderr io.Writer) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(stderr, "could not get working dir: %v\r\n", err)
		return
	}

	fmt.Fprintln(stdout, wd)
}

func cmdCd(stderr io.Writer, args []string) {
	if len(args) == 0 {
		return
	}

	targetDir := args[0]
	if targetDir == "~" {
		targetDir, _ = os.UserHomeDir()
	}

	if err := os.Chdir(targetDir); err != nil {
		fmt.Fprintf(stderr, "cd: %s: No such file or directory\r\n", targetDir)
		return
	}
}

func cmdExec(stdin io.Reader, stdout io.Writer, stderr io.Writer, execName string, args []string) {
	filePath := findFile(execName)
	if filePath == "" {
		fmt.Fprintf(stdout, "%s: command not found\r\n", execName)
		return
	}

	c := exec.Command(execName, args...)
	c.Stdin = stdin
	c.Stdout = stdout
	c.Stderr = stderr

	c.Run()
}

func findFile(fileName string) string {
	pathVal, pathExists := os.LookupEnv("PATH")
	if !pathExists {
		return ""
	}

	paths := strings.Split(pathVal, string(os.PathListSeparator))
	for _, p := range paths {
		files, _ := os.ReadDir(p)
		if len(files) == 0 {
			continue
		}

		for _, f := range files {
			if !f.IsDir() && f.Name() == fileName {
				return filepath.Join(p, f.Name())
			}
		}
	}

	return ""
}

func listPathCommands() []string {
	pathVal, pathExists := os.LookupEnv("PATH")
	if !pathExists {
		return []string{}
	}

	cmds := make([]string, 0)

	for _, p := range strings.Split(pathVal, string(os.PathListSeparator)) {
		files, _ := os.ReadDir(p)
		if len(files) == 0 {
			continue
		}

		for _, f := range files {
			if !f.IsDir() {
				cmds = append(cmds, f.Name())
			}
		}
	}

	return cmds
}
