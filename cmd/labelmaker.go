package cmd

import (
	"fmt"

	"github.com/misonikomipan/homebox-cli/internal/client"
	"github.com/spf13/cobra"
)

func newLabelmakerCmd() *cobra.Command {
	lm := &cobra.Command{
		Use:   "labelmaker",
		Short: "Manage labelmaker configurations",
	}

	lm.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all labelmaker configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/labelmakers", nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	lm.AddCommand(&cobra.Command{
		Use:   "get <id>",
		Short: "Get labelmaker configuration details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/labelmakers/"+args[0], nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	var createName, createConfig string
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new labelmaker configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			if createName == "" {
				fmt.Print("Name: ")
				fmt.Scanln(&createName)
			}
			c, err := client.New(true)
			if err != nil {
				return err
			}
			payload := map[string]any{
				"name":   createName,
				"config": createConfig,
			}
			data, err := c.Post("/v1/labelmakers", payload)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "Configuration name")
	createCmd.Flags().StringVarP(&createConfig, "config", "c", "", "Configuration JSON")
	lm.AddCommand(createCmd)

	var updateName, updateConfig string
	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a labelmaker configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/labelmakers/"+args[0], nil)
			if err != nil {
				return err
			}
			var payload map[string]any
			if err := unmarshalJSON(data, &payload); err != nil {
				return err
			}
			if cmd.Flags().Changed("name") {
				payload["name"] = updateName
			}
			if cmd.Flags().Changed("config") {
				payload["config"] = updateConfig
			}
			out, err := c.Put("/v1/labelmakers/"+args[0], payload)
			if err != nil {
				return err
			}
			client.PrintJSON(out)
			return nil
		},
	}
	updateCmd.Flags().StringVarP(&updateName, "name", "n", "", "Configuration name")
	updateCmd.Flags().StringVarP(&updateConfig, "config", "c", "", "Configuration JSON")
	lm.AddCommand(updateCmd)

	var deleteYes bool
	deleteCmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a labelmaker configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !deleteYes {
				if !confirm("Delete labelmaker configuration " + args[0] + "?") {
					return nil
				}
			}
			c, err := client.New(true)
			if err != nil {
				return err
			}
			if _, err := c.Delete("/v1/labelmakers/" + args[0]); err != nil {
				return err
			}
			fmt.Printf("{\"message\": \"Labelmaker configuration %s deleted\"}\n", args[0])
			return nil
		},
	}
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip confirmation")
	lm.AddCommand(deleteCmd)

	return lm
}
