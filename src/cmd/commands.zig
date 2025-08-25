const std = @import("std");
const print = std.debug.print;
const process = std.process;
const deps = @import("../deps/deps.zig");
const project = @import("../project/project.zig");
const types = @import("../types.zig");

const VERSION = "0.2.3";

pub const Command = enum {
    add,
    new,
    version,
    passthrough,
};

pub fn printUsage() void {
    print("zh {s}\n", .{VERSION});
    print("A Cargo-like package manager and build tool wrapper for Zig\n\n", .{});
    print("Usage: zh [command] [arguments]\n\n", .{});
    print("Commands:\n", .{});
    print("  new [name] [options]    Create a new Zig project\n", .{});
    print("    Options:\n", .{});
    print("      --bin               Create a binary project (default)\n", .{});
    print("      --lib               Create a library project\n", .{});
    print("  add [repo] [options]    Add a dependency\n", .{});
    print("    Options:\n", .{});
    print("      --github, --gh      Use GitHub (default)\n", .{});
    print("      --gitlab, --gl      Use GitLab\n", .{});
    print("      --codeberg, --cb    Use Codeberg\n", .{});
    print("  version                 Show version information\n", .{});
    print("  [any zig command]       Pass arguments directly to the zig command\n\n", .{});
    print("Examples:\n", .{});
    print("  zh new myapp --bin      Create a new binary project\n", .{});
    print("  zh new mylib --lib      Create a new library project\n", .{});
    print("  zh add mitchellh/libxev Add a dependency from GitHub\n", .{});
    print("  zh build                Pass through to 'zig build'\n", .{});
}

pub fn parseCommand(args: [][:0]u8) Command {
    if (args.len < 2) return .passthrough;

    const cmd = args[1];
    if (std.mem.eql(u8, cmd, "add")) return .add;
    if (std.mem.eql(u8, cmd, "new")) return .new;
    if (std.mem.eql(u8, cmd, "version") or std.mem.eql(u8, cmd, "--version") or std.mem.eql(u8, cmd, "-v")) return .version;

    return .passthrough;
}

pub fn handleNewCommand(allocator: std.mem.Allocator, args: [][:0]u8) !void {
    if (args.len < 3) {
        std.debug.print("Usage: zh new [name] [options]\n", .{});
        std.debug.print("  Options:\n", .{});
        std.debug.print("    --bin    Create a binary project (default)\n", .{});
        std.debug.print("    --lib    Create a library project\n", .{});
        process.exit(1);
    }

    const project_name = args[2];
    var project_type = types.ProjectType.binary;

    for (args[3..]) |arg| {
        if (std.mem.eql(u8, arg, "--bin")) {
            project_type = .binary;
        } else if (std.mem.eql(u8, arg, "--lib")) {
            project_type = .library;
        } else {
            std.debug.print("Unknown option: {s}\n", .{arg});
            process.exit(1);
        }
    }

    try project.createProject(allocator, project_name, project_type);
}

pub fn handleAddCommand(allocator: std.mem.Allocator, args: [][:0]u8) !void {
    if (args.len < 3) {
        std.debug.print("Usage: zh add [repo] [options]\n", .{});
        std.debug.print("  Options:\n", .{});
        std.debug.print("    --github or --gh      Use GitHub (default)\n", .{});
        std.debug.print("    --gitlab or --gl      Use GitLab\n", .{});
        std.debug.print("    --codeberg or --cb    Use Codeberg\n", .{});
        process.exit(1);
    }

    const repo_url = args[2];
    var provider = types.GitProvider.github;

    for (args[3..]) |arg| {
        if (std.mem.eql(u8, arg, "--github") or std.mem.eql(u8, arg, "--gh")) {
            provider = .github;
        } else if (std.mem.eql(u8, arg, "--gitlab") or std.mem.eql(u8, arg, "--gl")) {
            provider = .gitlab;
        } else if (std.mem.eql(u8, arg, "--codeberg") or std.mem.eql(u8, arg, "--cb")) {
            provider = .codeberg;
        }
    }

    try deps.addDependency(allocator, repo_url, provider);
}

pub fn handleVersionCommand(allocator: std.mem.Allocator) !void {
    var stdout_buffer: [1024]u8 = undefined;
    var stdout_writer = std.fs.File.stdout().writer(&stdout_buffer);
    const stdout = &stdout_writer.interface;

    try stdout.print("zh {s}\n", .{VERSION});

    var child = std.process.Child.init(&[_][]const u8{ "zig", "version" }, allocator);
    child.stdout_behavior = .Pipe;
    child.stderr_behavior = .Pipe;

    try child.spawn();

    const zig_stdout = try child.stdout.?.readToEndAlloc(allocator, 1024);
    defer allocator.free(zig_stdout);

    const result = try child.wait();
    if (result == .Exited and result.Exited == 0) {
        try stdout.print("zig {s}", .{zig_stdout});
    } else {
        try stdout.print("zig: (not found or error)\n", .{});
    }

    try stdout.flush();
}

pub fn passthroughToZig(allocator: std.mem.Allocator, args: [][:0]u8) void {
    const zig_args = allocator.alloc([]const u8, args.len) catch {
        std.debug.print("Error: Failed to allocate memory for zig args\n", .{});
        process.exit(1);
    };
    defer allocator.free(zig_args);

    zig_args[0] = "zig";
    for (args[1..], 1..) |arg, i| {
        zig_args[i] = @as([]const u8, arg);
    }

    var child = std.process.Child.init(zig_args, allocator);
    child.stdin_behavior = .Inherit;
    child.stdout_behavior = .Inherit;
    child.stderr_behavior = .Inherit;

    const result = child.spawnAndWait() catch |err| {
        std.debug.print("Error spawning zig command: {}\n", .{err});
        std.debug.print("Command: zig", .{});
        for (args[1..]) |arg| {
            std.debug.print(" {s}", .{arg});
        }
        std.debug.print("\n", .{});
        std.debug.print("Make sure 'zig' is in your PATH\n", .{});
        process.exit(1);
    };

    switch (result) {
        .Exited => |code| process.exit(code),
        .Signal, .Stopped, .Unknown => process.exit(1),
    }
}
