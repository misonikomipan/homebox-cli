package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/misonikomipan/homebox-cli/internal/client"
	"github.com/misonikomipan/homebox-cli/internal/config"
	"github.com/spf13/cobra"
)

func newItemsCmd() *cobra.Command {
	items := &cobra.Command{
		Use:   "items",
		Short: "Manage inventory items",
	}

	// list
	var query, locationID, labelID string
	var page, pageSize int
	var archived bool
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all items",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			q := url.Values{
				"page":     {strconv.Itoa(page)},
				"pageSize": {strconv.Itoa(pageSize)},
			}
			if query != "" {
				q.Set("q", query)
			}
			if locationID != "" {
				q.Set("locations", locationID)
			}
			if labelID != "" {
				q.Set("labels", labelID)
			}
			if archived {
				q.Set("includeArchived", "true")
			}
			data, err := c.Get("/v1/items", q)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	listCmd.Flags().StringVarP(&query, "query", "q", "", "Search query")
	listCmd.Flags().StringVarP(&locationID, "location", "l", "", "Filter by location ID")
	listCmd.Flags().StringVar(&labelID, "label", "", "Filter by label ID")
	listCmd.Flags().IntVar(&page, "page", 1, "Page number")
	listCmd.Flags().IntVar(&pageSize, "page-size", 10, "Items per page")
	listCmd.Flags().BoolVar(&archived, "archived", false, "Include archived items")
	items.AddCommand(listCmd)

	// get
	items.AddCommand(&cobra.Command{
		Use:   "get <id>",
		Short: "Get item details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/items/"+args[0], nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	// create
	var createName, createDesc, createLocID, createNotes, createAssetID, createPurchaseFrom string
	var createLabels []string
	var createQty int
	var createPrice float64
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new item",
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
				"quantity":    createQty,
				"notes":       createNotes,
			}
			if createLocID != "" {
				payload["locationId"] = createLocID
			}
			if len(createLabels) > 0 {
				payload["labelIds"] = createLabels
			}
			if createAssetID != "" {
				payload["assetId"] = createAssetID
			}
			if cmd.Flags().Changed("purchase-price") {
				payload["purchasePrice"] = createPrice
			}
			if createPurchaseFrom != "" {
				payload["purchaseFrom"] = createPurchaseFrom
			}
			data, err := c.Post("/v1/items", payload)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "Item name")
	createCmd.Flags().StringVarP(&createDesc, "description", "d", "", "Description")
	createCmd.Flags().StringVarP(&createLocID, "location", "l", "", "Location ID")
	createCmd.Flags().StringArrayVar(&createLabels, "label", nil, "Label IDs (repeatable)")
	createCmd.Flags().IntVarP(&createQty, "quantity", "q", 1, "Quantity")
	createCmd.Flags().StringVar(&createAssetID, "asset-id", "", "Asset ID")
	createCmd.Flags().Float64Var(&createPrice, "purchase-price", 0, "Purchase price")
	createCmd.Flags().StringVar(&createPurchaseFrom, "purchase-from", "", "Purchased from")
	createCmd.Flags().StringVar(&createNotes, "notes", "", "Notes")
	items.AddCommand(createCmd)

	// update
	var updateName, updateDesc, updateLocID, updateNotes, updatePurchaseFrom string
	var updateLabels []string
	var updateQty int
	var updatePrice, updateSoldPrice, updateInsuredFor float64
	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			// fetch current
			data, err := c.Get("/v1/items/"+args[0], nil)
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
			if cmd.Flags().Changed("location") {
				payload["locationId"] = updateLocID
			}
			if cmd.Flags().Changed("label") {
				payload["labelIds"] = updateLabels
			}
			if cmd.Flags().Changed("quantity") {
				payload["quantity"] = updateQty
			}
			if cmd.Flags().Changed("notes") {
				payload["notes"] = updateNotes
			}
			if cmd.Flags().Changed("purchase-price") {
				payload["purchasePrice"] = updatePrice
			}
			if cmd.Flags().Changed("purchase-from") {
				payload["purchaseFrom"] = updatePurchaseFrom
			}
			if cmd.Flags().Changed("sold-price") {
				payload["soldPrice"] = updateSoldPrice
			}
			if cmd.Flags().Changed("insured-for") {
				payload["insuredFor"] = updateInsuredFor
			}
			out, err := c.Put("/v1/items/"+args[0], payload)
			if err != nil {
				return err
			}
			client.PrintJSON(out)
			return nil
		},
	}
	updateCmd.Flags().StringVarP(&updateName, "name", "n", "", "Item name")
	updateCmd.Flags().StringVarP(&updateDesc, "description", "d", "", "Description")
	updateCmd.Flags().StringVarP(&updateLocID, "location", "l", "", "Location ID")
	updateCmd.Flags().StringArrayVar(&updateLabels, "label", nil, "Label IDs")
	updateCmd.Flags().IntVarP(&updateQty, "quantity", "q", 1, "Quantity")
	updateCmd.Flags().StringVar(&updateNotes, "notes", "", "Notes")
	updateCmd.Flags().Float64Var(&updatePrice, "purchase-price", 0, "Purchase price")
	updateCmd.Flags().StringVar(&updatePurchaseFrom, "purchase-from", "", "Purchased from")
	updateCmd.Flags().Float64Var(&updateSoldPrice, "sold-price", 0, "Sold price")
	updateCmd.Flags().Float64Var(&updateInsuredFor, "insured-for", 0, "Insured for")
	items.AddCommand(updateCmd)

	// delete
	var deleteYes bool
	deleteCmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !deleteYes {
				if !confirm("Delete item " + args[0] + "?") {
					return nil
				}
			}
			c, err := client.New(true)
			if err != nil {
				return err
			}
			if _, err := c.Delete("/v1/items/" + args[0]); err != nil {
				return err
			}
			fmt.Printf(`{"message": "Item %s deleted"}`+"\n", args[0])
			return nil
		},
	}
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip confirmation")
	items.AddCommand(deleteCmd)

	// duplicate
	items.AddCommand(&cobra.Command{
		Use:   "duplicate <id>",
		Short: "Duplicate an item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Post("/v1/items/"+args[0]+"/duplicate", nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	// path
	items.AddCommand(&cobra.Command{
		Use:   "path <id>",
		Short: "Get item hierarchy path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/items/"+args[0]+"/path", nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	// maintenance
	items.AddCommand(&cobra.Command{
		Use:   "maintenance <id>",
		Short: "List maintenance logs for an item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/items/"+args[0]+"/maintenance", nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	// export
	var exportOutput string
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export all items as CSV",
		RunE: func(cmd *cobra.Command, args []string) error {
			token := config.GetToken()
			if token == "" {
				return fmt.Errorf("not authenticated — run 'hb login' first")
			}
			endpoint := config.GetEndpoint()
			req, _ := http.NewRequest("GET", endpoint+"/api/v1/items/export", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp, err := (&http.Client{Timeout: 60 * time.Second}).Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			if exportOutput != "" {
				if err := os.WriteFile(exportOutput, body, 0644); err != nil {
					return err
				}
				fmt.Printf(`{"message": "Exported to %s"}`+"\n", exportOutput)
			} else {
				fmt.Print(string(body))
			}
			return nil
		},
	}
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file path")
	items.AddCommand(exportCmd)

	// import
	items.AddCommand(&cobra.Command{
		Use:   "import <csv-file>",
		Short: "Import items from CSV file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			token := config.GetToken()
			if token == "" {
				return fmt.Errorf("not authenticated — run 'hb login' first")
			}
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()

			var buf bytes.Buffer
			w := multipart.NewWriter(&buf)
			part, _ := w.CreateFormFile("csv", filepath.Base(args[0]))
			io.Copy(part, f)
			w.Close()

			endpoint := config.GetEndpoint()
			req, _ := http.NewRequest("POST", endpoint+"/api/v1/items/import", &buf)
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", w.FormDataContentType())
			resp, err := (&http.Client{Timeout: 60 * time.Second}).Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			client.PrintJSON(body)
			return nil
		},
	})

	// asset
	items.AddCommand(&cobra.Command{
		Use:   "asset <asset-id>",
		Short: "Get item by asset ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Get("/v1/assets/"+args[0], nil)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	})

	// attachments sub-group
	attachments := &cobra.Command{
		Use:   "attachments",
		Short: "Manage item attachments",
	}

	var attachType, attachName string
	uploadCmd := &cobra.Command{
		Use:   "upload <item-id> <file>",
		Short: "Upload an attachment to an item",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			token := config.GetToken()
			if token == "" {
				return fmt.Errorf("not authenticated — run 'hb login' first")
			}
			f, err := os.Open(args[1])
			if err != nil {
				return err
			}
			defer f.Close()

			name := attachName
			if name == "" {
				name = filepath.Base(args[1])
			}
			var buf bytes.Buffer
			w := multipart.NewWriter(&buf)
			w.WriteField("type", attachType)
			part, _ := w.CreateFormFile("file", name)
			io.Copy(part, f)
			w.Close()

			endpoint := config.GetEndpoint()
			req, _ := http.NewRequest("POST", endpoint+"/api/v1/items/"+args[0]+"/attachments", &buf)
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", w.FormDataContentType())
			resp, err := (&http.Client{Timeout: 60 * time.Second}).Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			client.PrintJSON(body)
			return nil
		},
	}
	uploadCmd.Flags().StringVarP(&attachType, "type", "t", "attachment", "Attachment type")
	uploadCmd.Flags().StringVar(&attachName, "name", "", "Override file name")
	attachments.AddCommand(uploadCmd)

	var attachDeleteYes bool
	attachDeleteCmd := &cobra.Command{
		Use:   "delete <item-id> <attachment-id>",
		Short: "Delete an item attachment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !attachDeleteYes {
				if !confirm("Delete attachment " + args[1] + "?") {
					return nil
				}
			}
			c, err := client.New(true)
			if err != nil {
				return err
			}
			if _, err := c.Delete("/v1/items/" + args[0] + "/attachments/" + args[1]); err != nil {
				return err
			}
			fmt.Printf(`{"message": "Attachment %s deleted"}`+"\n", args[1])
			return nil
		},
	}
	attachDeleteCmd.Flags().BoolVarP(&attachDeleteYes, "yes", "y", false, "Skip confirmation")
	attachments.AddCommand(attachDeleteCmd)
	items.AddCommand(attachments)

	return items
}
