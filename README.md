# hb — Homebox CLI

[日本語版はこちら](README.ja.md)

A command-line interface for the [Homebox](https://homebox.software) REST API, written in Go.

## Installation

### Build from source

```bash
git clone https://github.com/misonikomipan/homebox-cli.git
cd homebox-cli
go build -o hb .
mv hb /usr/local/bin/hb
```

## Quick Start

```bash
# 1. (Optional) Set a custom endpoint — default is https://homebox.mizobuchi.dev
hb config --endpoint https://your-homebox.example.com

# 2. Login
hb login --email you@example.com

# 3. Verify connection
hb status

# 4. List items
hb items list
```

## Configuration

Configuration is stored at `~/.config/hb/config.json` (permissions `0600`).

| Key        | Description            |
|------------|------------------------|
| `endpoint` | Homebox server URL     |
| `token`    | Bearer auth token      |

### Environment variables

Environment variables take precedence over the config file.

| Variable      | Description                  |
|---------------|------------------------------|
| `HB_ENDPOINT` | Override API endpoint URL    |
| `HB_TOKEN`    | Override authentication token|

```bash
HB_ENDPOINT=https://homebox.example.com hb items list
```

## Commands

```
hb [command]

Top-level:
  login           Login and store authentication token
  logout          Logout and clear stored token
  status          Get API status
  config          Show or update CLI configuration
  guide           Show quick-start guide with usage examples
  currency        Get available currency information
  barcode-search  Search for a product by barcode/EAN
```

### auth

```bash
hb auth me                          # Current user info
hb auth refresh                     # Refresh token
hb auth update-me --name "New Name" # Update profile
hb auth change-password             # Change password
```

### items

```bash
hb items list                                      # List all items
hb items list --query "laptop" --page-size 20      # Search
hb items list --location <id>                      # Filter by location
hb items list --label <id>                         # Filter by tag
hb items get <id>
hb items create --name "MacBook Pro" --location <id>
hb items create --name "Camera" --quantity 1 --purchase-price 80000 --notes "Sony A7"
hb items update <id> --name "New Name"
hb items delete <id> --yes
hb items duplicate <id>
hb items path <id>                                 # Hierarchy path
hb items maintenance <id>                          # Maintenance logs
hb items export --output items.csv
hb items import items.csv
hb items asset <asset-id>                          # Lookup by asset ID
hb items attachments upload <item-id> photo.jpg
hb items attachments delete <item-id> <attachment-id>
```

### locations

```bash
hb locations list
hb locations tree
hb locations tree --with-items
hb locations get <id>
hb locations create --name "書斎"
hb locations create --name "棚A" --parent <parent-id>
hb locations update <id> --name "新名前"
hb locations delete <id> --yes
```

### tags

```bash
hb tags list
hb tags get <id>
hb tags create --name "Electronics" --color "#3b82f6"
hb tags update <id> --name "Gadgets"
hb tags delete <id> --yes
```

### groups

```bash
hb groups info
hb groups stats
hb groups members
hb groups update --name "Home" --currency "JPY"
hb groups invite --uses 3 --expiry-days 7
```

### maintenance

```bash
hb maintenance list
hb maintenance create --item <id> --name "Oil change" --cost 3000
hb maintenance update <id> --completed-date 2026-03-12
hb maintenance delete <id> --yes
```

### notifiers

```bash
hb notifiers list
hb notifiers create --name "Slack" --url https://hooks.slack.com/...
hb notifiers update <id> --active=false
hb notifiers test
hb notifiers delete <id> --yes
```

### templates

```bash
hb templates list
hb templates get <id>
hb templates create --name "PC Template"
hb templates update <id> --name "Laptop Template"
hb templates create-item <template-id> --location <id>
hb templates delete <id> --yes
```

## Tips

```bash
# Pretty-print with jq
hb items list | jq '.items[]?.name'

# Get all location IDs
hb locations list | jq '.[].id'

# Show all commands for a subgroup
hb items --help
hb items create --help
```
