package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/misonikomipan/homebox-cli/internal/client"
	"github.com/spf13/cobra"
)

func newLabelmakerCmd() *cobra.Command {
	lm := &cobra.Command{
		Use:   "labelmaker",
		Short: "Generate labels (v0.26 labelmaker endpoints)",
	}

	var labelType, output string
	var print bool
	getCmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Generate a label PNG for an entity/item/location/asset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			t := strings.ToLower(labelType)
			switch t {
			case "entity", "item", "location", "asset":
			default:
				return fmt.Errorf("--type must be one of: entity, item, location, asset")
			}
			c, err := client.New(true)
			if err != nil {
				return err
			}
			q := url.Values{}
			if print {
				q.Set("print", "true")
			}
			data, status, err := c.Raw("GET", "/v1/labelmaker/"+t+"/"+args[0], q, nil, "")
			if err != nil {
				return err
			}
			if status >= 400 {
				return fmt.Errorf("HTTP %d: %s", status, string(data))
			}
			if output != "" {
				if err := os.WriteFile(output, data, 0644); err != nil {
					return err
				}
				fmt.Printf(`{"message": "Label saved to %s"}`+"\n", output)
			} else {
				os.Stdout.Write(data)
			}
			return nil
		},
	}
	getCmd.Flags().StringVarP(&labelType, "type", "t", "entity", "Label type (entity, item, location, asset)")
	getCmd.Flags().StringVarP(&output, "output", "o", "", "Output PNG file path")
	getCmd.Flags().BoolVar(&print, "print", false, "Send the label to the configured printer")
	lm.AddCommand(getCmd)

	return lm
}
