package cmd

import (
	"github.com/spf13/cobra"

	"github.com/emilkje/cwc/pkg/config"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear the configuration and remove the stored API key",
	Long: `Logout will clear the configuration and remove the stored API key.
This will require you to login again to use the chat with context tool.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := config.ClearConfig()

		if err != nil {
			return err
		}
		return nil
	},
}
