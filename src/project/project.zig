const std = @import("std");
const print = std.debug.print;
const types = @import("../types.zig");

pub fn createProject(allocator: std.mem.Allocator, name: []const u8, project_type: types.ProjectType) !void {
    std.fs.cwd().makeDir(name) catch |err| switch (err) {
        error.PathAlreadyExists => {},
        else => return err,
    };

    var original_dir = try std.fs.cwd().openDir(".", .{});
    defer original_dir.close();

    var project_dir = try std.fs.cwd().openDir(name, .{});
    defer project_dir.close();

    try std.posix.fchdir(project_dir.fd);
    defer std.posix.fchdir(original_dir.fd) catch {};

    try runZigInit();

    try createTemplates(allocator, name, project_type);

    try cleanBuildZigZon(allocator, name);

    print("     Created {s} project `{s}`\n", .{ project_type.toString(), name });
}

fn runZigInit() !void {
    var child = std.process.Child.init(&[_][]const u8{ "zig", "init" }, std.heap.page_allocator);
    child.stdout_behavior = .Ignore;
    child.stderr_behavior = .Ignore;
    _ = try child.spawnAndWait();
}

fn createTemplates(allocator: std.mem.Allocator, name: []const u8, project_type: types.ProjectType) !void {
    std.fs.cwd().deleteTree("src") catch {};
    try std.fs.cwd().makeDir("src");

    switch (project_type) {
        .binary => {
            const main_content =
                \\const std = @import("std");
                \\
                \\pub fn main() !void {
                \\    std.debug.print("Hello, world!\n", .{});
                \\}
                \\
            ;
            try std.fs.cwd().writeFile(.{ .sub_path = "src/main.zig", .data = main_content });

            const build_content = try generateBinaryBuildZig(allocator, name);
            defer allocator.free(build_content);
            try std.fs.cwd().writeFile(.{ .sub_path = "build.zig", .data = build_content });
        },
        .library => {
            const root_content =
                \\const std = @import("std");
                \\
                \\pub fn add(a: i32, b: i32) i32 {
                \\    return a + b;
                \\}
                \\
                \\test "basic add functionality" {
                \\    try std.testing.expect(add(3, 7) == 10);
                \\}
                \\
            ;
            try std.fs.cwd().writeFile(.{ .sub_path = "src/root.zig", .data = root_content });

            const build_content = try generateLibraryBuildZig(allocator, name);
            defer allocator.free(build_content);
            try std.fs.cwd().writeFile(.{ .sub_path = "build.zig", .data = build_content });
        },
    }
}

fn generateBinaryBuildZig(allocator: std.mem.Allocator, project_name: []const u8) ![]u8 {
    return try std.fmt.allocPrint(allocator,
        \\const std = @import("std");
        \\
        \\pub fn build(b: *std.Build) void {{
        \\    const target = b.standardTargetOptions(.{{}});
        \\    const optimize = b.standardOptimizeOption(.{{}});
        \\
        \\    const exe = b.addExecutable(.{{
        \\        .name = "{s}",
        \\        .root_source_file = b.path("src/main.zig"),
        \\        .target = target,
        \\        .optimize = optimize,
        \\    }});
        \\
        \\    b.installArtifact(exe);
        \\
        \\    const run_cmd = b.addRunArtifact(exe);
        \\    run_cmd.step.dependOn(b.getInstallStep());
        \\
        \\    if (b.args) |args| {{
        \\        run_cmd.addArgs(args);
        \\    }}
        \\
        \\    const run_step = b.step("run", "Run the app");
        \\    run_step.dependOn(&run_cmd.step);
        \\
        \\    const unit_tests = b.addTest(.{{
        \\        .root_source_file = b.path("src/main.zig"),
        \\        .target = target,
        \\        .optimize = optimize,
        \\    }});
        \\
        \\    const run_unit_tests = b.addRunArtifact(unit_tests);
        \\    const test_step = b.step("test", "Run unit tests");
        \\    test_step.dependOn(&run_unit_tests.step);
        \\}}
        \\
    , .{project_name});
}

fn generateLibraryBuildZig(allocator: std.mem.Allocator, project_name: []const u8) ![]u8 {
    return try std.fmt.allocPrint(allocator,
        \\const std = @import("std");
        \\
        \\pub fn build(b: *std.Build) void {{
        \\    const target = b.standardTargetOptions(.{{}});
        \\    const optimize = b.standardOptimizeOption(.{{}});
        \\
        \\    const lib = b.addStaticLibrary(.{{
        \\        .name = "{s}",
        \\        .root_source_file = b.path("src/root.zig"),
        \\        .target = target,
        \\        .optimize = optimize,
        \\    }});
        \\
        \\    b.installArtifact(lib);
        \\
        \\    const unit_tests = b.addTest(.{{
        \\        .root_source_file = b.path("src/root.zig"),
        \\        .target = target,
        \\        .optimize = optimize,
        \\    }});
        \\
        \\    const run_unit_tests = b.addRunArtifact(unit_tests);
        \\    const test_step = b.step("test", "Run unit tests");
        \\    test_step.dependOn(&run_unit_tests.step);
        \\}}
        \\
    , .{project_name});
}

fn cleanBuildZigZon(allocator: std.mem.Allocator, name: []const u8) !void {
    const content = try std.fs.cwd().readFileAlloc(allocator, "build.zig.zon", 8192);
    defer allocator.free(content);

    var lines = std.mem.splitScalar(u8, content, '\n');
    var fingerprint: ?[]const u8 = null;

    while (lines.next()) |line| {
        if (std.mem.indexOf(u8, line, ".fingerprint = ")) |_| {
            fingerprint = std.mem.trim(u8, line, " \t\n\r");
            break;
        }
    }

    if (fingerprint == null) {
        return error.FingerprintNotFound;
    }

    const clean_zon = try std.fmt.allocPrint(allocator,
        \\.{{
        \\    .name = .{s},
        \\    .version = "0.0.1",
        \\    {s}
        \\    .minimum_zig_version = "0.14.0",
        \\    .dependencies = .{{}},
        \\    .paths = .{{
        \\        "build.zig",
        \\        "build.zig.zon",
        \\        "src",
        \\    }},
        \\}}
        \\
    , .{ name, fingerprint.? });
    defer allocator.free(clean_zon);

    try std.fs.cwd().writeFile(.{ .sub_path = "build.zig.zon", .data = clean_zon });
}
