package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/emilkje/cwc/pkg/config"
)

func createLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Clear the configuration and remove the stored API key",
		Long: `Logout will clear the configuration and remove the stored API key.
This will require you to login again to use the chat with context tool.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := config.ClearConfig()
			if err != nil {
				return fmt.Errorf("error clearing configuration: %w", err)
			}

			return nil
		},
	}

	return cmd
}
