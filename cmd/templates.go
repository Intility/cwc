package cmd

import (
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/intility/cwc/pkg/config"
	"github.com/intility/cwc/pkg/templates"
	cwcui "github.com/intility/cwc/pkg/ui"
)

type Template struct {
	template           templates.Template
	placement          string
	isOverridingGlobal bool
}

func createTemplatesCmd() *cobra.Command {
	ui := cwcui.NewUI() //nolint:varnamelen
	cmd := &cobra.Command{
		Use:   "templates",
		Short: "Lists available templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			tmpls := locateTemplates()

			if len(tmpls) == 0 {
				ui.PrintMessage("No templates found", cwcui.MessageTypeInfo)
				return nil
			}

			if cfgDir, err := config.GetConfigDir(); err == nil {
				ui.PrintMessage("global", cwcui.MessageTypeWarning)
				ui.PrintMessage(": the template is defined in "+
					filepath.Join(cfgDir, "templates.yaml")+"\n", cwcui.MessageTypeInfo)
			}

			ui.PrintMessage("local", cwcui.MessageTypeSuccess)
			ui.PrintMessage(": the template is defined in ./cwc/templates.yaml\n", cwcui.MessageTypeInfo)

			ui.PrintMessage("overridden", cwcui.MessageTypeError)
			ui.PrintMessage(": the local template is overriding a global template with the same name\n\n", cwcui.MessageTypeInfo)

			ui.PrintMessage("Available templates:\n", cwcui.MessageTypeInfo)

			for _, template := range tmpls {
				if template.isOverridingGlobal {
					template.placement = "overridden"
				}

				placementMessageType := cwcui.MessageTypeSuccess
				if template.placement == "global" {
					placementMessageType = cwcui.MessageTypeWarning
				}

				if template.isOverridingGlobal {
					placementMessageType = cwcui.MessageTypeError
				}

				ui.PrintMessage("- name: ", cwcui.MessageTypeInfo)
				ui.PrintMessage(template.template.Name, cwcui.MessageTypeInfo)
				ui.PrintMessage(" ("+template.placement+")\n", placementMessageType)
				printTemplateInfo(template.template)
			}

			return nil
		},
	}

	return cmd
}

func locateTemplates() map[string]Template {
	var localTemplates, globalTemplates []templates.Template

	cfgDir, err := config.GetConfigDir()
	if err == nil {
		globalTemplatesLocator := templates.NewYamlFileTemplateLocator(filepath.Join(cfgDir, "templates.yaml"))
		locatedTemplates, err := globalTemplatesLocator.ListTemplates()

		if err == nil {
			globalTemplates = locatedTemplates
		}
	}

	localTemplatesLocator := templates.NewYamlFileTemplateLocator(filepath.Join(".cwc", "templates.yaml"))
	locatedTemplates, err := localTemplatesLocator.ListTemplates()

	if err == nil {
		localTemplates = locatedTemplates
	}

	tmpls := make(map[string]Template)

	// populate the list of templates, marking the local ones as overriding the global ones if they have the same name
	for _, t := range globalTemplates {
		tmpls[t.Name] = Template{template: t, placement: "global", isOverridingGlobal: false}
	}

	for _, t := range localTemplates {
		_, exists := tmpls[t.Name]
		tmpls[t.Name] = Template{template: t, placement: "local", isOverridingGlobal: exists}
	}

	return tmpls
}

func printTemplateInfo(template templates.Template) {
	ui := cwcui.NewUI() //nolint:varnamelen
	ui.PrintMessage("  description: "+template.Description+"\n", cwcui.MessageTypeInfo)

	dfp := "no"
	if template.DefaultPrompt != "" {
		dfp = "yes"
	}

	ui.PrintMessage("  has_default_prompt: "+dfp+"\n", cwcui.MessageTypeInfo)

	variablesCount := len(template.Variables)

	ui.PrintMessage("  variables: "+strconv.Itoa(variablesCount)+"\n", cwcui.MessageTypeInfo)

	for _, variable := range template.Variables {
		ui.PrintMessage("  - name: ", cwcui.MessageTypeInfo)
		ui.PrintMessage(variable.Name, cwcui.MessageTypeInfo)
		ui.PrintMessage("\n", cwcui.MessageTypeInfo)
		ui.PrintMessage("    description: "+variable.Description+"\n", cwcui.MessageTypeInfo)

		dv := "no"
		if variable.DefaultValue != "" {
			dv = "yes"
		}

		ui.PrintMessage("    has_default_value: "+dv+"\n", cwcui.MessageTypeInfo)
	}

	ui.PrintMessage("\n", cwcui.MessageTypeInfo)
}
