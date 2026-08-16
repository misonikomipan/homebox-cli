package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/misonikomipan/homebox-cli/internal/client"
	"github.com/misonikomipan/homebox-cli/internal/config"
	"github.com/spf13/cobra"
)

// newEntityTypesCmd manages the v0.26 entity type resource. Every entity
// (item or location) belongs to an entity type, which carries the isLocation
// flag.
func newEntityTypesCmd() *cobra.Command {
	et := &cobra.Command{
		Use:   "entity-types",
		Short: "Manage entity types",
	}

	et.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all entity types",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/entity-types", nil)
			if err != nil {
				return err
			}

			if config.GetFormat() == "table" {
				var types []struct {
					ID         string `json:"id"`
					Name       string `json:"name"`
					IsLocation bool   `json:"isLocation"`
					Icon       string `json:"icon"`
				}
				if err := json.Unmarshal(data, &types); err == nil {
					headers := []string{"ID", "Name", "IsLocation", "Icon"}
					rows := make([][]any, len(types))
					for i, t := range types {
						rows[i] = []any{t.ID, t.Name, t.IsLocation, t.Icon}
					}
					client.Print(data, headers, rows)
					return nil
				}
			}

			client.Print(data, nil, nil)
			return nil
		},
	})

	var createName, createIcon, createTemplateID string
	var createLocation bool
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new entity type",
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
				"name":       createName,
				"isLocation": createLocation,
				"icon":       createIcon,
			}
			if createTemplateID != "" {
				payload["defaultTemplateId"] = createTemplateID
			}
			data, err := c.Post("/v1/entity-types", payload)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "Type name")
	createCmd.Flags().BoolVar(&createLocation, "location", false, "Mark as a location type")
	createCmd.Flags().StringVarP(&createIcon, "icon", "i", "", "Icon (emoji or URL)")
	createCmd.Flags().StringVar(&createTemplateID, "default-template", "", "Default template ID")
	et.AddCommand(createCmd)

	var updateName, updateIcon, updateTemplateID string
	var updateLocation bool
	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an entity type",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/entity-types", nil)
			if err != nil {
				return err
			}
			var types []map[string]any
			if err := json.Unmarshal(data, &types); err != nil {
				return err
			}
			var cur map[string]any
			for _, t := range types {
				if stringField(t, "id") == args[0] {
					cur = t
					break
				}
			}
			if cur == nil {
				return fmt.Errorf("entity type %s not found", args[0])
			}
			payload := map[string]any{
				"name":       stringField(cur, "name"),
				"isLocation": cur["isLocation"],
				"icon":       stringField(cur, "icon"),
			}
			if tid, ok := cur["defaultTemplateId"]; ok && tid != nil {
				payload["defaultTemplateId"] = tid
			}
			if cmd.Flags().Changed("name") {
				payload["name"] = updateName
			}
			if cmd.Flags().Changed("location") {
				payload["isLocation"] = updateLocation
			}
			if cmd.Flags().Changed("icon") {
				payload["icon"] = updateIcon
			}
			if cmd.Flags().Changed("default-template") {
				payload["defaultTemplateId"] = updateTemplateID
			}
			out, err := c.Put("/v1/entity-types/"+args[0], payload)
			if err != nil {
				return err
			}
			client.PrintJSON(out)
			return nil
		},
	}
	updateCmd.Flags().StringVarP(&updateName, "name", "n", "", "Type name")
	updateCmd.Flags().BoolVar(&updateLocation, "location", false, "Mark as a location type")
	updateCmd.Flags().StringVarP(&updateIcon, "icon", "i", "", "Icon")
	updateCmd.Flags().StringVar(&updateTemplateID, "default-template", "", "Default template ID")
	et.AddCommand(updateCmd)

	var deleteYes bool
	deleteCmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an entity type",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !deleteYes {
				if !confirm("Delete entity type " + args[0] + "?") {
					return nil
				}
			}
			c, err := client.New(true)
			if err != nil {
				return err
			}
			if _, err := c.Delete("/v1/entity-types/" + args[0]); err != nil {
				return err
			}
			fmt.Printf(`{"message": "Entity type %s deleted"}`+"\n", args[0])
			return nil
		},
	}
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip confirmation")
	et.AddCommand(deleteCmd)

	return et
}
