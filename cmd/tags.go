package cmd

import (
	"fmt"

	"github.com/misonikomipan/homebox-cli/internal/client"
	"github.com/spf13/cobra"
)

func newTagsCmd() *cobra.Command {
	tags := &cobra.Command{
		Use:   "tags",
		Short: "Manage tags/labels",
	}

	tags.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/tags", nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	tags.AddCommand(&cobra.Command{
		Use:   "get <id>",
		Short: "Get tag details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/tags/"+args[0], nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	var createName, createColor, createDesc string
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new tag",
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
				"name":        createName,
				"description": createDesc,
			}
			if createColor != "" {
				payload["color"] = createColor
			}
			data, err := c.Post("/v1/tags", payload)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "Tag name")
	createCmd.Flags().StringVarP(&createColor, "color", "c", "", "Hex color (e.g. #ff0000)")
	createCmd.Flags().StringVarP(&createDesc, "description", "d", "", "Description")
	tags.AddCommand(createCmd)

	var updateName, updateColor, updateDesc string
	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a tag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/tags/"+args[0], nil)
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
			if cmd.Flags().Changed("color") {
				payload["color"] = updateColor
			}
			if cmd.Flags().Changed("description") {
				payload["description"] = updateDesc
			}
			out, err := c.Put("/v1/tags/"+args[0], payload)
			if err != nil {
				return err
			}
			client.PrintJSON(out)
			return nil
		},
	}
	updateCmd.Flags().StringVarP(&updateName, "name", "n", "", "Tag name")
	updateCmd.Flags().StringVarP(&updateColor, "color", "c", "", "Hex color")
	updateCmd.Flags().StringVarP(&updateDesc, "description", "d", "", "Description")
	tags.AddCommand(updateCmd)

	var deleteYes bool
	deleteCmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a tag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !deleteYes {
				if !confirm("Delete tag " + args[0] + "?") {
					return nil
				}
			}
			c, err := client.New(true)
			if err != nil {
				return err
			}
			if _, err := c.Delete("/v1/tags/" + args[0]); err != nil {
				return err
			}
			fmt.Printf(`{"message": "Tag %s deleted"}`+"\n", args[0])
			return nil
		},
	}
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip confirmation")
	tags.AddCommand(deleteCmd)

	return tags
}
