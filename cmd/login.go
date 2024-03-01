package cmd

import (
	"github.com/emilkje/cwc/pkg/config"
	"github.com/emilkje/cwc/pkg/errors"
	"github.com/emilkje/cwc/pkg/ui"
	"github.com/spf13/cobra"
)

const (
	serviceName = "cwc"
)

var (
	apiKeyFlag          string
	endpointFlag        string
	apiVersionFlag      string
	modelDeploymentFlag string
)

// Declaration of a new cobra command for login
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Azure OpenAI",
	Long: `Login will prompt you to enter your Azure OpenAI API key and other relevant information required for authentication. 
Your credentials will be stored securely in your keyring and will never be exposed on the file system directly.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Prompt for other required authentication details (apiKey, endpoint, version, and deployment)
		apiKey := apiKeyFlag
		endpoint := endpointFlag
		apiVersion := apiVersionFlag
		modelDeployment := modelDeploymentFlag

		if apiKeyFlag == "" {
			ui.PrintMessage("Enter the Azure OpenAI API Key: ", ui.MessageTypeInfo)
			apiKey = config.SanitizeInput(ui.ReadUserInput())
		}

		if endpointFlag == "" {
			ui.PrintMessage("Enter the Azure OpenAI API Endpoint: ", ui.MessageTypeInfo)
			endpoint = config.SanitizeInput(ui.ReadUserInput())
		}

		if apiVersionFlag == "" {
			ui.PrintMessage("Enter the Azure OpenAI API Version: ", ui.MessageTypeInfo)
			apiVersion = config.SanitizeInput(ui.ReadUserInput())
		}

		if modelDeploymentFlag == "" {
			ui.PrintMessage("Enter the Azure OpenAI Model Deployment: ", ui.MessageTypeInfo)
			modelDeployment = config.SanitizeInput(ui.ReadUserInput())
		}

		cfg := config.NewConfig(endpoint, apiVersion, modelDeployment)
		cfg.SetAPIKey(apiKey)

		err := config.SaveConfig(cfg)
		if err != nil {
			if validationErr, ok := errors.AsConfigValidationError(err); ok {
				for _, e := range validationErr.Errors {
					ui.PrintMessage(e+"\n", ui.MessageTypeError)
				}
				return nil // suppress the error
			}
			return err
		}

		ui.PrintMessage("config saved successfully\n", ui.MessageTypeSuccess)

		return nil
	},
}

func init() {

	// Add flags to the login command
	loginCmd.Flags().StringVarP(&apiKeyFlag, "api-key", "k", "", "Azure OpenAI API Key")
	loginCmd.Flags().StringVarP(&endpointFlag, "endpoint", "e", "", "Azure OpenAI API Endpoint")
	loginCmd.Flags().StringVarP(&apiVersionFlag, "api-version", "v", "", "Azure OpenAI API Version")
	loginCmd.Flags().StringVarP(&modelDeploymentFlag, "model-deployment", "m", "", "Azure OpenAI Model Deployment")
}
