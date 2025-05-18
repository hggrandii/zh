package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hggrandii/zh/internal/deps"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: zh [command] [arguments]")
		fmt.Fprintln(os.Stderr, "\nCommands:")
		fmt.Fprintln(os.Stderr, "  add [github-url]    Add a dependency from GitHub using the latest commit")
		fmt.Fprintln(os.Stderr, "  [any zig command]   Pass arguments directly to the zig command")
		os.Exit(1)
	}

	if os.Args[1] == "add" {
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: zh add [github-url]")
			os.Exit(1)
		}
		err := deps.AddDependency(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to add dependency: %v\n", err)
			os.Exit(1)
		}
		return
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
		fmt.Fprintf(os.Stderr, "Failed to execute zig: %v\n", err)
		os.Exit(1)
	}
}
