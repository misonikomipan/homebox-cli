package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/misonikomipan/homebox-cli/internal/client"
	"github.com/misonikomipan/homebox-cli/internal/config"
	"github.com/spf13/cobra"
)

func newMaintenanceCmd() *cobra.Command {
	m := &cobra.Command{
		Use:   "maintenance",
		Short: "Manage maintenance entries",
	}

	var listStatus string
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all maintenance entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			q := url.Values{}
			if listStatus != "" {
				q.Set("status", listStatus)
			}
			data, err := c.Get("/v1/maintenance", q)
			if err != nil {
				return err
			}

			if config.GetFormat() == "table" {
				var entries []struct {
					ID     string  `json:"id"`
					Name   string  `json:"name"`
					Cost   float64 `json:"cost"`
					Status string  `json:"status"`
				}
				if err := json.Unmarshal(data, &entries); err == nil {
					headers := []string{"ID", "Name", "Cost", "Status"}
					rows := make([][]any, len(entries))
					for i, e := range entries {
						rows[i] = []any{e.ID, e.Name, e.Cost, e.Status}
					}
					client.Print(data, headers, rows)
					return nil
				}
			}

			client.Print(data, nil, nil)
			return nil
		},
	}
	listCmd.Flags().StringVarP(&listStatus, "status", "s", "", "Filter by status")
	m.AddCommand(listCmd)

	var createItemID, createName, createNotes, createScheduled, createCompleted string
	var createCost float64
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a maintenance entry for an item",
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
				"name":  createName,
				"notes": createNotes,
				"cost":  createCost,
			}
			if createScheduled != "" {
				payload["scheduledDate"] = createScheduled
			}
			if createCompleted != "" {
				payload["completedDate"] = createCompleted
			}
			data, err := c.Post("/v1/items/"+createItemID+"/maintenance", payload)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	createCmd.Flags().StringVarP(&createItemID, "item", "i", "", "Item ID")
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "Entry name")
	createCmd.Flags().StringVar(&createNotes, "notes", "", "Notes")
	createCmd.Flags().Float64Var(&createCost, "cost", 0, "Cost")
	createCmd.Flags().StringVar(&createScheduled, "scheduled-date", "", "Scheduled date (ISO 8601)")
	createCmd.Flags().StringVar(&createCompleted, "completed-date", "", "Completed date (ISO 8601)")
	_ = createCmd.MarkFlagRequired("item")
	m.AddCommand(createCmd)

	var updateName, updateNotes, updateScheduled, updateCompleted string
	var updateCost float64
	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a maintenance entry",
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
			if cmd.Flags().Changed("notes") {
				payload["notes"] = updateNotes
			}
			if cmd.Flags().Changed("cost") {
				payload["cost"] = updateCost
			}
			if cmd.Flags().Changed("scheduled-date") {
				payload["scheduledDate"] = updateScheduled
			}
			if cmd.Flags().Changed("completed-date") {
				payload["completedDate"] = updateCompleted
			}
			data, err := c.Put("/v1/maintenance/"+args[0], payload)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	updateCmd.Flags().StringVarP(&updateName, "name", "n", "", "Entry name")
	updateCmd.Flags().StringVar(&updateNotes, "notes", "", "Notes")
	updateCmd.Flags().Float64Var(&updateCost, "cost", 0, "Cost")
	updateCmd.Flags().StringVar(&updateScheduled, "scheduled-date", "", "Scheduled date (ISO 8601)")
	updateCmd.Flags().StringVar(&updateCompleted, "completed-date", "", "Completed date (ISO 8601)")
	m.AddCommand(updateCmd)

	var deleteYes bool
	deleteCmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a maintenance entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !deleteYes {
				if !confirm("Delete maintenance entry " + args[0] + "?") {
					return nil
				}
			}
			c, err := client.New(true)
			if err != nil {
				return err
			}
			if _, err := c.Delete("/v1/maintenance/" + args[0]); err != nil {
				return err
			}
			fmt.Printf(`{"message": "Maintenance entry %s deleted"}`+"\n", args[0])
			return nil
		},
	}
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip confirmation")
	m.AddCommand(deleteCmd)

	return m
}
