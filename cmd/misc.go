package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/misonikomipan/homebox-cli/internal/client"
	"github.com/misonikomipan/homebox-cli/internal/config"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Get API status",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(false)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/status", nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
}

func newConfigCmd() *cobra.Command {
	var endpoint string
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Show or update CLI configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if endpoint != "" {
				if err := config.SetEndpoint(endpoint); err != nil {
					return err
				}
				out, _ := json.MarshalIndent(map[string]string{"message": "Endpoint set to " + endpoint}, "", "  ")
				fmt.Println(string(out))
				return nil
			}
			out, _ := json.MarshalIndent(map[string]any{
				"endpoint":      config.GetEndpoint(),
				"authenticated": config.GetToken() != "",
			}, "", "  ")
			fmt.Println(string(out))
			return nil
		},
	}
	cmd.Flags().StringVar(&endpoint, "endpoint", "", "Set API endpoint URL")
	return cmd
}

func newCurrencyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "currency",
		Short: "Get available currency information",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			// v0.26 renamed this endpoint to /v1/currencies
			data, err := c.Get("/v1/currencies", nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
}

func newBarcodeSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "barcode-search <barcode>",
		Short: "Search for a product by barcode/EAN",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			q := url.Values{"s": {args[0]}}
			data, err := c.Get("/v1/products/search-from-barcode", q)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
}

func newGuideCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "guide",
		Short: "Show quick-start guide with usage examples",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(`
┌─────────────────────────────────────┐
│   hb - Homebox CLI Quick Reference  │
└─────────────────────────────────────┘

SETUP
  Set endpoint        hb config --endpoint https://homebox.example.com
  Login               hb login --email you@example.com
  Use API key         hb auth token hb_xxxx (v0.26 API keys)
  Check config        hb config
  API status          hb status

ITEMS (v0.26: entities)
  List all            hb items list
  Search              hb items list --query laptop --page-size 20
  Filter location     hb items list --location <parent-id>
  Filter tag          hb items list --label <tag-id>
  Get                 hb items get <id>
  Create              hb items create --name "MacBook Pro" --location <id>
  Create (full)       hb items create --name "Camera" --quantity 1 --purchase-price 80000
  Update              hb items update <id> --name "New Name"
  Delete              hb items delete <id> --yes
  Duplicate           hb items duplicate <id> --copy-fields
  Path                hb items path <id>
  Export CSV          hb items export --output items.csv
  Import CSV          hb items import items.csv
  By asset ID         hb items asset <asset-id>
  Upload attachment   hb items attachments upload <item-id> photo.jpg
  Delete attachment   hb items attachments delete <item-id> <attachment-id>
  Add field           hb items fields add <item-id> --label "SN" --value "123"

LOCATIONS (v0.26: entities with isLocation=true)
  List                hb locations list
  Tree                hb locations tree
  Tree with items     hb locations tree --with-items
  Get                 hb locations get <id>
  Create              hb locations create --name "書斎"
  Create nested       hb locations create --name "棚A" --parent <parent-id>
  Update              hb locations update <id> --name "新名前"
  Delete              hb locations delete <id> --yes

ENTITY TYPES
  List                hb entity-types list
  Create              hb entity-types create --name "Electronics" --icon "🔌"

TAGS
  List                hb tags list
  Create              hb tags create --name "Electronics" --color "#3b82f6"
  Update              hb tags update <id> --name "Gadgets"
  Delete              hb tags delete <id> --yes

GROUPS
  List                hb groups list
  Info                hb groups info
  Statistics          hb groups stats
  Members             hb groups members
  Invite              hb groups invite --uses 3 --expiry-days 7

MAINTENANCE
  List                hb maintenance list            (status: both|scheduled|completed)
  Create              hb maintenance create --item <id> --name "Oil change" --cost 3000 --completed-date 2026-03-12
  Update              hb maintenance update <id> --cost 4000
  Delete              hb maintenance delete <id> --yes

NOTIFIERS
  List                hb notifiers list
  Create              hb notifiers create --name "Slack" --url https://hooks.slack.com/...
  Test                hb notifiers test
  Delete              hb notifiers delete <id> --yes

TEMPLATES
  List                hb templates list
  Create item         hb templates create-item <template-id> --name "Widget" --location <loc-id>

LABELMAKER
  Label PNG           hb labelmaker get <id> -o label.png (--type entity|item|location|asset)

AUTH
  API keys            hb auth api-keys list|create|delete
  Refresh token       hb auth refresh
  Current user        hb auth me
  Settings            hb auth settings
  Update profile      hb auth update-me --name "New Name"
  Change password     hb auth change-password
  Logout              hb logout

ENV VARS
  Override endpoint   HB_ENDPOINT=https://homebox.example.com hb items list
  Override token      HB_TOKEN=<token> hb items list

TIPS
  Pretty-print        hb items list | jq '.items[]?.name'
  Get IDs             hb locations list | jq '.items[].id'
  Sub-command help    hb items create --help
`)
		},
	}
}
