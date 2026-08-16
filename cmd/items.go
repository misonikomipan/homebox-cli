package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
	var query string
	var locationIDs, labelIDs []string
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
			// v0.26: location filters became parentIds, label filters became tags
			for _, id := range locationIDs {
				q.Add("parentIds", id)
			}
			for _, id := range labelIDs {
				q.Add("tags", id)
			}
			if archived {
				q.Set("includeArchived", "true")
			}
			data, err := c.Get(entityBasePath, q)
			if err != nil {
				return err
			}

			if config.GetFormat() == "table" {
				var resp struct {
					Items []struct {
						ID       string  `json:"id"`
						AssetID  string  `json:"assetId"`
						Name     string  `json:"name"`
						Quantity float64 `json:"quantity"`
						Parent   *struct {
							Name string `json:"name"`
						} `json:"parent"`
					} `json:"items"`
				}
				if err := json.Unmarshal(data, &resp); err == nil {
					headers := []string{"ID", "Asset", "Name", "Qty", "Location"}
					rows := make([][]any, len(resp.Items))
					for i, it := range resp.Items {
						loc := ""
						if it.Parent != nil {
							loc = it.Parent.Name
						}
						rows[i] = []any{it.ID, it.AssetID, it.Name, it.Quantity, loc}
					}
					client.Print(data, headers, rows)
					return nil
				}
			}

			client.Print(data, nil, nil)
			return nil
		},
	}
	listCmd.Flags().StringVarP(&query, "query", "q", "", "Search query")
	listCmd.Flags().StringArrayVarP(&locationIDs, "location", "l", nil, "Filter by location ID (repeatable)")
	listCmd.Flags().StringArrayVar(&labelIDs, "label", nil, "Filter by label/tag ID (repeatable)")
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
			data, err := c.Get(entityBasePath+"/"+args[0], nil)
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
	var createQty float64
	var createPrice, createSoldPrice float64
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
				"name":     createName,
				"quantity": createQty,
			}
			if createDesc != "" {
				payload["description"] = createDesc
			}
			if createLocID != "" {
				payload["parentId"] = createLocID
			}
			if len(createLabels) > 0 {
				payload["tagIds"] = createLabels
			}
			data, err := c.Post(entityBasePath, payload)
			if err != nil {
				return err
			}

			// v0.26: entity create only accepts name/description/quantity/
			// parentId/entityTypeId/tagIds. Purchase and other advanced fields
			// are set with a follow-up update.
			if cmd.Flags().Changed("purchase-price") || cmd.Flags().Changed("purchase-from") ||
				cmd.Flags().Changed("notes") || cmd.Flags().Changed("sold-price") ||
				cmd.Flags().Changed("asset-id") {
				var created map[string]any
				if err := json.Unmarshal(data, &created); err != nil {
					return err
				}
				id, _ := created["id"].(string)
				up := updatePayload(created)
				if cmd.Flags().Changed("purchase-price") {
					up["purchasePrice"] = createPrice
				}
				if cmd.Flags().Changed("purchase-from") {
					up["purchaseFrom"] = createPurchaseFrom
				}
				if cmd.Flags().Changed("notes") {
					up["notes"] = createNotes
				}
				if cmd.Flags().Changed("sold-price") {
					up["soldPrice"] = createSoldPrice
				}
				if cmd.Flags().Changed("asset-id") {
					up["assetId"] = createAssetID
				}
				data, err = putEntity(c, id, up)
				if err != nil {
					return err
				}
			}
			client.PrintJSON(data)
			return nil
		},
	}
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "Item name")
	createCmd.Flags().StringVarP(&createDesc, "description", "d", "", "Description")
	createCmd.Flags().StringVarP(&createLocID, "location", "l", "", "Location (parent) ID")
	createCmd.Flags().StringArrayVar(&createLabels, "label", nil, "Label/tag IDs (repeatable)")
	createCmd.Flags().Float64VarP(&createQty, "quantity", "q", 1, "Quantity")
	createCmd.Flags().StringVar(&createAssetID, "asset-id", "", "Asset ID (applied after creation)")
	createCmd.Flags().Float64Var(&createPrice, "purchase-price", 0, "Purchase price")
	createCmd.Flags().StringVar(&createPurchaseFrom, "purchase-from", "", "Purchased from")
	createCmd.Flags().Float64Var(&createSoldPrice, "sold-price", 0, "Sold price")
	createCmd.Flags().StringVar(&createNotes, "notes", "", "Notes")
	items.AddCommand(createCmd)

	// update
	var updateName, updateDesc, updateLocID, updateNotes, updatePurchaseFrom string
	var updateLabels []string
	var updateQty float64
	var updatePrice, updateSoldPrice float64
	var updateArchived bool
	var updateSerial, updateModel, updateManufacturer, updatePurchaseDate, updateSoldDate string
	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an item",
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
			if cmd.Flags().Changed("location") {
				payload["parentId"] = updateLocID
			}
			if cmd.Flags().Changed("label") {
				payload["tagIds"] = updateLabels
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
			if cmd.Flags().Changed("archived") {
				payload["archived"] = updateArchived
			}
			if cmd.Flags().Changed("serial-number") {
				payload["serialNumber"] = updateSerial
			}
			if cmd.Flags().Changed("model-number") {
				payload["modelNumber"] = updateModel
			}
			if cmd.Flags().Changed("manufacturer") {
				payload["manufacturer"] = updateManufacturer
			}
			if cmd.Flags().Changed("purchase-date") {
				payload["purchaseDate"] = updatePurchaseDate
			}
			if cmd.Flags().Changed("sold-date") {
				payload["soldDate"] = updateSoldDate
			}
			out, err := putEntity(c, args[0], payload)
			if err != nil {
				return err
			}
			client.PrintJSON(out)
			return nil
		},
	}
	updateCmd.Flags().StringVarP(&updateName, "name", "n", "", "Item name")
	updateCmd.Flags().StringVarP(&updateDesc, "description", "d", "", "Description")
	updateCmd.Flags().StringVarP(&updateLocID, "location", "l", "", "Location (parent) ID")
	updateCmd.Flags().StringArrayVar(&updateLabels, "label", nil, "Label/tag IDs (repeatable)")
	updateCmd.Flags().Float64VarP(&updateQty, "quantity", "q", 1, "Quantity")
	updateCmd.Flags().StringVar(&updateNotes, "notes", "", "Notes")
	updateCmd.Flags().Float64Var(&updatePrice, "purchase-price", 0, "Purchase price")
	updateCmd.Flags().StringVar(&updatePurchaseFrom, "purchase-from", "", "Purchased from")
	updateCmd.Flags().Float64Var(&updateSoldPrice, "sold-price", 0, "Sold price")
	updateCmd.Flags().BoolVar(&updateArchived, "archived", false, "Archive state")
	updateCmd.Flags().StringVar(&updateSerial, "serial-number", "", "Serial number")
	updateCmd.Flags().StringVar(&updateModel, "model-number", "", "Model number")
	updateCmd.Flags().StringVar(&updateManufacturer, "manufacturer", "", "Manufacturer")
	updateCmd.Flags().StringVar(&updatePurchaseDate, "purchase-date", "", "Purchase date (YYYY-MM-DD)")
	updateCmd.Flags().StringVar(&updateSoldDate, "sold-date", "", "Sold date (YYYY-MM-DD)")
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
			if _, err := c.Delete(entityBasePath + "/" + args[0]); err != nil {
				return err
			}
			fmt.Printf(`{"message": "Item %s deleted"}`+"\n", args[0])
			return nil
		},
	}
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip confirmation")
	items.AddCommand(deleteCmd)

	// duplicate
	var dupMaintenance, dupAttachments, dupFields bool
	var dupPrefix string
	duplicateCmd := &cobra.Command{
		Use:   "duplicate <id>",
		Short: "Duplicate an item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, err := c.Post(entityBasePath+"/"+args[0]+"/duplicate", map[string]any{
				"copyMaintenance":  dupMaintenance,
				"copyAttachments":  dupAttachments,
				"copyCustomFields": dupFields,
				"copyPrefix":       dupPrefix,
			})
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	duplicateCmd.Flags().BoolVar(&dupMaintenance, "copy-maintenance", false, "Copy maintenance entries")
	duplicateCmd.Flags().BoolVar(&dupAttachments, "copy-attachments", false, "Copy attachments")
	duplicateCmd.Flags().BoolVar(&dupFields, "copy-fields", false, "Copy custom fields")
	duplicateCmd.Flags().StringVar(&dupPrefix, "prefix", "", "Prefix for the copy name")
	items.AddCommand(duplicateCmd)

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
			data, err := c.Get(entityBasePath+"/"+args[0]+"/path", nil)
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
			data, err := c.Get(entityBasePath+"/"+args[0]+"/maintenance?status=both", nil)
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
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, status, err := c.Raw("GET", entityBasePath+"/export", nil, nil, "")
			if err != nil {
				return err
			}
			if status >= 400 {
				return fmt.Errorf("HTTP %d: %s", status, string(data))
			}
			if exportOutput != "" {
				if err := os.WriteFile(exportOutput, data, 0644); err != nil {
					return err
				}
				fmt.Printf(`{"message": "Exported to %s"}`+"\n", exportOutput)
			} else {
				fmt.Print(string(data))
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
			c, err := client.New(true)
			if err != nil {
				return err
			}
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()

			var buf bytes.Buffer
			w := multipart.NewWriter(&buf)
			part, _ := w.CreateFormFile("csv", filepath.Base(args[0]))
			if _, err := io.Copy(part, f); err != nil {
				return err
			}
			w.Close()

			data, status, err := c.Raw("POST", entityBasePath+"/import", nil, &buf, w.FormDataContentType())
			if err != nil {
				return err
			}
			if status >= 400 {
				return fmt.Errorf("HTTP %d: %s", status, string(data))
			}
			fmt.Println(`{"message": "Import complete"}`)
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
	var attachPrimary bool
	uploadCmd := &cobra.Command{
		Use:   "upload <item-id> <file>",
		Short: "Upload an attachment to an item",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
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
			w.WriteField("name", name) // required by v0.26
			if attachType != "" {
				w.WriteField("type", attachType)
			}
			if attachPrimary {
				w.WriteField("primary", "true")
			}
			part, _ := w.CreateFormFile("file", name)
			if _, err := io.Copy(part, f); err != nil {
				return err
			}
			w.Close()

			data, status, err := c.Raw("POST", entityBasePath+"/"+args[0]+"/attachments", nil, &buf, w.FormDataContentType())
			if err != nil {
				return err
			}
			if status >= 400 {
				return fmt.Errorf("HTTP %d: %s", status, string(data))
			}
			client.PrintJSON(data)
			return nil
		},
	}
	uploadCmd.Flags().StringVarP(&attachType, "type", "t", "", "Attachment type (photo, attachment, ...)")
	uploadCmd.Flags().StringVar(&attachName, "name", "", "Override file name")
	uploadCmd.Flags().BoolVar(&attachPrimary, "primary", false, "Set as primary attachment")
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
			if _, err := c.Delete(entityBasePath + "/" + args[0] + "/attachments/" + args[1]); err != nil {
				return err
			}
			fmt.Printf(`{"message": "Attachment %s deleted"}`+"\n", args[1])
			return nil
		},
	}
	attachDeleteCmd.Flags().BoolVarP(&attachDeleteYes, "yes", "y", false, "Skip confirmation")
	attachments.AddCommand(attachDeleteCmd)

	var attUpdType, attUpdTitle string
	var attUpdPrimary bool
	attUpdateCmd := &cobra.Command{
		Use:   "update <item-id> <attachment-id>",
		Short: "Update an item attachment (type/title/primary)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			// PUT /v1/entities/{id}/attachments/{attachment_id} is a full
			// update: fetch the current attachment so the type field (a
			// required enum) round-trips even when it was not changed.
			cur, err := fetchEntity(c, args[0])
			if err != nil {
				return err
			}
			var current map[string]any
			if atts, ok := cur["attachments"].([]any); ok {
				for _, a := range atts {
					if am, ok := a.(map[string]any); ok && stringField(am, "id") == args[1] {
						current = am
						break
					}
				}
			}
			if current == nil {
				return fmt.Errorf("attachment %s not found on item %s", args[1], args[0])
			}
			payload := map[string]any{
				"type":    stringField(current, "type"),
				"title":   stringField(current, "title"),
				"primary": current["primary"],
			}
			if cmd.Flags().Changed("type") {
				payload["type"] = attUpdType
			}
			if cmd.Flags().Changed("title") {
				payload["title"] = attUpdTitle
			}
			if cmd.Flags().Changed("primary") {
				payload["primary"] = attUpdPrimary
			}
			data, err := c.Put(entityBasePath+"/"+args[0]+"/attachments/"+args[1], payload)
			if err != nil {
				return err
			}
			client.PrintJSON(data)
			return nil
		},
	}
	attUpdateCmd.Flags().StringVarP(&attUpdType, "type", "t", "", "Attachment type")
	attUpdateCmd.Flags().StringVar(&attUpdTitle, "title", "", "Attachment title")
	attUpdateCmd.Flags().BoolVar(&attUpdPrimary, "primary", false, "Set as primary")
	attachments.AddCommand(attUpdateCmd)

	var attGetOutput string
	attGetCmd := &cobra.Command{
		Use:   "get <item-id> <attachment-id>",
		Short: "Download an item attachment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(true)
			if err != nil {
				return err
			}
			data, status, err := c.Raw("GET", entityBasePath+"/"+args[0]+"/attachments/"+args[1], nil, nil, "")
			if err != nil {
				return err
			}
			if status >= 400 {
				return fmt.Errorf("HTTP %d: %s", status, string(data))
			}
			if attGetOutput != "" {
				if err := os.WriteFile(attGetOutput, data, 0644); err != nil {
					return err
				}
				fmt.Printf(`{"message": "Attachment saved to %s"}`+"\n", attGetOutput)
			} else {
				os.Stdout.Write(data)
			}
			return nil
		},
	}
	attGetCmd.Flags().StringVarP(&attGetOutput, "output", "o", "", "Output file path")
	attachments.AddCommand(attGetCmd)
	items.AddCommand(attachments)

	// fields sub-group (v0.26: fields are managed through the entity update)
	fields := &cobra.Command{
		Use:   "fields",
		Short: "Manage item custom fields",
	}

	fields.AddCommand(&cobra.Command{
		Use:   "list <item-id>",
		Short: "List custom fields of an item",
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
			fieldsList, _ := cur["fields"]
			out, _ := json.MarshalIndent(fieldsList, "", "  ")
			fmt.Println(string(out))
			return nil
		},
	})

	var fieldLabel, fieldValue, fieldType string
	fieldAddCmd := &cobra.Command{
		Use:   "add <item-id>",
		Short: "Add a custom field to an item",
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
			fieldsList, _ := payload["fields"].([]any)
			fieldsList = append(fieldsList, fieldData(fieldType, fieldLabel, fieldValue, ""))
			payload["fields"] = fieldsList
			out, err := putEntity(c, args[0], payload)
			if err != nil {
				return err
			}
			client.PrintJSON(out)
			return nil
		},
	}
	fieldAddCmd.Flags().StringVarP(&fieldLabel, "label", "l", "", "Field label")
	fieldAddCmd.Flags().StringVarP(&fieldValue, "value", "v", "", "Field value")
	fieldAddCmd.Flags().StringVarP(&fieldType, "type", "t", "text", "Field type (text, number, boolean)")
	fieldAddCmd.MarkFlagRequired("label")
	fieldAddCmd.MarkFlagRequired("value")
	fields.AddCommand(fieldAddCmd)

	fieldUpdateCmd := &cobra.Command{
		Use:   "update <item-id> <field-id>",
		Short: "Update a custom field of an item",
		Args:  cobra.ExactArgs(2),
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
			fieldsList, _ := payload["fields"].([]any)
			found := false
			for i, f := range fieldsList {
				if fm, ok := f.(map[string]any); ok && stringField(fm, "id") == args[1] {
					if cmd.Flags().Changed("label") {
						fm["name"] = fieldLabel
					}
					if cmd.Flags().Changed("value") {
						setFieldValue(fm, fieldType, fieldValue)
					}
					if cmd.Flags().Changed("type") {
						fm["type"] = fieldType
					}
					fieldsList[i] = fm
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("field %s not found on item %s", args[1], args[0])
			}
			payload["fields"] = fieldsList
			out, err := putEntity(c, args[0], payload)
			if err != nil {
				return err
			}
			client.PrintJSON(out)
			return nil
		},
	}
	fieldUpdateCmd.Flags().StringVarP(&fieldLabel, "label", "l", "", "Field label")
	fieldUpdateCmd.Flags().StringVarP(&fieldValue, "value", "v", "", "Field value")
	fieldUpdateCmd.Flags().StringVarP(&fieldType, "type", "t", "text", "Field type (text, number, boolean)")
	fields.AddCommand(fieldUpdateCmd)

	var fieldDeleteYes bool
	fieldDeleteCmd := &cobra.Command{
		Use:   "delete <item-id> <field-id>",
		Short: "Delete a custom field from an item",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !fieldDeleteYes {
				if !confirm("Delete field " + args[1] + "?") {
					return nil
				}
			}
			c, err := client.New(true)
			if err != nil {
				return err
			}
			cur, err := fetchEntity(c, args[0])
			if err != nil {
				return err
			}
			payload := updatePayload(cur)
			fieldsList, _ := payload["fields"].([]any)
			outList := fieldsList[:0]
			for _, f := range fieldsList {
				if fm, ok := f.(map[string]any); ok && stringField(fm, "id") == args[1] {
					continue
				}
				outList = append(outList, f)
			}
			payload["fields"] = outList
			out, err := putEntity(c, args[0], payload)
			if err != nil {
				return err
			}
			fmt.Printf(`{"message": "Field %s deleted"}`+"\n", args[1])
			_ = out
			return nil
		},
	}
	fieldDeleteCmd.Flags().BoolVarP(&fieldDeleteYes, "yes", "y", false, "Skip confirmation")
	fields.AddCommand(fieldDeleteCmd)

	items.AddCommand(fields)

	return items
}

// fieldData builds an EntityFieldData map for the given label/value/type.
func fieldData(fieldType, label, value, id string) map[string]any {
	f := map[string]any{
		"type":         fieldType,
		"name":         label,
		"textValue":    "",
		"numberValue":  0,
		"booleanValue": false,
	}
	if id != "" {
		f["id"] = id
	}
	setFieldValue(f, fieldType, value)
	return f
}

// setFieldValue stores value in the correct field of an EntityFieldData map
// depending on the field type.
func setFieldValue(f map[string]any, fieldType, value string) {
	switch strings.ToLower(fieldType) {
	case "number":
		if n, err := strconv.Atoi(value); err == nil {
			f["numberValue"] = n
		} else if fl, err := strconv.ParseFloat(value, 64); err == nil {
			f["numberValue"] = int(fl)
		}
	case "boolean":
		f["booleanValue"] = strings.EqualFold(value, "true") || value == "1"
	default:
		f["type"] = "text"
		f["textValue"] = value
		f["numberValue"] = 0
		f["booleanValue"] = false
	}
}
