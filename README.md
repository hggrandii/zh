# zh

A simple hobby package manager and build tool wrapper for Zig.

## What it does

- Creates new Zig projects with sensible defaults
- Adds dependencies from GitHub/GitLab/Codeberg with automatic module detection
- Wraps `zig build` commands for convenience
- Prefers stable releases over latest commits when available

## Installation

### Option 1: Build from source
```bash
git clone https://github.com/hggrandii/zh
cd zh
zig build -Doptimize=ReleaseSafe
sudo cp zig-out/bin/zh /usr/local/bin/
```

### Option 2: User install (no sudo)
```bash
git clone https://github.com/hggrandii/zh
cd zh
zig build -Doptimize=ReleaseSafe
mkdir -p ~/.local/bin
cp zig-out/bin/zh ~/.local/bin/
```

## Usage

### Create a new project
```bash
zh new myproject          # Creates a binary project
zh new mylib --lib        # Creates a library project
```

### Add dependencies
```bash
zh add zigzap/zap         # Adds latest stable release
zh add user/repo --gitlab # From GitLab
zh add user/repo --cb     # From Codeberg
```

### Build and run
```bash
zh build                  # Same as 'zig build'
zh run                    # Same as 'zig run'
zh test                   # Same as 'zig test'
# ... any other zig command
```

### Version info
```bash
zh version               # Shows zh and zig versions
```

## Features

- **Smart dependency resolution**: Automatically detects the correct module name by parsing the package's build files
- **Stable releases first**: Prefers tagged releases over latest commits when available
- **Simple project templates**: Clean, minimal project structure
- **Multiple git providers**: GitHub, GitLab, and Codeberg support

## Example

```bash
zh new myapp
cd myapp
zh add zigzap/zap
zh build
zh run
```

## Notes

This is a hobby project. It's not trying to replace the official Zig package manager or compete with serious tools. It just makes some common workflows a bit more convenient.

## Version

Current version: 0.2.0
