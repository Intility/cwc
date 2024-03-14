package main

import (
	"fmt"
	"github.com/emilkje/cwc/pkg/ui"
	"os"

	"github.com/emilkje/cwc/cmd"
)

//go:generate ./bin/lang-gen

func main() {

	err := cmd.CwcCmd.Execute()

	if err != nil {
		ui.PrintMessage(fmt.Sprintf("Error: %s\n", err), ui.MessageTypeError)
		os.Exit(1)
	}
}
