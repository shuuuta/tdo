# Tdo

A simple command-line todo manager for projects and global tasks.

Tdo helps you organize tasks both globally and per project. Tasks are automatically associated with the current Git repository, or you can manage global tasks that are available everywhere.

## Installation

```bash
go install github.com/shuuuta/tdo@latest
```

## Usage

### Add tasks

```bash
# Add project tasks (in a Git repository)
tdo add "Fix authentication bug" "Update documentation"

# Add global tasks
tdo add -g "Buy groceries" "Call dentist"
```

### List tasks

```bash
# List project tasks
tdo list

# List global tasks
tdo list -g
```

### Mark tasks as done

```bash
# Mark project tasks as done (by index)
tdo done 1 3

# Mark global tasks as done
tdo done -g 2
```

## How it works

- **Project tasks**: When in a Git repository, tasks are automatically associated with that project
- **Global tasks**: Use the `-g` flag to manage tasks that are available everywhere
- **Automatic fallback**: If not in a Git repository, tasks are automatically saved as global tasks

## License

MIT
