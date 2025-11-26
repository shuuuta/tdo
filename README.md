# Pask

A simple command-line task manager for projects and global tasks.

Pask (short for "project task") helps you organize tasks both globally and per project. Tasks are automatically associated with the current Git repository, or you can manage global tasks that are available everywhere.

## Installation

```bash
go install github.com/shuuuta/pask@latest
```

## Usage

### Add tasks

```bash
# Add project tasks (in a Git repository)
pask add "Fix authentication bug" "Update documentation"

# Add global tasks
pask add -g "Buy groceries" "Call dentist"
```

### List tasks

```bash
# List project tasks
pask list

# List global tasks
pask list -g
```

### Mark tasks as done

```bash
# Mark project tasks as done (by index)
pask done 1 3

# Mark global tasks as done
pask done -g 2
```

## How it works

- **Project tasks**: When in a Git repository, tasks are automatically associated with that project
- **Global tasks**: Use the `-g` flag to manage tasks that are available everywhere
- **Automatic fallback**: If not in a Git repository, tasks are automatically saved as global tasks

## License

MIT
