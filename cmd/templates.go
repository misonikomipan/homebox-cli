package cmd

import (
	"fmt"

	"github.com/misonikomipan/homebox-cli/internal/client"
	"github.com/spf13/cobra"
)

func newTemplatesCmd() *cobra.Command {
	t := &cobra.Command{
		Use:   "templates",
		Short: "Manage item templates",
	}

	t.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/templates", nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	t.AddCommand(&cobra.Command{
		Use:   "get <id>",
		Short: "Get template details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/templates/"+args[0], nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	var createName, createDesc string
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new item template",
		RunE: func(cmd *cobra.Command, args []string) error {
			if createName == "" {
				fmt.Print("Name: ")
				fmt.Scanln(&createName)
			}
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Post("/v1/templates", map[string]any{
				"name":        createName,
				"description": createDesc,
			})
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "Template name")
	createCmd.Flags().StringVarP(&createDesc, "description", "d", "", "Description")
	t.AddCommand(createCmd)

	var updateName, updateDesc string
	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/templates/"+args[0], nil)
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
			if cmd.Flags().Changed("description") {
				payload["description"] = updateDesc
			}
			out, err := c.Put("/v1/templates/"+args[0], payload)
			if err != nil {
				return err
			}
			client.PrintJSON(out)
			return nil
		},
	}
	updateCmd.Flags().StringVarP(&updateName, "name", "n", "", "Template name")
	updateCmd.Flags().StringVarP(&updateDesc, "description", "d", "", "Description")
	t.AddCommand(updateCmd)

	var deleteYes bool
	deleteCmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !deleteYes {
				if !confirm("Delete template " + args[0] + "?") {
					return nil
				}
			}
			c, err := client.New(true)
			if err != nil {
				return err
			}
			if _, err := c.Delete("/v1/templates/" + args[0]); err != nil {
				return err
			}
			fmt.Printf(`{"message": "Template %s deleted"}`+"\n", args[0])
			return nil
		},
	}
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip confirmation")
	t.AddCommand(deleteCmd)

	var createItemLocID string
	createItemCmd := &cobra.Command{
		Use:   "create-item <template-id>",
		Short: "Create an item from a template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			payload := map[string]any{}
			if createItemLocID != "" {
				payload["locationId"] = createItemLocID
			}
			data, err := c.Post("/v1/templates/"+args[0]+"/create-item", payload)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	createItemCmd.Flags().StringVarP(&createItemLocID, "location", "l", "", "Location ID")
	t.AddCommand(createItemCmd)

	return t
}
