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
				} else {
					if !pathExists {
						fmt.Fprintf(os.Stdout, "%s: not found\r\n", arg)
						continue
					}

					var found bool
					for _, p := range osPaths {
						files, _ := os.ReadDir(p)
						if len(files) == 0 {
							continue
						}

						for _, f := range files {
							if !f.IsDir() && f.Name() == arg {
								fmt.Fprintf(os.Stdout, "%s is %s\r\n", arg, filepath.Join(p, f.Name()))
								found = true
								break
							}
						}
					}

					if !found {
						fmt.Fprintf(os.Stdout, "%s: not found\r\n", arg)
					}
				}
			}
		default:
			fileName := command.Key

			var found bool
			for _, p := range osPaths {
				files, _ := os.ReadDir(p)
				if len(files) == 0 {
					continue
				}

				for _, f := range files {
					if !f.IsDir() && f.Name() == fileName {
						fullPath := filepath.Join(p, f.Name())
						out, err := exec.Command(fullPath, command.Args...).Output()
						if err != nil {
							fmt.Fprintf(os.Stdout, "Could not execute command %s: %v\r\n", fileName, err)
							continue
						}
						found = true
						fmt.Fprint(os.Stdout, string(out), "\r\n")
						break
					}
				}
			}

			if !found {
				fmt.Fprintf(os.Stdout, "%s: command not found\r\n", command.Key)
			}
		}
	}
}

func fileType(command, path string) FileType {
	if IsBuiltIn(command) {
		return FileTypeBuiltIn
	}

	return FileTypeExecutable
}
