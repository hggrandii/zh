package project

import (
	"fmt"
	"os"
	"os/exec"
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

	defer func() {
		os.Chdir(originalDir)
	}()

	cmd := exec.Command("zig", "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run zig init: %w", err)
	}

	if projectType == Binary {
		if err := convertToBinary(name); err != nil {
			return fmt.Errorf("failed to convert to binary project: %w", err)
		}
	}

	fmt.Printf("Created %s project '%s'\n", getProjectTypeString(projectType), name)
	return nil
}

func convertToBinary(projectName string) error {
	if err := os.Remove("src/root.zig"); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove src/root.zig: %w", err)
	}

	buildZigContent := generateBinaryBuildZig(projectName)
	if err := os.WriteFile("build.zig", []byte(buildZigContent), 0644); err != nil {
		return fmt.Errorf("failed to write build.zig: %w", err)
	}

	return nil
}

func generateBinaryBuildZig(projectName string) string {
	return fmt.Sprintf(`const std = @import("std");

pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});

    const optimize = b.standardOptimizeOption(.{});

    const exe_mod = b.createModule(.{
        .root_source_file = b.path("src/main.zig"),
        .target = target,
        .optimize = optimize,
    });

    const exe = b.addExecutable(.{
        .name = "%s",
        .root_module = exe_mod,
    });

    b.installArtifact(exe);

    const run_cmd = b.addRunArtifact(exe);

    run_cmd.step.dependOn(b.getInstallStep());

    if (b.args) |args| {
        run_cmd.addArgs(args);
    }

    const run_step = b.step("run", "Run the app");
    run_step.dependOn(&run_cmd.step);

    const exe_unit_tests = b.addTest(.{
        .root_module = exe_mod,
    });

    const run_exe_unit_tests = b.addRunArtifact(exe_unit_tests);

    const test_step = b.step("test", "Run unit tests");
    test_step.dependOn(&run_exe_unit_tests.step);
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
