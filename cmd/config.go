package cmd

import (
	stdErrors "errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/intility/cwc/pkg/config"
	"github.com/intility/cwc/pkg/errors"
	cwcui "github.com/intility/cwc/pkg/ui"
)

func createConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Get or set config variables",
		Long:  `Get or set config variables`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Usage()
			if err != nil {
				return fmt.Errorf("failed to print usage: %w", err)
			}

			return nil
		},
	}

	cmd.AddCommand(createGetConfigCommand())
	cmd.AddCommand(createSetConfigCommand())

	return cmd
}

func createGetConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Print current config",
		Long:  "Print current config",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := config.NewDefaultProvider()
			cfg, err := provider.GetConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			printConfig(cfg)

			return nil
		},
	}

	return cmd
}

func createSetConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set config variables",
		Long:  "Set config variables",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgProvider := config.NewDefaultProvider()
			cfg, err := cfgProvider.GetConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// if no args are given, print the help and exit
			if len(args) == 0 {
				err = cmd.Help()
				if err != nil {
					return fmt.Errorf("failed to print help: %w", err)
				}

				return nil
			}

			err = processKeyValuePairs(cfg, args)

			if err != nil {
				var suppressedError errors.SuppressedError
				if ok := stdErrors.As(err, &suppressedError); ok {
					cmd.SilenceUsage = true
					cmd.SilenceErrors = true
				}

				return err
			}

			return nil
		},
	}

	return cmd
}

func processKeyValuePairs(cfg *config.Config, kvPairs []string) error {
	// iterate over each argument and process them as key=value pairs
	argKvSubstrCount := 2
	for _, arg := range kvPairs {
		kvPair := strings.SplitN(arg, "=", argKvSubstrCount)
		if len(kvPair) != argKvSubstrCount {
			return errors.ArgParseError{Message: fmt.Sprintf("invalid argument format: %s, expected key=value", arg)}
		}

		key := kvPair[0]
		value := kvPair[1]

		err := setConfigValue(cfg, key, value)
		if err != nil {
			return fmt.Errorf("failed to set config value: %w", err)
		}
	}

	provider := config.NewDefaultProvider()

	err := provider.SaveConfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	printConfig(cfg)

	return nil
}

func setConfigValue(cfg *config.Config, key, value string) error {
	ui := cwcui.NewUI() //nolint:varnamelen

	switch key {
	case "endpoint":
		cfg.Endpoint = value
	case "deploymentName":
		cfg.ModelDeployment = value
	case "apiKey":
		cfg.SetAPIKey(value)
	case "useGitignore":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return errors.ArgParseError{Message: "invalid boolean value for useGitignore: " + value}
		}

		cfg.UseGitignore = b
	case "excludeGitDir":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return errors.ArgParseError{Message: "invalid boolean value for excludeGitDir: " + value}
		}

		cfg.ExcludeGitDir = b
	default:
		ui.PrintMessage(fmt.Sprintf("Unknown config key: %s\n", key), cwcui.MessageTypeError)

		validKeys := []string{
			"endpoint",
			"deploymentName",
			"apiKey",
			"useGitignore",
			"excludeGitDir",
		}

		ui.PrintMessage("Valid keys are: "+strings.Join(validKeys, ", "), cwcui.MessageTypeInfo)

		return errors.SuppressedError{}
	}

	return nil
}

func printConfig(cfg *config.Config) {
	table := [][]string{
		{"Name", "Value"},
		{"endpoint", cfg.Endpoint},
		{"deploymentName", cfg.ModelDeployment},
		{"apiKey", cfg.APIKey()},
		{"SEP", ""},
		{"useGitignore", fmt.Sprintf("%t", cfg.UseGitignore)},
		{"excludeGitDir", fmt.Sprintf("%t", cfg.ExcludeGitDir)},
	}

	printTable(table)
}

func printTable(table [][]string) {
	ui := cwcui.NewUI() //nolint:varnamelen
	columnLengths := calculateColumnLengths(table)

	var lineLength int

	additionalChars := 3 // +3 for 3 additional characters before and after each field: "| %s "
	for _, c := range columnLengths {
		lineLength += c + additionalChars // +3 for 3 additional characters before and after each field: "| %s "
	}

	lineLength++                               // +1 for the last "|" in the line
	singleLineLength := lineLength - len("++") // -2 because of "+" as first and last character

	for lineIndex, line := range table {
		if lineIndex == 0 { // table header
			// lineLength-2 because of "+" as first and last charactr
			ui.PrintMessage(fmt.Sprintf("+%s+\n", strings.Repeat("-", singleLineLength)), cwcui.MessageTypeInfo)
		}

	lineLoop:
		for rowIndex, val := range line {
			if val == "SEP" {
				// lineLength-2 because of "+" as first and last character
				ui.PrintMessage(fmt.Sprintf("+%s+\n", strings.Repeat("-", singleLineLength)), cwcui.MessageTypeInfo)
				break lineLoop
			}

			ui.PrintMessage(fmt.Sprintf("| %-*s ", columnLengths[rowIndex], val), cwcui.MessageTypeInfo)
			if rowIndex == len(line)-1 {
				ui.PrintMessage("|\n", cwcui.MessageTypeInfo)
			}
		}

		if lineIndex == 0 || lineIndex == len(table)-1 { // table header or last line
			// lineLength-2 because of "+" as first and last character
			ui.PrintMessage(fmt.Sprintf("+%s+\n", strings.Repeat("-", singleLineLength)), cwcui.MessageTypeInfo)
		}
	}
}

func calculateColumnLengths(table [][]string) []int {
	columnLengths := make([]int, len(table[0]))

	for _, line := range table {
		for i, val := range line {
			if len(val) > columnLengths[i] {
				columnLengths[i] = len(val)
			}
		}
	}

	return columnLengths
}
