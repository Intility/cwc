package cmd

import (
	"path/filepath"
	"strconv"

	"github.com/intility/cwc/pkg/config"
	"github.com/intility/cwc/pkg/templates"
	"github.com/intility/cwc/pkg/ui"
	"github.com/spf13/cobra"
)

type Template struct {
	template           templates.Template
	placement          string
	isOverridingGlobal bool
}

func createTemplatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "templates",
		Short: "Lists available templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			tmpls := locateTemplates()

			if len(tmpls) == 0 {
				ui.PrintMessage("No templates found", ui.MessageTypeInfo)
				return nil
			}

			if cfgDir, err := config.GetConfigDir(); err == nil {
				ui.PrintMessage("global", ui.MessageTypeWarning)
				ui.PrintMessage(": the template is defined in "+
					filepath.Join(cfgDir, "templates.yaml")+"\n", ui.MessageTypeInfo)
			}

			ui.PrintMessage("local", ui.MessageTypeSuccess)
			ui.PrintMessage(": the template is defined in ./cwc/templates.yaml\n", ui.MessageTypeInfo)

			ui.PrintMessage("overridden", ui.MessageTypeError)
			ui.PrintMessage(": the local template is overriding a global template with the same name\n\n", ui.MessageTypeInfo)

			ui.PrintMessage("Available templates:\n", ui.MessageTypeInfo)

			for _, template := range tmpls {
				if template.isOverridingGlobal {
					template.placement = "overridden"
				}

				placementMessageType := ui.MessageTypeSuccess
				if template.placement == "global" {
					placementMessageType = ui.MessageTypeWarning
				}

				if template.isOverridingGlobal {
					placementMessageType = ui.MessageTypeError
				}

				ui.PrintMessage("- name: ", ui.MessageTypeInfo)
				ui.PrintMessage(template.template.Name, ui.MessageTypeInfo)
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
	ui.PrintMessage("  description: "+template.Description+"\n", ui.MessageTypeInfo)

	dfp := "no"
	if template.DefaultPrompt != "" {
		dfp = "yes"
	}

	ui.PrintMessage("  has_default_prompt: "+dfp+"\n", ui.MessageTypeInfo)

	variablesCount := len(template.Variables)

	ui.PrintMessage("  variables: "+strconv.Itoa(variablesCount)+"\n", ui.MessageTypeInfo)

	for _, variable := range template.Variables {
		ui.PrintMessage("  - name: ", ui.MessageTypeInfo)
		ui.PrintMessage(variable.Name, ui.MessageTypeInfo)
		ui.PrintMessage("\n", ui.MessageTypeInfo)
		ui.PrintMessage("    description: "+variable.Description+"\n", ui.MessageTypeInfo)

		dv := "no"
		if variable.DefaultValue != "" {
			dv = "yes"
		}

		ui.PrintMessage("    has_default_value: "+dv+"\n", ui.MessageTypeInfo)
	}

	ui.PrintMessage("\n", ui.MessageTypeInfo)
}
