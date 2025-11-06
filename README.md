# cursync

A CLI tool for synchronizing Cursor IDE rules between a source directory and git project's `.cursor/rules` directory. Supports bidirectional sync, file pattern filtering, YAML header preservation, and configuration management.

## How to start

### Prerequisites

- Go>=1.22.4
- Git

### Install
```bash
go install github.com/yanodintsovmercuryo/cursync@latest
```

### Example cxonfiguration
```bash 
# Create rules dir
mkdir ~/my-rules
cd ~/my-rules
git init

# Set default configs
cursync cfg -d ~/my-rules -p "local_*.mdc" -o false -w false

# Upload rules from exist project
cd ~/dev/exist-project
cursync push

# Download rules to new project
cd ~/dev/new-project
cursync pull
```

## Commands

### pull

```bash
cursync pull -d ~/my-rules -p "local-*.mdc" -o false
```

Synchronizes files from source directory to project `.cursor/rules` directory. Deletes extra files in project that don't exist in source. 

Flags:

- **`--rules-dir` / `-d`** - Path to rules directory (overrides config file)
- **`--file-patterns` / `-p`** - Comma-separated file patterns to sync (e.g., `local_*.mdc,translate/*.md`) (overrides config file)
- **`--overwrite-headers` / `-o`** - Overwrite headers instead of preserving them


### push

```bash
cursync push -d ~/my-rules -p "local_*.mdc" -o false -w true
```

Synchronizes files from project `.cursor/rules` directory to source directory. Deletes extra files in source that don't exist in project. Automatically commits changes to git repository. 

Flags:

- **`--rules-dir` / `-d`** - Path to rules directory (overrides config file)
- **`--file-patterns` / `-p`** - Comma-separated file patterns to sync (e.g., `local_*.mdc,translate/*.md`) (overrides config file)
- **`--overwrite-headers` / `-o`** - Overwrite headers instead of preserving them
- **`--git-without-push` / `-w`** - Commit changes but don't push to remote

### cfg

```bash
# View current configuration
cursync cfg

# Set default configuration values
cursync cfg -d ~/my-rules -p "local_*.mdc" -o false -w false
```

Manages default configuration values stored in `~/.config/cursync.toml`:
- View configuration: Run `cursync cfg` without flags
- Set defaults: Use flags to set default values
- Clear defaults: Set empty value for string flags or `false` for bool flags

Flags:

- **`--rules-dir` / `-d`** - Set default rules directory (empty value clears it)
- **`--file-patterns` / `-p`** - Set default file patterns (empty value clears it)
- **`--overwrite-headers` / `-o`** - Set default overwrite-headers flag (use `true`, `1`, `false`, `0`, or empty to clear)
- **`--git-without-push` / `-w`** - Set default git-without-push flag (use `true`, `1`, `false`, `0`, or empty to clear)


## Development

### Prerequisites

- Go 1.22.4+
- Task runner (for using Taskfile.yml)

### Building

```bash
task build
```

### Running Tests

```bash
task test
```

### Formatting and Linting

```bash
task fmt
task lint
```

### Generating Mocks

```bash
task generate
```

## Architecture

The tool follows a clean architecture pattern:

- **`pkg/`** - Static utilities without dependencies (file operations, path utilities, git operations, output formatting)
- **`service/`** - Business logic with dependencies:
  - **`service/file/`** - File operations facade (comparator, copier, filter sub-services)
  - **`service/sync/`** - Main synchronization service orchestrating pull/push operations
- **`models/`** - Data structures and types

### Synchronization Flow

1. **Pull Flow:**
   - Get rules source directory from flag or config
   - Detect git root directory
   - Find source files (with optional pattern filtering)
   - Clean up extra files in destination
   - Copy files maintaining directory structure
   - Skip identical files

2. **Push Flow:**
   - Get rules source directory from flag or config
   - Detect git root directory
   - Verify project `.cursor/rules` directory exists
   - Find project files (with optional pattern filtering)
   - Clean up extra files in source directory
   - Copy files maintaining directory structure
   - Commit changes to git repository (with optional push)

## Troubleshooting

### "rules directory not specified" error

Ensure `--rules-dir` flag is provided or set default via `cursync cfg --rules-dir <path>`.

### "failed to find git root" error

The tool searches recursively for either a `.git` directory or `.cursor` folder starting from the current directory. Ensure you're running the command from within a git repository or a directory containing a `.cursor` folder.

### "project rules directory not found" error (push command)

The push command requires the project's `.cursor/rules` directory to exist. Create it first if needed.

### Git commit failures

Check git repository status and ensure you have proper permissions. The tool will continue synchronization even if commit fails, but will display an error message.
