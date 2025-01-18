package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func main() {

	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			slog.Error("Cannot read input", slog.Any("err", err))
			os.Exit(1)
		}

		command := strings.TrimSpace(input)

		switch command {
		default:
			{
				fmt.Fprintf(os.Stdout, "%s: command not found\r\n", command)
			}
		}
	}

}
