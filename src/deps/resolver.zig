const std = @import("std");
const types = @import("../types.zig");

pub fn parseRepoURL(allocator: std.mem.Allocator, url: []const u8, provider: types.GitProvider) !types.RepoInfo {
    var clean_url = url;

    if (std.mem.startsWith(u8, clean_url, "https://")) {
        clean_url = clean_url[8..];
    } else if (std.mem.startsWith(u8, clean_url, "http://")) {
        clean_url = clean_url[7..];
    }

    if (std.mem.indexOf(u8, clean_url, ".") == null) {
        var parts = std.mem.splitScalar(u8, clean_url, '/');
        const owner = parts.next() orelse return error.InvalidFormat;
        const repo = parts.next() orelse return error.InvalidFormat;

        return types.RepoInfo{
            .owner = owner,
            .repo = repo,
            .provider = provider,
            .allocator = allocator,
        };
    }

    const domain = provider.domain();
    if (std.mem.startsWith(u8, clean_url, domain)) {
        clean_url = clean_url[domain.len..];
    }

    var parts = std.mem.splitScalar(u8, clean_url, '/');
    const owner = parts.next() orelse return error.InvalidFormat;
    var repo = parts.next() orelse return error.InvalidFormat;

    if (std.mem.endsWith(u8, repo, ".git")) {
        repo = repo[0 .. repo.len - 4];
    }

    return types.RepoInfo{
        .owner = owner,
        .repo = repo,
        .provider = provider,
        .allocator = allocator,
    };
}

pub fn getLatestCommitInfo(allocator: std.mem.Allocator, repo: *types.RepoInfo) !void {
    if (getLatestStableRelease(allocator, repo)) {
        return;
    } else |_| {
        std.debug.print("// Could not find a stable release, using latest commit\n", .{});
        try getLatestCommitFromBranch(allocator, repo);
    }
}

fn getLatestStableRelease(allocator: std.mem.Allocator, repo: *types.RepoInfo) !void {
    const url = switch (repo.provider) {
        .github => try std.fmt.allocPrint(allocator, "https://api.github.com/repos/{s}/{s}/releases/latest", .{ repo.owner, repo.repo }),
        .gitlab => try std.fmt.allocPrint(allocator, "https://gitlab.com/api/v4/projects/{s}%2F{s}/releases", .{ repo.owner, repo.repo }),
        .codeberg => try std.fmt.allocPrint(allocator, "https://codeberg.org/api/v1/repos/{s}/{s}/releases", .{ repo.owner, repo.repo }),
    };
    defer allocator.free(url);

    const response = makeHttpRequest(allocator, url) catch {
        return error.NoStableRelease;
    };
    defer allocator.free(response);

    const parsed = std.json.parseFromSlice(std.json.Value, allocator, response, .{}) catch {
        return error.NoStableRelease;
    };
    defer parsed.deinit();

    switch (repo.provider) {
        .github => {
            const tag_name = parsed.value.object.get("tag_name") orelse return error.NoStableRelease;
            const prerelease = parsed.value.object.get("prerelease") orelse return error.NoStableRelease;

            if (prerelease.bool) return error.NoStableRelease;

            repo.commit_hash = try allocator.dupe(u8, tag_name.string);
        },
        .gitlab, .codeberg => {
            const releases = parsed.value.array;
            if (releases.items.len == 0) return error.NoStableRelease;

            for (releases.items) |release| {
                const tag_name = release.object.get("tag_name") orelse continue;

                if (repo.provider == .gitlab) {
                    if (release.object.get("upcoming_release")) |upcoming| {
                        if (upcoming.bool) continue;
                    }
                }

                repo.commit_hash = try allocator.dupe(u8, tag_name.string);
                return;
            }
            return error.NoStableRelease;
        },
    }
}

fn getLatestCommitFromBranch(allocator: std.mem.Allocator, repo: *types.RepoInfo) !void {
    try getDefaultBranch(allocator, repo);

    const url = switch (repo.provider) {
        .github => try std.fmt.allocPrint(allocator, "https://api.github.com/repos/{s}/{s}/branches/{s}", .{ repo.owner, repo.repo, repo.default_branch.? }),
        .gitlab => try std.fmt.allocPrint(allocator, "https://gitlab.com/api/v4/projects/{s}%2F{s}/repository/branches/{s}", .{ repo.owner, repo.repo, repo.default_branch.? }),
        .codeberg => try std.fmt.allocPrint(allocator, "https://codeberg.org/api/v1/repos/{s}/{s}/branches/{s}", .{ repo.owner, repo.repo, repo.default_branch.? }),
    };
    defer allocator.free(url);

    const response = try makeHttpRequest(allocator, url);
    defer allocator.free(response);

    const parsed = try std.json.parseFromSlice(std.json.Value, allocator, response, .{});
    defer parsed.deinit();

    const commit_obj = switch (repo.provider) {
        .github => parsed.value.object.get("commit").?.object,
        .gitlab, .codeberg => parsed.value.object.get("commit").?.object,
    };

    const commit_hash = switch (repo.provider) {
        .github => commit_obj.get("sha").?.string,
        .gitlab, .codeberg => commit_obj.get("id").?.string,
    };

    repo.commit_hash = try allocator.dupe(u8, commit_hash);
}

fn getDefaultBranch(allocator: std.mem.Allocator, repo: *types.RepoInfo) !void {
    const url = switch (repo.provider) {
        .github => try std.fmt.allocPrint(allocator, "https://api.github.com/repos/{s}/{s}", .{ repo.owner, repo.repo }),
        .gitlab => try std.fmt.allocPrint(allocator, "https://gitlab.com/api/v4/projects/{s}%2F{s}", .{ repo.owner, repo.repo }),
        .codeberg => try std.fmt.allocPrint(allocator, "https://codeberg.org/api/v1/repos/{s}/{s}", .{ repo.owner, repo.repo }),
    };
    defer allocator.free(url);

    const response = try makeHttpRequest(allocator, url);
    defer allocator.free(response);

    const parsed = try std.json.parseFromSlice(std.json.Value, allocator, response, .{});
    defer parsed.deinit();

    const default_branch = parsed.value.object.get("default_branch").?.string;
    repo.default_branch = try allocator.dupe(u8, default_branch);
}

fn makeHttpRequest(allocator: std.mem.Allocator, url: []const u8) ![]u8 {
    var client = std.http.Client{ .allocator = allocator };
    defer client.deinit();

    const uri = try std.Uri.parse(url);

    const header_buffer = try allocator.alloc(u8, 8192);
    defer allocator.free(header_buffer);
    var req = try client.open(.GET, uri, .{
        .server_header_buffer = header_buffer,
    });
    defer req.deinit();

    try req.send();
    try req.finish();
    try req.wait();

    const body = try req.reader().readAllAlloc(allocator, 1024 * 1024);
    return body;
}
