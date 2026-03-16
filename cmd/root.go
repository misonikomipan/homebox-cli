package cmd

import (
	"github.com/misonikomipan/homebox-cli/internal/config"
	"github.com/spf13/cobra"
)

var outputFormat string

var rootCmd = &cobra.Command{
	Use:   "hb",
	Short: "Homebox REST API CLI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("format") {
			return config.SetFormat(outputFormat)
		}
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "json", "Output format (json, table)")

	rootCmd.AddCommand(		newLoginCmd(),
		newLogoutCmd(),
		newStatusCmd(),
		newConfigCmd(),
		newGuideCmd(),
		newCurrencyCmd(),
		newBarcodeSearchCmd(),
		newItemsCmd(),
		newLocationsCmd(),
		newTagsCmd(),
		newGroupsCmd(),
		newMaintenanceCmd(),
		newNotifiersCmd(),
		newTemplatesCmd(),
		newAuthCmd(),
	)
}
