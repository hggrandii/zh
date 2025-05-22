package deps

import (
	"fmt"
	"os"
	"os/exec"
)

func fetchDependency(repo *RepoInfo) error {
	archiveURL := generateArchiveURL(repo)
	fmt.Printf("Generated archive URL: %s\n", archiveURL)

	cmd := exec.Command("zig", "fetch", "--save", archiveURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to fetch dependency: %w", err)
	}

	return nil
}
