package main

import (
	stdErrors "errors"
	"fmt"
	"os"

	"github.com/intility/cwc/cmd"
	"github.com/intility/cwc/pkg/errors"
	"github.com/intility/cwc/pkg/ui"
)

//go:generate ./bin/lang-gen

func main() {
	command := cmd.CreateRootCommand()

	err := command.Execute()
	if err != nil {
		// if error is of type suppressedError, do not print error message
		var suppressedError errors.SuppressedError
		if ok := stdErrors.As(err, &suppressedError); !ok {
			ui.PrintMessage(fmt.Sprintf("Error: %s\n", err), ui.MessageTypeError)
		}

		os.Exit(1)
	}
}
