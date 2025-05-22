package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hggrandii/zh/internal/deps"
	"github.com/hggrandii/zh/internal/project"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		handleAddCommand()
		return
	case "new":
		handleNewCommand()
		return
	default:
		// Pass through to zig command
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
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: zh [command] [arguments]")
	fmt.Fprintln(os.Stderr, "\nCommands:")
	fmt.Fprintln(os.Stderr, "  new [name] [options]    Create a new Zig project")
	fmt.Fprintln(os.Stderr, "    Options:")
	fmt.Fprintln(os.Stderr, "      --bin               Create a binary project (default)")
	fmt.Fprintln(os.Stderr, "      --lib               Create a library project")
	fmt.Fprintln(os.Stderr, "  add [repo] [options]    Add a dependency")
	fmt.Fprintln(os.Stderr, "    Options:")
	fmt.Fprintln(os.Stderr, "      --github, --gh      Use GitHub (default)")
	fmt.Fprintln(os.Stderr, "      --gitlab, --gl      Use GitLab")
	fmt.Fprintln(os.Stderr, "      --codeberg, --cb    Use Codeberg")
	fmt.Fprintln(os.Stderr, "  [any zig command]       Pass arguments directly to the zig command")
}

func handleNewCommand() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: zh new [name] [options]")
		fmt.Fprintln(os.Stderr, "  Options:")
		fmt.Fprintln(os.Stderr, "    --bin    Create a binary project (default)")
		fmt.Fprintln(os.Stderr, "    --lib    Create a library project")
		os.Exit(1)
	}

	projectName := os.Args[2]
	projectType := project.Binary

	for i := 3; i < len(os.Args); i++ {
		arg := strings.ToLower(os.Args[i])
		switch arg {
		case "--bin":
			projectType = project.Binary
		case "--lib":
			projectType = project.Library
		default:
			fmt.Fprintf(os.Stderr, "Unknown option: %s\n", arg)
			os.Exit(1)
		}
	}

	err := project.CreateProject(projectName, projectType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create project: %v\n", err)
		os.Exit(1)
	}
}

func handleAddCommand() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: zh add [repo] [options]")
		fmt.Fprintln(os.Stderr, "  Options:")
		fmt.Fprintln(os.Stderr, "    --github or --gh      Use GitHub (default)")
		fmt.Fprintln(os.Stderr, "    --gitlab or --gl      Use GitLab")
		fmt.Fprintln(os.Stderr, "    --codeberg or --cb    Use Codeberg")
		os.Exit(1)
	}

	repoURL := os.Args[2]
	provider := deps.GitHub

	for i := 3; i < len(os.Args); i++ {
		arg := strings.ToLower(os.Args[i])
		switch arg {
		case "--github", "--gh":
			provider = deps.GitHub
		case "--gitlab", "--gl":
			provider = deps.GitLab
		case "--codeberg", "--cb":
			provider = deps.Codeberg
		}
	}

	err := deps.AddDependency(repoURL, provider)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to add dependency: %v\n", err)
		os.Exit(1)
	}
}
