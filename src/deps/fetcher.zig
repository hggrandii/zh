const std = @import("std");
const print = std.debug.print;
const types = @import("../types.zig");

pub fn fetchDependency(allocator: std.mem.Allocator, repo: *types.RepoInfo) !void {
    const archive_url = try generateArchiveURL(allocator, repo);
    defer allocator.free(archive_url);

    print("Generated archive URL: {s}\n", .{archive_url});

    var child = std.process.Child.init(&[_][]const u8{ "zig", "fetch", "--save", archive_url }, allocator);
    child.stdout_behavior = .Inherit;
    child.stderr_behavior = .Inherit;

    const result = try child.spawnAndWait();
    if (result != .Exited or result.Exited != 0) {
        return error.FetchFailed;
    }
}

fn generateArchiveURL(allocator: std.mem.Allocator, repo: *types.RepoInfo) ![]u8 {
    const is_tag = std.mem.startsWith(u8, repo.commit_hash.?, "v") or
        std.mem.indexOf(u8, repo.commit_hash.?, ".") != null;

    return switch (repo.provider) {
        .github => if (is_tag)
            try std.fmt.allocPrint(allocator, "https://github.com/{s}/{s}/archive/refs/tags/{s}.tar.gz", .{ repo.owner, repo.repo, repo.commit_hash.? })
        else
            try std.fmt.allocPrint(allocator, "https://github.com/{s}/{s}/archive/{s}.tar.gz", .{ repo.owner, repo.repo, repo.commit_hash.? }),

        .gitlab => if (is_tag)
            try std.fmt.allocPrint(allocator, "https://gitlab.com/{s}/{s}/-/archive/{s}/{s}-{s}.tar.gz", .{ repo.owner, repo.repo, repo.commit_hash.?, repo.repo, repo.commit_hash.? })
        else
            try std.fmt.allocPrint(allocator, "https://gitlab.com/{s}/{s}/-/archive/{s}/{s}-{s}.tar.gz", .{ repo.owner, repo.repo, repo.commit_hash.?, repo.repo, repo.commit_hash.? }),

        .codeberg => try std.fmt.allocPrint(allocator, "https://codeberg.org/{s}/{s}/archive/{s}.tar.gz", .{ repo.owner, repo.repo, repo.commit_hash.? }),
    };
}
