package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hb",
	Short: "Homebox REST API CLI",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(
		newLoginCmd(),
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
