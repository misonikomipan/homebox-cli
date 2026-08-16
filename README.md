# homebox-cli

[日本語 (Japanese)](README.ja.md)

[![Go Report Card]
(https://goreportcard.com/badge/github.com/misonikomipan/homebox-cli)](https://goreportcard.com/report/github.com/misonikomipan/homebox-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A powerful, user-friendly Command-Line Interface (CLI) for managing your [Homebox](https://github.com/sysadminsmedia/homebox) inventory system.

**Compatible with Homebox v0.26.x** (the CLI targets the v0.26 "entities" API, where items and locations were merged into a single resource).

## Features

- **Resource Management**: CRUD operations for Items, Locations, Tags, Maintenance, Notifiers, Templates, and Entity Types.
- **Custom Fields**: Full support for custom fields on entities (`hb items fields`).
- **API Keys**: Create, list, and revoke API keys (`hb auth api-keys`) and use them directly (`hb auth token hb_...`).
- **Labelmaker**: Generate item / location / asset labels as PNG (`hb labelmaker get`).
- **Flexible Output**: Choose between `json` (for scripting) or `table` (for readability).
- **Shell Autocompletion**: Support for Bash, Zsh, Fish, and PowerShell.
- **Hierarchy Support**: View location trees with or without items.
- **Data Portability**: Export and Import inventory via CSV.

## Homebox v0.26 changes

Homebox v0.26 merged the `/v1/items` and `/v1/locations` APIs into a single
[`/v1/entities`](https://github.com/sysadminsmedia/homebox/pull/1414) resource.
This CLI was updated to match:

| Old endpoint (pre-v0.26)         | v0.26 endpoint                                   |
| -------------------------------- | ------------------------------------------------ |
| `GET/POST /v1/items`             | `GET/POST /v1/entities`                          |
| `GET/PUT/DELETE /v1/items/{id}`  | `GET/PUT/PATCH/DELETE /v1/entities/{id}`         |
| `GET /v1/items/{id}/path`        | `GET /v1/entities/{id}/path`                     |
| `POST /v1/items/{id}/duplicate`  | `POST /v1/entities/{id}/duplicate`               |
| `GET/POST /v1/items/{id}/maintenance` | `GET/POST /v1/entities/{id}/maintenance`     |
| `GET/POST /v1/items/export|import` | `GET/POST /v1/entities/export|import`          |
| `GET/POST /v1/items/{id}/attachments` | `POST /v1/entities/{id}/attachments` (now requires the `name` form field) |
| `GET /v1/locations`              | `GET /v1/entities?isLocation=true`               |
| `GET /v1/locations/tree`         | `GET /v1/entities/tree`                          |
| `POST /v1/locations`             | `POST /v1/entities` (with a location entity type)|
| `GET /v1/currency`               | `GET /v1/currencies`                             |
| `PUT /v1/users/change-password`  | `PUT /v1/users/self/change-password`             |
| `GET/POST/PUT/DELETE /v1/labelmakers` | removed — use `GET /v1/labelmaker/{entity|item|location|asset}/{id}` |
| item custom fields CRUD          | managed through `PUT /v1/entities/{id}` (fields array) |

New in v0.26 and supported by this CLI:

- **API keys** — login is not required for automation: store a key with
  `hb auth token hb_...` or `HB_TOKEN`, and manage keys with `hb auth api-keys`.
- **Entity types** — `hb entity-types list|create|update|delete`.
- **User settings** — `hb auth settings`.

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

### 2. Authenticate

Either log in with your credentials:

```bash
hb login --email your-email@example.com
```

…or use a v0.26 API key (recommended for automation):

```bash
hb auth token hb_your_api_key_here
# or
export HB_TOKEN=hb_your_api_key_here
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

# Generate a label PNG for an item
hb labelmaker get <item-id> -o label.png

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
- `HB_TOKEN`: Authentication token (session token or `hb_` API key)
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
