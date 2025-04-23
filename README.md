# Repo Mapper

A command-line utility for creating a structured overview of your project in Markdown
format.

## Description

Repo Mapper scans the specified project directory, analyzes its structure, and creates a
Markdown document containing:

- Visual representation of the directory and file structure as a tree
- Contents of all project files (except for ignored files)

The utility automatically takes into account rules from .gitignore to filter files that
should not be included in the report.

## Installation

```bash
go install github.com/dsbasko/repo-mapper@latest
```

## Usage

### Basic usage

```bash
repo-mapper
```

This command will scan the current directory and create a `summary.md` file with an
overview of the project.

### With parameters

```bash
repo-mapper -dir /path/to/your/project -o report.md
```

With custom ignore patterns:

```bash
repo-mapper -ignore "*.log,tmp/*,!important.log"
```

### Available parameters

| Parameter        | Description                                                    | Default                 |
| ---------------- | -------------------------------------------------------------- | ----------------------- |
| `-dir`           | Root directory of the project for scanning                     | `.` (current directory) |
| `-o`             | Output file name                                               | `summary.md`            |
| `-ignore`, `-i`  | Additional ignore patterns (comma-separated, .gitignore syntax)| (none)                  |
| `-h`, `--help`   | Show help                                                      |                         |

## Using .gitignore

Repo Mapper automatically reads the `.gitignore` file from the project's root directory
and does not include files and directories that match Git ignore rules in the report.

## Custom Ignore Patterns

In addition to using `.gitignore` files, you can specify custom ignore patterns directly through command-line parameters:

```bash
repo-mapper -i "*.log,build/*,node_modules/"
```

These patterns follow the same syntax as .gitignore:

- Use `*` to match any number of characters (e.g., `*.log` ignores all log files)
- Use `!` to negate a pattern (e.g., `!important.log` includes this file even if other patterns would ignore it)
- Use `**` to match nested directories (e.g., `**/node_modules` ignores node_modules at any level)
- Use `/` at the end to specify directories (e.g., `tmp/` ignores only directories named tmp)

Custom ignore patterns are combined with patterns from the project's `.gitignore` file and default ignore patterns.

### Examples

Ignore all JavaScript files except those in a specific directory:
```bash
repo-mapper -i "*.js,!src/components/*.js"
```

Ignore build directories and temporary files:
```bash
repo-mapper -i "build/,dist/,tmp/,*.tmp"
```

Ignore multiple file types:
```bash
repo-mapper -i "*.log,*.tmp,*.cache"
```