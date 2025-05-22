package deps

import (
	"fmt"
	"os"
	"strings"
)

func addToBuildZig(dependencyName string) error {
	buildZigPath := "build.zig"
	content, err := os.ReadFile(buildZigPath)
	if err != nil {
		return fmt.Errorf("failed to read build.zig: %w", err)
	}

	buildZigContent := string(content)

	if strings.Contains(buildZigContent, fmt.Sprintf(`"%s"`, dependencyName)) {
		return fmt.Errorf("dependency '%s' may already exist in build.zig", dependencyName)
	}

	modified, err := injectDependency(buildZigContent, dependencyName)
	if err != nil {
		return err
	}

	if err := os.WriteFile(buildZigPath, []byte(modified), 0644); err != nil {
		return fmt.Errorf("failed to write build.zig: %w", err)
	}

	return nil
}

func injectDependency(buildZigContent, dependencyName string) (string, error) {
	lines := strings.Split(buildZigContent, "\n")
	var result []string

	for i, line := range lines {
		result = append(result, line)

		if strings.Contains(line, "const exe_mod = b.createModule(.{") {
			for j := i + 1; j < len(lines); j++ {
				result = append(result, lines[j])
				if strings.Contains(lines[j], "});") {
					result = append(result, "")
					result = append(result, fmt.Sprintf("    // Add %s dependency", dependencyName))
					result = append(result, fmt.Sprintf(`    exe_mod.addImport("%s", b.dependency("%s", .{}).module("root"));`, dependencyName, dependencyName))

					for k := j + 1; k < len(lines); k++ {
						result = append(result, lines[k])
					}

					return strings.Join(result, "\n"), nil
				}
			}
		}
	}

	return "", fmt.Errorf("could not find exe_mod creation pattern in build.zig")
}

func generateDependencyName(repo *RepoInfo) string {
	return repo.Repo
}
