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
go install github.com/dsbasko/repo-mapper
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

### Available parameters

| Parameter      | Description                                | Default                 |
| -------------- | ------------------------------------------ | ----------------------- |
| `-dir`         | Root directory of the project for scanning | `.` (current directory) |
| `-o`           | Output file name                           | `summary.md`            |
| `-h`, `--help` | Show help                                  |                         |

## Using .gitignore

Repo Mapper automatically reads the `.gitignore` file from the project's root directory
and does not include files and directories that match Git ignore rules in the report.
