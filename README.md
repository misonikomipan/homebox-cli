# homebox-cli

[日本語 (Japanese)](README.ja.md)

[![Go Report Card]
(https://goreportcard.com/badge/github.com/misonikomipan/homebox-cli)](https://goreportcard.com/report/github.com/misonikomipan/homebox-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A powerful, user-friendly Command-Line Interface (CLI) for managing your [Homebox](https://github.com/sysadminsmedia/homebox) inventory system.

## Features

- **Resource Management**: CRUD operations for Items, Locations, Tags, Maintenance, Notifiers, and Templates.
- **Custom Fields**: Full support for Item custom fields (`hb items fields`).
- **Labelmaker**: Manage labelmaker configurations.
- **Flexible Output**: Choose between `json` (for scripting) or `table` (for readability).
- **Shell Autocompletion**: Support for Bash, Zsh, Fish, and PowerShell.
- **Hierarchy Support**: View location trees with or without items.
- **Data Portability**: Export and Import inventory items via CSV.

## Installation

### From Source

Ensure you have [Go](https://go.dev/doc/install) 1.21 or later installed.

```bash
git clone https://github.com/misonikomipan/homebox-cli.git
cd homebox-cli
go build -o hb main.go
mv hb /usr/local/bin/ # Optional: move to a directory in your PATH
```

## Quick Start

### 1. Configure the Endpoint

Set your Homebox instance URL:

```bash
hb config --endpoint https://homebox.example.com
```

### 2. Login

Authenticate with your email and password:

```bash
hb login --email your-email@example.com
```

### 3. Basic Commands

```bash
# List items in a beautiful table
hb items list --format table

# Search for an item
hb items list --query "laptop" --format table

# View location tree
hb locations tree --with-items

# Add a custom field to an item
hb items fields add <item-id> --label "Serial Number" --value "XYZ-123"

# Generate shell completion
hb completion zsh > ~/.zshrc.d/_hb
```

## Usage

For detailed help on any command, use the `--help` flag:

```bash
hb --help
hb items --help
hb items create --help
```

## Configuration

Settings are stored in `~/.config/hb/config.json`.

You can also use environment variables:
- `HB_ENDPOINT`: API endpoint URL
- `HB_TOKEN`: Authentication token
- `HB_FORMAT`: Default output format (`json` or `table`)

## Development

### Git Hooks

We use pre-commit and pre-push hooks to ensure code quality.

```bash
# Hooks are automatically enabled if you run the following after cloning:
chmod +x scripts/hooks/*
git config core.hooksPath scripts/hooks
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
