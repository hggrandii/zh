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
    allocator: std.mem.Allocator,

    pub fn deinit(self: *RepoInfo) void {
        if (self.commit_hash) |hash| self.allocator.free(hash);
        if (self.default_branch) |branch| self.allocator.free(branch);
    }
};
