package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/misonikomipan/homebox-cli/internal/client"
	"github.com/misonikomipan/homebox-cli/internal/config"
	"github.com/spf13/cobra"
)

func newLocationsCmd() *cobra.Command {
	loc := &cobra.Command{
		Use:   "locations",
		Short: "Manage locations",
	}

	loc.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all locations",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/locations", nil)
			if err != nil {
				return err
			}

			if config.GetFormat() == "table" {
				var locations []struct {
					ID          string `json:"id"`
					Name        string `json:"name"`
					Description string `json:"description"`
				}
				if err := json.Unmarshal(data, &locations); err == nil {
					headers := []string{"ID", "Name", "Description"}
					rows := make([][]any, len(locations))
					for i, l := range locations {
						rows[i] = []any{l.ID, l.Name, l.Description}
					}
					client.Print(data, headers, rows)
					return nil
				}
			}

			client.Print(data, nil, nil)
			return nil
		},
	})

	var withItems bool
	treeCmd := &cobra.Command{
		Use:   "tree",
		Short: "Get location tree structure",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			q := url.Values{}
			if withItems {
				q.Set("withItems", "true")
			}
			data, err := c.Get("/v1/locations/tree", q)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	treeCmd.Flags().BoolVar(&withItems, "with-items", false, "Include items in response")
	loc.AddCommand(treeCmd)

	loc.AddCommand(&cobra.Command{
		Use:   "get <id>",
		Short: "Get location details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/locations/"+args[0], nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	var createName, createDesc, createParent string
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new location",
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
			if createParent != "" {
				payload["parentId"] = createParent
			}
			data, err := c.Post("/v1/locations", payload)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "Location name")
	createCmd.Flags().StringVarP(&createDesc, "description", "d", "", "Description")
	createCmd.Flags().StringVarP(&createParent, "parent", "p", "", "Parent location ID")
	loc.AddCommand(createCmd)

	var updateName, updateDesc, updateParent string
	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a location",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/locations/"+args[0], nil)
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
			if cmd.Flags().Changed("parent") {
				payload["parentId"] = updateParent
			}
			out, err := c.Put("/v1/locations/"+args[0], payload)
			if err != nil {
				return err
			}
			client.PrintJSON(out)
			return nil
		},
	}
	updateCmd.Flags().StringVarP(&updateName, "name", "n", "", "Location name")
	updateCmd.Flags().StringVarP(&updateDesc, "description", "d", "", "Description")
	updateCmd.Flags().StringVarP(&updateParent, "parent", "p", "", "Parent location ID")
	loc.AddCommand(updateCmd)

	var deleteYes bool
	deleteCmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a location",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !deleteYes {
				if !confirm("Delete location " + args[0] + "?") {
					return nil
				}
			}
			c, err := client.New(true)
			if err != nil {
				return err
			}
			if _, err := c.Delete("/v1/locations/" + args[0]); err != nil {
				return err
			}
			fmt.Printf(`{"message": "Location %s deleted"}`+"\n", args[0])
			return nil
		},
	}
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip confirmation")
	loc.AddCommand(deleteCmd)

	return loc
}
