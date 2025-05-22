package deps

import (
	"fmt"
)

func AddDependency(url string, provider GitProvider) error {
	repo, err := parseRepoURL(url, provider)
	if err != nil {
		return fmt.Errorf("failed to parse repository URL: %w", err)
	}

	if err := getLatestCommitInfo(repo); err != nil {
		return fmt.Errorf("failed to get repository information: %w", err)
	}

	if err := fetchDependency(repo); err != nil {
		return fmt.Errorf("failed to fetch dependency: %w", err)
	}

	dependencyName := generateDependencyName(repo)
	if err := addToBuildZig(dependencyName); err != nil {
		fmt.Printf("Warning: Failed to add dependency to build.zig: %v\n", err)
		fmt.Printf("You may need to manually add the dependency to your build.zig\n")
		fmt.Printf("Add this line after your exe_mod creation:\n")
		fmt.Printf(`    exe_mod.addImport("%s", b.dependency("%s", .{}).module("root"));`+"\n", dependencyName, dependencyName)
	} else {
		fmt.Printf("✓ Added dependency '%s' to build.zig\n", dependencyName)
	}

	fmt.Printf("✓ Dependency '%s' added successfully!\n", dependencyName)
	fmt.Printf("You can now use it in your Zig code:\n")
	fmt.Printf(`    const %s = @import("%s");`+"\n", dependencyName, dependencyName)

	return nil
}
