package config

import (
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

const (
	APIKeyEnvVar       = "AOAI_API_KEY" // #nosec
	EndpointEnvVar     = "AOAI_ENDPOINT"
	APIVersionEnvVar   = "AOAI_API_VERSION"
	ModelDeploymentVar = "AOAI_MODEL_DEPLOYMENT"
)

type missingEnvVarError struct {
	// missing is a list of environment variables that are required but missing
	missing []string
}

func (e missingEnvVarError) Error() string {
	return fmt.Sprintf("missing required environment variables: %v", e.missing)
}

func NewFromEnv() (openai.ClientConfig, error) {

	// check for the presence of the required environment variables
	apiKey := os.Getenv(APIKeyEnvVar)
	endpoint := os.Getenv(EndpointEnvVar)
	apiVersion := os.Getenv(APIVersionEnvVar)
	modelDeployment := os.Getenv(ModelDeploymentVar)

	// if any of the required environment variables are missing, return an error
	requiredSlice := []string{apiKey, endpoint, apiVersion, modelDeployment}
	missingSlice := []string{}
	for i, v := range requiredSlice {
		if v == "" {
			missingSlice = append(missingSlice, requiredSlice[i])
		}
	}
	if len(missingSlice) > 0 {
		return openai.ClientConfig{}, missingEnvVarError{missing: missingSlice}
	}

	config := openai.DefaultAzureConfig(apiKey, endpoint)

	config.APIVersion = apiVersion

	// If you use a deployment name different from the model name, you can customize the AzureModelMapperFunc function
	config.AzureModelMapperFunc = func(model string) string {
		azureModelMapping := map[string]string{
			openai.GPT4TurboPreview: modelDeployment,
		}

		return azureModelMapping[model]
	}

	return config, nil
}
