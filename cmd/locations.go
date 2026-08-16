package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/misonikomipan/homebox-cli/internal/client"
	"github.com/misonikomipan/homebox-cli/internal/config"
	"github.com/spf13/cobra"
)

func newLocationsCmd() *cobra.Command {
	loc := &cobra.Command{
		Use:   "locations",
		Short: "Manage locations (entity type: location)",
	}

	// list
	var page, pageSize int
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all locations",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			// v0.26: locations are entities with isLocation=true
			q := url.Values{
				"isLocation": {"true"},
				"page":       {strconv.Itoa(page)},
				"pageSize":   {strconv.Itoa(pageSize)},
			}
			data, err := c.Get(entityBasePath, q)
			if err != nil {
				return err
			}

			if config.GetFormat() == "table" {
				var resp struct {
					Items []struct {
						ID          string  `json:"id"`
						Name        string  `json:"name"`
						Description string  `json:"description"`
						ItemCount   float64 `json:"itemCount"`
					} `json:"items"`
				}
				if err := json.Unmarshal(data, &resp); err == nil {
					headers := []string{"ID", "Name", "Description", "Items"}
					rows := make([][]any, len(resp.Items))
					for i, l := range resp.Items {
						rows[i] = []any{l.ID, l.Name, l.Description, l.ItemCount}
					}
					client.Print(data, headers, rows)
					return nil
				}
			}

			client.Print(data, nil, nil)
			return nil
		},
	}
	listCmd.Flags().IntVar(&page, "page", 1, "Page number")
	listCmd.Flags().IntVar(&pageSize, "page-size", 10, "Items per page")
	loc.AddCommand(listCmd)

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
			data, err := c.Get(entityBasePath+"/tree", q)
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
			data, err := c.Get(entityBasePath+"/"+args[0], nil)
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
			etID, err := locationEntityTypeID(c)
			if err != nil {
				return err
			}
			payload := map[string]any{
				"name":         createName,
				"description":  createDesc,
				"entityTypeId": etID,
			}
			if createParent != "" {
				payload["parentId"] = createParent
			}
			data, err := c.Post(entityBasePath, payload)
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
			cur, err := fetchEntity(c, args[0])
			if err != nil {
				return err
			}
			payload := updatePayload(cur)
			if cmd.Flags().Changed("name") {
				payload["name"] = updateName
			}
			if cmd.Flags().Changed("description") {
				payload["description"] = updateDesc
			}
			if cmd.Flags().Changed("parent") {
				payload["parentId"] = updateParent
			}
			out, err := putEntity(c, args[0], payload)
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
			if _, err := c.Delete(entityBasePath + "/" + args[0]); err != nil {
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
