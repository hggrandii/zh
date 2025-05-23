package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hggrandii/zh/internal/deps"
	"github.com/hggrandii/zh/internal/project"
)

const VERSION = "0.1.2"

func HandleNewCommand() {
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

func HandleAddCommand() {
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

func HandleVersionCommand() {
	fmt.Printf("zh %s\n", VERSION)

	cmd := exec.Command("zig", "version")
	zigVersion, err := cmd.Output()
	if err != nil {
		fmt.Printf("zig: (not found or error: %v)\n", err)
	} else {
		fmt.Printf("zig %s", string(zigVersion))
	}
}
