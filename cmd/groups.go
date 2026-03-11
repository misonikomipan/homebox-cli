package cmd

import (
	"fmt"
	"time"

	"github.com/misonikomipan/homebox-cli/internal/client"
	"github.com/spf13/cobra"
)

func newGroupsCmd() *cobra.Command {
	groups := &cobra.Command{
		Use:   "groups",
		Short: "Manage groups and members",
	}

	groups.AddCommand(&cobra.Command{
		Use:   "info",
		Short: "Get current group information",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/groups", nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	var updateName, updateCurrency string
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update group settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			payload := map[string]any{}
			if updateName != "" {
				payload["name"] = updateName
			}
			if updateCurrency != "" {
				payload["currency"] = updateCurrency
			}
			data, err := c.Put("/v1/groups", payload)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	updateCmd.Flags().StringVarP(&updateName, "name", "n", "", "Group name")
	updateCmd.Flags().StringVarP(&updateCurrency, "currency", "c", "", "Currency code")
	groups.AddCommand(updateCmd)

	groups.AddCommand(&cobra.Command{
		Use:   "stats",
		Short: "Get group statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/groups/statistics", nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	groups.AddCommand(&cobra.Command{
		Use:   "members",
		Short: "List group members",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/groups/members", nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	var inviteUses, inviteExpiry int
	inviteCmd := &cobra.Command{
		Use:   "invite",
		Short: "Create a group invitation link",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			expiry := time.Now().UTC().Add(time.Duration(inviteExpiry) * 24 * time.Hour).Format(time.RFC3339)
			data, err := c.Post("/v1/groups/invitations", map[string]any{
				"uses":      inviteUses,
				"expiresAt": expiry,
			})
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	inviteCmd.Flags().IntVarP(&inviteUses, "uses", "u", 1, "Number of uses")
	inviteCmd.Flags().IntVar(&inviteExpiry, "expiry-days", 7, "Days until expiry")
	groups.AddCommand(inviteCmd)

	_ = fmt.Sprintf // suppress unused import
	return groups
}
