const std = @import("std");
const process = std.process;
const cmd = @import("cmd/commands.zig");

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    defer _ = gpa.deinit();
    const allocator = gpa.allocator();

    const args = try process.argsAlloc(allocator);
    defer process.argsFree(allocator, args);

    if (args.len < 2) {
        cmd.printUsage();
        process.exit(1);
    }

    const command = cmd.parseCommand(args);

    switch (command) {
        .add => try cmd.handleAddCommand(allocator, args),
        .new => try cmd.handleNewCommand(allocator, args),
        .version => try cmd.handleVersionCommand(allocator),
        .passthrough => cmd.passthroughToZig(allocator, args),
    }
}
