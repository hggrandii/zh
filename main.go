package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hggrandii/zh/cmd"
)

func main() {
	if len(os.Args) < 2 {
		cmd.PrintUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		cmd.HandleAddCommand()
		return
	case "new":
		cmd.HandleNewCommand()
		return
	case "version", "--version", "-v":
		cmd.HandleVersionCommand()
		return
	default:
		args := os.Args[1:]
		zigCmd := exec.Command("zig", args...)
		zigCmd.Stdin = os.Stdin
		zigCmd.Stdout = os.Stdout
		zigCmd.Stderr = os.Stderr
		if err := zigCmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			fmt.Fprintf(os.Stderr, "Failed to execute zig: %v\n", err)
			os.Exit(1)
		}
	}
}
