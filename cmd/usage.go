package cmd

import (
	"fmt"
	"os"
)

func PrintUsage() {
	fmt.Fprintf(os.Stderr, "zh %s\n", VERSION)
	fmt.Fprintln(os.Stderr, "A Cargo-like package manager and build tool wrapper for Zig")
	fmt.Fprintln(os.Stderr, "\nUsage: zh [command] [arguments]")
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
	fmt.Fprintln(os.Stderr, "  version                 Show version information")
	fmt.Fprintln(os.Stderr, "  [any zig command]       Pass arguments directly to the zig command")
	fmt.Fprintln(os.Stderr, "\nExamples:")
	fmt.Fprintln(os.Stderr, "  zh new myapp --bin      Create a new binary project")
	fmt.Fprintln(os.Stderr, "  zh new mylib --lib      Create a new library project")
	fmt.Fprintln(os.Stderr, "  zh add mitchellh/libxev Add a dependency from GitHub")
	fmt.Fprintln(os.Stderr, "  zh build                Pass through to 'zig build'")
}
