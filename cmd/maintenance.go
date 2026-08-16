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

// maintenancePayload builds the JSON body accepted by the v0.26 maintenance
// endpoints. Cost must be serialized as a string and at least one of
// completedDate / scheduledDate is required by the server.
func maintenancePayload(name, description, completed, scheduled string, cost float64) map[string]any {
	return map[string]any{
		"name":          name,
		"description":   description,
		"cost":          strconv.FormatFloat(cost, 'f', -1, 64),
		"completedDate": completed,
		"scheduledDate": scheduled,
	}
}

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
			// v0.26 rejects an empty status filter; default to "both".
			status := listStatus
			if status == "" {
				status = "both"
			}
			q.Set("status", status)
			data, err := c.Get("/v1/maintenance", q)
			if err != nil {
				return err
			}

			if config.GetFormat() == "table" {
				var raw []map[string]any
				if err := json.Unmarshal(data, &raw); err == nil {
					headers := []string{"ID", "Name", "Cost", "Item"}
					rows := make([][]any, 0, len(raw))
					for _, e := range raw {
						cost, _ := e["cost"].(string)
						rows = append(rows, []any{
							stringField(e, "id"),
							stringField(e, "name"),
							cost,
							stringField(e, "itemName"),
						})
					}
					client.Print(data, headers, rows)
					return nil
				}
			}

			client.Print(data, nil, nil)
			return nil
		},
	}
	listCmd.Flags().StringVarP(&listStatus, "status", "s", "both", "Filter by status (scheduled, completed, both)")
	m.AddCommand(listCmd)

	var createItemID, createName, createNotes, createScheduled, createCompleted string
	var createCost float64
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a maintenance entry for an item",
		RunE: func(cmd *cobra.Command, args []string) error {
			if createItemID == "" {
				return fmt.Errorf("--item is required")
			}
			if createName == "" {
				fmt.Print("Name: ")
				fmt.Scanln(&createName)
			}
			if createCompleted == "" && createScheduled == "" {
				return fmt.Errorf("at least one of --completed-date or --scheduled-date is required")
			}
			c, err := client.New(true)
			if err != nil {
				return err
			}
			payload := maintenancePayload(createName, createNotes, createCompleted, createScheduled, createCost)
			data, err := c.Post(entityBasePath+"/"+createItemID+"/maintenance", payload)
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
	createCmd.Flags().StringVar(&createScheduled, "scheduled-date", "", "Scheduled date (YYYY-MM-DD)")
	createCmd.Flags().StringVar(&createCompleted, "completed-date", "", "Completed date (YYYY-MM-DD)")
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
			// There is no GET /maintenance/{id}; fetch the full list and pick
			// the entry so the mandatory date fields round-trip.
			data, err := c.Get("/v1/maintenance", url.Values{"status": {"both"}})
			if err != nil {
				return err
			}
			var entries []map[string]any
			if err := json.Unmarshal(data, &entries); err != nil {
				return err
			}
			var cur map[string]any
			for _, e := range entries {
				if stringField(e, "id") == args[0] {
					cur = e
					break
				}
			}
			if cur == nil {
				return fmt.Errorf("maintenance entry %s not found", args[0])
			}

			name := stringField(cur, "name")
			notes := stringField(cur, "description")
			completed := stringField(cur, "completedDate")
			scheduled := stringField(cur, "scheduledDate")
			cost := 0.0
			if cs, ok := cur["cost"].(string); ok {
				cost, _ = strconv.ParseFloat(cs, 64)
			}

			if cmd.Flags().Changed("name") {
				name = updateName
			}
			if cmd.Flags().Changed("notes") {
				notes = updateNotes
			}
			if cmd.Flags().Changed("completed-date") {
				completed = updateCompleted
			}
			if cmd.Flags().Changed("scheduled-date") {
				scheduled = updateScheduled
			}
			if cmd.Flags().Changed("cost") {
				cost = updateCost
			}
			if completed == "" && scheduled == "" {
				return fmt.Errorf("at least one of --completed-date or --scheduled-date is required")
			}
			payload := maintenancePayload(name, notes, completed, scheduled, cost)
			out, err := c.Put("/v1/maintenance/"+args[0], payload)
			if err != nil {
				return err
			}
			client.PrintJSON(out)
			return nil
		},
	}
	updateCmd.Flags().StringVarP(&updateName, "name", "n", "", "Entry name")
	updateCmd.Flags().StringVar(&updateNotes, "notes", "", "Notes")
	updateCmd.Flags().Float64Var(&updateCost, "cost", 0, "Cost")
	updateCmd.Flags().StringVar(&updateScheduled, "scheduled-date", "", "Scheduled date (YYYY-MM-DD)")
	updateCmd.Flags().StringVar(&updateCompleted, "completed-date", "", "Completed date (YYYY-MM-DD)")
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
