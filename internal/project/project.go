package project

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ProjectType int

const (
	Binary ProjectType = iota
	Library
)

func CreateProject(name string, projectType ProjectType) error {
	if err := os.MkdirAll(name, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := os.Chdir(name); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := runZigInit(); err != nil {
		return fmt.Errorf("failed to run zig init: %w", err)
	}

	if err := createOurTemplates(name, projectType); err != nil {
		return fmt.Errorf("failed to create templates: %w", err)
	}

	if err := cleanBuildZigZon(name); err != nil {
		return fmt.Errorf("failed to clean build.zig.zon: %w", err)
	}

	fmt.Printf("     Created %s project `%s`\n", getProjectTypeString(projectType), name)

	return nil
}

func runZigInit() error {
	cmd := exec.Command("zig", "init")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func createOurTemplates(name string, projectType ProjectType) error {
	os.RemoveAll("src")
	if err := os.MkdirAll("src", 0755); err != nil {
		return err
	}

	if projectType == Binary {
		mainContent := `const std = @import("std");

pub fn main() !void {
    std.debug.print("Hello, world!\n", .{});
}
`
		if err := os.WriteFile("src/main.zig", []byte(mainContent), 0644); err != nil {
			return err
		}

		buildContent := generateBinaryBuildZig(name)
		if err := os.WriteFile("build.zig", []byte(buildContent), 0644); err != nil {
			return err
		}

	} else {
		rootContent := `const std = @import("std");

pub fn add(a: i32, b: i32) i32 {
    return a + b;
}

test "basic add functionality" {
    try std.testing.expect(add(3, 7) == 10);
}
`
		if err := os.WriteFile("src/root.zig", []byte(rootContent), 0644); err != nil {
			return err
		}

		buildContent := generateLibraryBuildZig(name)
		if err := os.WriteFile("build.zig", []byte(buildContent), 0644); err != nil {
			return err
		}
	}

	return nil
}

func cleanBuildZigZon(name string) error {
	content, err := os.ReadFile("build.zig.zon")
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var fingerprint string
	for _, line := range lines {
		if strings.Contains(line, ".fingerprint = ") {
			fingerprint = strings.TrimSpace(line)
			break
		}
	}

	if fingerprint == "" {
		return fmt.Errorf("could not find fingerprint in generated build.zig.zon")
	}

	cleanZon := fmt.Sprintf(`.{
    .name = .%s,
    .version = "0.0.1",
    %s
    .minimum_zig_version = "0.14.0",
    .dependencies = .{},
    .paths = .{
        "build.zig",
        "build.zig.zon",
        "src",
    },
}
`, name, fingerprint)

	return os.WriteFile("build.zig.zon", []byte(cleanZon), 0644)
}

func generateBinaryBuildZig(projectName string) string {
	return fmt.Sprintf(`const std = @import("std");

pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    const exe = b.addExecutable(.{
        .name = "%s",
        .root_source_file = b.path("src/main.zig"),
        .target = target,
        .optimize = optimize,
    });

    b.installArtifact(exe);

    const run_cmd = b.addRunArtifact(exe);
    run_cmd.step.dependOn(b.getInstallStep());

    if (b.args) |args| {
        run_cmd.addArgs(args);
    }

    const run_step = b.step("run", "Run the app");
    run_step.dependOn(&run_cmd.step);

    const unit_tests = b.addTest(.{
        .root_source_file = b.path("src/main.zig"),
        .target = target,
        .optimize = optimize,
    });

    const run_unit_tests = b.addRunArtifact(unit_tests);
    const test_step = b.step("test", "Run unit tests");
    test_step.dependOn(&run_unit_tests.step);
}
`, projectName)
}

func generateLibraryBuildZig(projectName string) string {
	return fmt.Sprintf(`const std = @import("std");

pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    const lib = b.addStaticLibrary(.{
        .name = "%s",
        .root_source_file = b.path("src/root.zig"),
        .target = target,
        .optimize = optimize,
    });

    b.installArtifact(lib);

    const unit_tests = b.addTest(.{
        .root_source_file = b.path("src/root.zig"),
        .target = target,
        .optimize = optimize,
    });

    const run_unit_tests = b.addRunArtifact(unit_tests);
    const test_step = b.step("test", "Run unit tests");
    test_step.dependOn(&run_unit_tests.step);
}
`, projectName)
}

func getProjectTypeString(projectType ProjectType) string {
	switch projectType {
	case Binary:
		return "binary"
	case Library:
		return "library"
	default:
		return "unknown"
	}
}
