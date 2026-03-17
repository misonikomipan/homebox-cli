package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/misonikomipan/homebox-cli/internal/client"
	"github.com/misonikomipan/homebox-cli/internal/config"
	"github.com/spf13/cobra"
)

func newNotifiersCmd() *cobra.Command {
	n := &cobra.Command{
		Use:   "notifiers",
		Short: "Manage notifiers/webhooks",
	}

	n.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all notifiers",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/notifiers", nil)
			if err != nil {
				return err
			}

			if config.GetFormat() == "table" {
				var notifiers []struct {
					ID       string `json:"id"`
					Name     string `json:"name"`
					URL      string `json:"url"`
					IsActive bool   `json:"isActive"`
				}
				if err := json.Unmarshal(data, &notifiers); err == nil {
					headers := []string{"ID", "Name", "URL", "Active"}
					rows := make([][]any, len(notifiers))
					for i, n := range notifiers {
						rows[i] = []any{n.ID, n.Name, n.URL, n.IsActive}
					}
					client.Print(data, headers, rows)
					return nil
				}
			}

			client.Print(data, nil, nil)
			return nil
		},
	})

	var createName, createURL string
	var createActive bool
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new notifier",
		RunE: func(cmd *cobra.Command, args []string) error {
			if createName == "" {
				fmt.Print("Name: ")
				fmt.Scanln(&createName)
			}
			if createURL == "" {
				fmt.Print("URL: ")
				fmt.Scanln(&createURL)
			}
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Post("/v1/notifiers", map[string]any{
				"name":     createName,
				"url":      createURL,
				"isActive": createActive,
			})
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "Notifier name")
	createCmd.Flags().StringVarP(&createURL, "url", "u", "", "Webhook URL")
	createCmd.Flags().BoolVar(&createActive, "active", true, "Is active")
	n.AddCommand(createCmd)

	var updateName, updateURL string
	var updateActive bool
	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a notifier",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			payload := map[string]any{}
			if cmd.Flags().Changed("name") {
				payload["name"] = updateName
			}
			if cmd.Flags().Changed("url") {
				payload["url"] = updateURL
			}
			if cmd.Flags().Changed("active") {
				payload["isActive"] = updateActive
			}
			data, err := c.Put("/v1/notifiers/"+args[0], payload)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	updateCmd.Flags().StringVarP(&updateName, "name", "n", "", "Notifier name")
	updateCmd.Flags().StringVarP(&updateURL, "url", "u", "", "Webhook URL")
	updateCmd.Flags().BoolVar(&updateActive, "active", true, "Is active")
	n.AddCommand(updateCmd)

	var deleteYes bool
	deleteCmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a notifier",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !deleteYes {
				if !confirm("Delete notifier " + args[0] + "?") {
					return nil
				}
			}
			c, err := client.New(true)
			if err != nil {
				return err
			}
			if _, err := c.Delete("/v1/notifiers/" + args[0]); err != nil {
				return err
			}
			fmt.Printf(`{"message": "Notifier %s deleted"}`+"\n", args[0])
			return nil
		},
	}
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip confirmation")
	n.AddCommand(deleteCmd)

	n.AddCommand(&cobra.Command{
		Use:   "test",
		Short: "Send a test notification",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Post("/v1/notifiers/test", nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	return n
}
