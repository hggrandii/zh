const std = @import("std");

pub const GitProvider = enum {
    github,
    gitlab,
    codeberg,

    pub fn toString(self: GitProvider) []const u8 {
        return switch (self) {
            .github => "github",
            .gitlab => "gitlab",
            .codeberg => "codeberg",
        };
    }

    pub fn domain(self: GitProvider) []const u8 {
        return switch (self) {
            .github => "github.com/",
            .gitlab => "gitlab.com/",
            .codeberg => "codeberg.org/",
        };
    }
};

pub const ProjectType = enum {
    binary,
    library,

    pub fn toString(self: ProjectType) []const u8 {
        return switch (self) {
            .binary => "binary",
            .library => "library",
        };
    }
};

pub const RepoInfo = struct {
    owner: []const u8,
    repo: []const u8,
    provider: GitProvider,
    commit_hash: ?[]const u8 = null,
    default_branch: ?[]const u8 = null,
    latest_release: ?[]const u8 = null,
    is_release: bool = false,
    allocator: std.mem.Allocator,

    pub fn deinit(self: *RepoInfo) void {
        self.allocator.free(self.owner);
        self.allocator.free(self.repo);
        if (self.commit_hash) |hash| self.allocator.free(hash);
        if (self.default_branch) |branch| self.allocator.free(branch);
        if (self.latest_release) |release| self.allocator.free(release);
    }
};
