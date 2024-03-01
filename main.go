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

//
//func Execute(defCmd string) error {
//	var cmdFound bool
//	commands := cmd.RootCmd.Commands()
//
//	for _, a := range commands {
//		for _, b := range os.Args[1:] {
//			if a.Name() == b {
//				cmdFound = true
//				break
//			}
//		}
//	}
//	if !cmdFound {
//		args := append([]string{defCmd}, os.Args[1:]...)
//		cmd.RootCmd.SetArgs(args)
//	}
//	if err := cmd.RootCmd.Execute(); err != nil {
//		return err
//	}
//
//	return nil
//}
