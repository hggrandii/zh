const std = @import("std");
const print = std.debug.print;
const types = @import("../types.zig");
const resolver = @import("resolver.zig");
const buildzig = @import("buildzig.zig");
const fetcher = @import("fetcher.zig");

pub fn addDependency(allocator: std.mem.Allocator, url: []const u8, provider: types.GitProvider) !void {
    var repo = try resolver.parseRepoURL(allocator, url, provider);
    defer repo.deinit();

    try resolver.getLatestCommitInfo(allocator, &repo);
    try fetcher.fetchDependency(allocator, &repo);

    const dependency_name = repo.repo;
    buildzig.addToBuildZig(allocator, dependency_name) catch |err| {
        print("Warning: Failed to add dependency to build.zig: {}\n", .{err});
        print("You may need to manually add the dependency to your build.zig\n", .{});
        print("Add this line after your exe/lib creation:\n", .{});
        print("    exe.root_module.addImport(\"{s}\", b.dependency(\"{s}\", .{{}}).module(\"root\"));\n", .{ dependency_name, dependency_name });
        return;
    };

    print("✓ Added dependency '{s}' to build.zig\n", .{dependency_name});
    print("✓ Dependency '{s}' added successfully!\n", .{dependency_name});
    print("You can now use it in your Zig code:\n", .{});
    print("    const {s} = @import(\"{s}\");\n", .{ dependency_name, dependency_name });
}
