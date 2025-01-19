package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type FileType int

const (
	FileTypeBuiltIn FileType = iota
	FileTypeExecutable
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

				val, exists := os.LookupEnv("PATH")
				osPaths := strings.Split(val, string(os.PathListSeparator))

				for _, arg := range command.Args {
					if IsBuiltIn(arg) {
						fmt.Fprintf(os.Stdout, "%s is a shell builtin\r\n", arg)
					} else {
						if !exists {
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
									fmt.Fprintf(os.Stdout, "%s is %s\r\n", arg, fileNameDisplay(p, f.Name()))
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
			}
		default:
			fmt.Fprintf(os.Stdout, "%s: command not found\r\n", command.Key)
		}
	}
}

func fileType(command, path string) FileType {
	if IsBuiltIn(command) {
		return FileTypeBuiltIn
	}

	return FileTypeExecutable
}

func fileNameDisplay(path, fName string) string {
	if path[len(path)-1:] == string(os.PathSeparator) {
		return path + fName
	}
	return path + string(os.PathSeparator) + fName
}
