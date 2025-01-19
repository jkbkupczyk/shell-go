package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type FileType int

const (
	FileTypeBuiltIn FileType = iota
	FileTypeExecutable
)

type ControlFlow int

const (
	FlowContinue ControlFlow = iota
	FlowBreak
)

func main() {
	for {
		if _, err := fmt.Fprint(os.Stdout, "$ "); err != nil {
			fmt.Fprintf(os.Stdout, "write error: %v\r\n", err)
			continue
		}

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stdout, "cannot read input: %v\r\n", err)
			continue
		}

		command, err := toCmd(strings.TrimSpace(input))
		if err != nil {
			fmt.Fprintf(os.Stdout, "invalid command: %v\r\n", err)
			continue
		}

		pathVal, pathExists := os.LookupEnv("PATH")
		osPaths := strings.Split(pathVal, string(os.PathListSeparator))

		switch command.Key {
		case CmdExit:
			{
				if len(command.Args) == 0 {
					os.Exit(0)
				}

				exitCode, err := strconv.Atoi(command.Args[0])
				if err != nil {
					fmt.Fprintf(os.Stdout, "invalid exit code value: %s\r\n", command.Args[0])
					continue
				}

				os.Exit(exitCode)
			}
		case CmdEcho:
			fmt.Fprint(os.Stdout, strings.Join(command.Args, " "), "\r\n")
		case CmdType:
			{
				if len(command.Args) == 0 {
					continue
				}

				arg := command.Args[0]

				if IsBuiltIn(arg) {
					fmt.Fprintf(os.Stdout, "%s is a shell builtin\r\n", arg)
					continue
				}

				if !pathExists {
					fmt.Fprintf(os.Stdout, "%s: not found\r\n", arg)
					continue
				}

				filePath := findFile(arg, osPaths)
				if filePath == "" {
					fmt.Fprintf(os.Stdout, "%s: not found\r\n", arg)
					continue
				}

				fmt.Fprintf(os.Stdout, "%s is %s\r\n", arg, filePath)
			}
		case CmdPwd:
			{
				wd, err := os.Getwd()
				if err != nil {
					fmt.Fprintf(os.Stdout, "could not get working dir: %v\r\n", err)
					return
				}

				fmt.Fprintln(os.Stdout, wd)
			}
		case CmdCd:
			{
				if len(command.Args) == 0 {
					continue
				}

				arg := command.Args[0]

				if arg == "~" {
					arg, _ = os.UserHomeDir()
				}

				if err := os.Chdir(arg); err != nil {
					fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\r\n", arg)
					continue
				}
			}
		default:
			filePath := findFile(command.Key, osPaths)
			if filePath == "" {
				fmt.Fprintf(os.Stdout, "%s: command not found\r\n", command.Key)
				continue
			}

			c := exec.Command(command.Key, command.Args...)
			c.Stderr = os.Stderr
			c.Stdout = os.Stdout

			if err := c.Run(); err != nil {
				fmt.Fprintf(os.Stdout, "could not execute command %s: %v\r\n", filePath, err)
				continue
			}
		}
	}
}

func findFile(fileName string, paths []string) string {
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

func fileType(command, path string) FileType {
	if IsBuiltIn(command) {
		return FileTypeBuiltIn
	}

	return FileTypeExecutable
}
