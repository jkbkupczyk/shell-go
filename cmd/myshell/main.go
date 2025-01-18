package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		case "exit":
			{
				if len(command.Args) == 0 {
					continue
				}

				exitCode, err := strconv.Atoi(command.Args[0])
				if err != nil {
					fmt.Fprintf(os.Stdout, "invalid exit code value: %s\r\n", command.Args[0])
					continue
				}

				os.Exit(exitCode)
			}
		case "echo":
			{
				fmt.Fprint(os.Stdout, strings.Join(command.Args, " "), "\r\n")
			}
		default:
			{
				fmt.Fprintf(os.Stdout, "%s: command not found\r\n", command.Key)
			}
		}
	}

}
