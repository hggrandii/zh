package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: zh [zig arguements]")
		os.Exit(1)
	}

	args := os.Args[1:]

	cmd := exec.Command("zig", args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}

		fmt.Fprintf(os.Stderr, "Failed to execute zig : %v \n", err)
		os.Exit(1)
	}
}
