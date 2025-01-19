package main

import (
	"fmt"
)

func errCmdNotFound(cmd string) error {
	return fmt.Errorf("%s: command not found", cmd)
}
