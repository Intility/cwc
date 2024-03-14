package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/emilkje/cwc/pkg/config"
	"github.com/emilkje/cwc/pkg/errors"
	"github.com/emilkje/cwc/pkg/ui"
)

var (
	apiKeyFlag          string //nolint:gochecknoglobals
	endpointFlag        string //nolint:gochecknoglobals
	apiVersionFlag      string //nolint:gochecknoglobals
	modelDeploymentFlag string //nolint:gochecknoglobals
)

func createLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Azure OpenAI",
		Long: "Login will prompt you to enter your Azure OpenAI API key " +
			"and other relevant information required for authentication.\n" +
			"Your credentials will be stored securely in your keyring and will never be exposed on the file system directly.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Prompt for other required authentication details (apiKey, endpoint, version, and deployment)
			if apiKeyFlag == "" {
				ui.PrintMessage("Enter the Azure OpenAI API Key: ", ui.MessageTypeInfo)
				apiKeyFlag = config.SanitizeInput(ui.ReadUserInput())
			}

			if endpointFlag == "" {
				ui.PrintMessage("Enter the Azure OpenAI API Endpoint: ", ui.MessageTypeInfo)
				endpointFlag = config.SanitizeInput(ui.ReadUserInput())
			}

			if apiVersionFlag == "" {
				ui.PrintMessage("Enter the Azure OpenAI API Version: ", ui.MessageTypeInfo)
				apiVersionFlag = config.SanitizeInput(ui.ReadUserInput())
			}

			if modelDeploymentFlag == "" {
				ui.PrintMessage("Enter the Azure OpenAI Model Deployment: ", ui.MessageTypeInfo)
				modelDeploymentFlag = config.SanitizeInput(ui.ReadUserInput())
			}

			cfg := config.NewConfig(endpointFlag, apiVersionFlag, modelDeploymentFlag)
			cfg.SetAPIKey(apiKeyFlag)

			err := config.SaveConfig(cfg)
			if err != nil {
				if validationErr, ok := errors.AsConfigValidationError(err); ok {
					for _, e := range validationErr.Errors {
						ui.PrintMessage(e+"\n", ui.MessageTypeError)
					}

					return nil // suppress the error
				}

				return fmt.Errorf("error saving configuration: %w", err)
			}

			ui.PrintMessage("config saved successfully\n", ui.MessageTypeSuccess)

			return nil
		},
	}

	cmd.Flags().StringVarP(&apiKeyFlag, "api-key", "k", "", "Azure OpenAI API Key")
	cmd.Flags().StringVarP(&endpointFlag, "endpoint", "e", "", "Azure OpenAI API Endpoint")
	cmd.Flags().StringVarP(&apiVersionFlag, "api-version", "v", "", "Azure OpenAI API Version")
	cmd.Flags().StringVarP(&modelDeploymentFlag, "model-deployment", "m", "", "Azure OpenAI Model Deployment")

	return cmd
}
