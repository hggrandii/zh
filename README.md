# zh - A Cargo-like Package Manager for Zig

`zh` is a wrapper around the Zig toolchain that provides Cargo-like functionality for creating projects and managing dependencies.

## Features

- **Easy project creation**: `zh new myproject --bin` or `zh new mylib --lib`
- **Dependency management**: `zh add owner/repo` to add dependencies from GitHub, GitLab, or Codeberg
- **Zig passthrough**: Any unrecognized command is passed directly to `zig`
- **Clean templates**: Minimal, focused project templates without excessive comments

## Installation

### Using Go (Recommended)

```bash
go install github.com/hggrandii/zh@latest
```

### From Source

```bash
git clone https://github.com/hggrandii/zh.git
cd zh
go install .
```

## Requirements

- Go 1.19+ (for installation)
- Zig 0.14.0+ (for project creation and building)

## Usage

### Create a new project

```bash
# Create a binary project (default)
zh new myapp
zh new myapp --bin

# Create a library project
zh new mylib --lib
```

### Add dependencies

```bash
# Add from GitHub (default)
zh add mitchellh/libxev
zh add owner/repo --github

# Add from GitLab
zh add owner/repo --gitlab

# Add from Codeberg
zh add owner/repo --codeberg
```

### Build and run (passthrough to zig)

```bash
zh build
zh build run
zh test
zh --help
```

### Show version information

```bash
zh version
```

## Project Structure

When you create a new project, `zh` generates a clean, minimal structure:

### Binary Project (`zh new myapp --bin`)

```
myapp/
├── build.zig          # Clean build script
├── build.zig.zon      # Dependencies manifest
└── src/
    └── main.zig       # Simple "Hello, world!"
```

### Library Project (`zh new mylib --lib`)

```
mylib/
├── build.zig          # Library build script
├── build.zig.zon      # Dependencies manifest
└── src/
    └── root.zig       # Simple add function with test
```

## Comparison with Cargo

| Cargo | zh | Description |
|-------|----|----|
| `cargo new myapp` | `zh new myapp --bin` | Create binary project |
| `cargo new mylib --lib` | `zh new mylib --lib` | Create library project |
| `cargo add dep` | `zh add owner/repo` | Add dependency |
| `cargo build` | `zh build` | Build project |
| `cargo run` | `zh build run` | Run project |
| `cargo test` | `zh test` | Run tests |

## Why zh?

Zig's package manager is powerful but can be verbose and complex for simple use cases. `zh` provides:

- **Familiar workflow** for developers coming from Rust/Cargo
- **Clean project templates** without excessive boilerplate
- **Simplified dependency management** with automatic build.zig integration
- **Transparent zig passthrough** - you can use any zig command

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.


## Acknowledgments

- Inspired by Rust's Cargo package manager
- Built on top of Zig's excellent build system and package manager
