package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"

	"golang.org/x/term"
)

var errNoTargetFd = errors.New("redirect specified but no target file descriptor got")

func main() {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create term: %v\r\n", err)
		return
	}
	defer term.Restore((fd), oldState)

	var sb strings.Builder
	reader := bufio.NewReader(os.Stdin)

	for {
		if _, err := fmt.Fprint(os.Stdout, "$ "); err != nil {
			fmt.Fprintf(os.Stderr, "write error: %v\r\n", err)
			continue
		}

		for {
			r, _, err := reader.ReadRune()
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not read character: %v\r\n", err)
				continue
			}

			if r == unicode.ReplacementChar {
				fmt.Fprintf(os.Stderr, "invalid char input: %v\r\n", err)
				continue
			} else if r == 0x3 || r == 0x4 {
				return
			} else if r == 0x0A || r == 0xD {
				fmt.Fprintln(os.Stdout)
				break
			} else if r == 0x9 {
				missing := suggestMissing(sb.String())
				if missing == "" {
					continue
				}

				for _, v := range missing {
					sb.WriteRune(v)
					fmt.Fprint(os.Stdout, string(v))
				}

				sb.WriteRune(' ')
			} else {
				sb.WriteRune(r)
				fmt.Fprint(os.Stdout, string(r))
			}
		}

		input := sb.String()
		sb.Reset()

		command, err := toCmd(strings.TrimSpace(input))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid command: %v\r\n", err)
			continue
		}

		fdOut, fdErr, args, err := redirects(command.Args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed redirecting: %v\r\n", err)
			continue
		}

		execCommand(command.Key, args, fdOut, fdErr)
	}
}

func execCommand(programName string, args []string, fdOut *os.File, fdErr *os.File) {
	defer closeFile(fdOut)
	defer closeFile(fdErr)

	stdout := os.Stdout
	if fdOut != nil {
		stdout = fdOut
	}

	stderr := os.Stderr
	if fdErr != nil {
		stderr = fdErr
	}

	switch programName {
	case CmdExit:
		cmdExit(stderr, args)
	case CmdEcho:
		cmdEcho(stdout, args)
	case CmdType:
		cmdType(stdout, args)
	case CmdPwd:
		cmdPwd(stdout, stderr)
	case CmdCd:
		cmdCd(stderr, args)
	default:
		cmdExec(stdout, stderr, programName, args)
	}
}

func redirects(args []string) (*os.File, *os.File, []string, error) {
	if len(args) == 0 {
		return nil, nil, args, nil
	}

	newArgs := make([]string, 0)
	var targetOut, targetErr string
	var appendOut, appendErr bool

	for i := 0; i < len(args); i++ {
		arg := args[i]
		// TODO: refactor
		if arg == ">" || arg == "1>" {
			if i+1 >= len(args) {
				return nil, nil, args, errNoTargetFd
			}
			targetOut = args[i+1]
			i++
		} else if arg == "2>" {
			if i+1 >= len(args) {
				return nil, nil, args, errNoTargetFd
			}
			targetErr = args[i+1]
			i++
		} else if arg == ">>" || arg == "1>>" {
			if i+1 >= len(args) {
				return nil, nil, args, errNoTargetFd
			}
			appendOut = true
			targetOut = args[i+1]
			i++
		} else if arg == "2>>" {
			if i+1 >= len(args) {
				return nil, nil, args, errNoTargetFd
			}
			appendErr = true
			targetErr = args[i+1]
			i++
		} else {
			newArgs = append(newArgs, arg)
		}
	}

	fdOut, err := createFile(targetOut, appendOut)
	if err != nil {
		return nil, nil, newArgs, err
	}

	fdErr, err := createFile(targetErr, appendErr)
	if err != nil {
		closeFile(fdOut)
		return nil, nil, newArgs, err
	}

	return fdOut, fdErr, newArgs, err
}

func closeFile(f *os.File) {
	if f == nil {
		return
	}

	if err := f.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "could not close file %s: %v\r\n", f.Name(), err)
	}
}

func createFile(fileName string, append bool) (*os.File, error) {
	if fileName == "" {
		return nil, nil
	}

	var fd *os.File
	var err error

	if append {
		fd, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	} else {
		fd, err = os.Create(fileName)
	}

	if err != nil {
		return nil, err
	}

	return fd, nil
}
