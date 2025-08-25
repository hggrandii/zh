const std = @import("std");
const ArrayList = std.ArrayList;

pub fn addToBuildZig(allocator: std.mem.Allocator, dependency_name: []const u8, module_name: []const u8) !void {
    const build_zig_content = try std.fs.cwd().readFileAlloc(allocator, "build.zig", 8192);
    defer allocator.free(build_zig_content);

    const search_str = try std.fmt.allocPrint(allocator, "\"{s}\"", .{dependency_name});
    defer allocator.free(search_str);

    if (std.mem.indexOf(u8, build_zig_content, search_str)) |_| {
        return error.DependencyAlreadyExists;
    }

    const modified = try injectDependency(allocator, build_zig_content, dependency_name, module_name);
    defer allocator.free(modified);

    try std.fs.cwd().writeFile(.{ .sub_path = "build.zig", .data = modified });
}

fn injectDependency(allocator: std.mem.Allocator, content: []const u8, dependency_name: []const u8, module_name: []const u8) ![]u8 {
    var lines = std.mem.splitScalar(u8, content, '\n');
    var result: ArrayList([]const u8) = .empty;
    defer result.deinit(allocator);

    var allocated_strings: ArrayList([]const u8) = .empty;
    defer {
        for (allocated_strings.items) |str| {
            allocator.free(str);
        }
        allocated_strings.deinit(allocator);
    }

    const exe_pattern = "const exe = b.addExecutable(.{";
    const lib_pattern = "const lib = b.addLibrary(.{";

    while (lines.next()) |line| {
        try result.append(allocator, line);

        if (std.mem.indexOf(u8, line, exe_pattern)) |_| {
            while (lines.next()) |next_line| {
                try result.append(allocator, next_line);
                if (std.mem.indexOf(u8, next_line, "});")) |_| {
                    try result.append(allocator, "");

                    const comment = try std.fmt.allocPrint(allocator, "    // Added {s}", .{dependency_name});
                    try allocated_strings.append(allocator, comment);
                    try result.append(allocator, comment);

                    const dep_line = try std.fmt.allocPrint(allocator, "    const {s}_dep = b.dependency(\"{s}\", .{{}});", .{ dependency_name, dependency_name });
                    try allocated_strings.append(allocator, dep_line);
                    try result.append(allocator, dep_line);

                    const import_line = try std.fmt.allocPrint(allocator, "    exe.root_module.addImport(\"{s}\", {s}_dep.module(\"{s}\"));", .{ dependency_name, dependency_name, module_name });
                    try allocated_strings.append(allocator, import_line);
                    try result.append(allocator, import_line);
                    break;
                }
            }
        } else if (std.mem.indexOf(u8, line, lib_pattern)) |_| {
            while (lines.next()) |next_line| {
                try result.append(allocator, next_line);
                if (std.mem.indexOf(u8, next_line, "});")) |_| {
                    try result.append(allocator, "");

                    const comment = try std.fmt.allocPrint(allocator, "    // Added {s}", .{dependency_name});
                    try allocated_strings.append(allocator, comment);
                    try result.append(allocator, comment);

                    const dep_line = try std.fmt.allocPrint(allocator, "    const {s}_dep = b.dependency(\"{s}\", .{{}});", .{ dependency_name, dependency_name });
                    try allocated_strings.append(allocator, dep_line);
                    try result.append(allocator, dep_line);

                    const import_line = try std.fmt.allocPrint(allocator, "    lib.root_module.addImport(\"{s}\", {s}_dep.module(\"{s}\"));", .{ dependency_name, dependency_name, module_name });
                    try allocated_strings.append(allocator, import_line);
                    try result.append(allocator, import_line);
                    break;
                }
            }
        }
    }

    var total_len: usize = 0;
    for (result.items) |line| {
        total_len += line.len + 1;
    }

    const joined = try allocator.alloc(u8, total_len);
    var pos: usize = 0;
    for (result.items) |line| {
        @memcpy(joined[pos .. pos + line.len], line);
        pos += line.len;
        joined[pos] = '\n';
        pos += 1;
    }

    return joined;
}
