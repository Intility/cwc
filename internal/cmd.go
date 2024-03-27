package internal

import (
	"fmt"

	"github.com/intility/cwc/pkg/config"
	"github.com/intility/cwc/pkg/errors"
	"github.com/intility/cwc/pkg/filetree"
	"github.com/intility/cwc/pkg/pathmatcher"
	"github.com/intility/cwc/pkg/ui"
)

const (
	warnFileSizeThreshold = 100000
)

// askConfirmation prompts the user if they want to proceed with no files.
func askConfirmation(prompt string, messageType ui.MessageType) bool {
	ui.PrintMessage(prompt, messageType)

	if !ui.AskYesNo("Do you wish to proceed?", false) {
		ui.PrintMessage("See ya later!", ui.MessageTypeInfo)
		return false
	}

	return true
}

func excludeMatchersFromConfig() ([]pathmatcher.PathMatcher, error) {
	var excludeMatchers []pathmatcher.PathMatcher

	cfg, err := config.LoadConfig()
	if err != nil {
		return excludeMatchers, fmt.Errorf("error loading config: %w", err)
	}

	if cfg.UseGitignore {
		gitignoreMatcher, err := pathmatcher.NewGitignorePathMatcher()
		if err != nil {
			if errors.IsGitNotInstalledError(err) {
				ui.PrintMessage("warning: git not found in PATH, skipping .gitignore\n", ui.MessageTypeWarning)
			} else {
				return nil, fmt.Errorf("error creating gitignore matcher: %w", err)
			}
		}

		excludeMatchers = append(excludeMatchers, gitignoreMatcher)
	}

	if cfg.ExcludeGitDir {
		gitDirMatcher, err := pathmatcher.NewRegexPathMatcher(`^\.git(/|\\)`)
		if err != nil {
			return nil, fmt.Errorf("error creating git directory matcher: %w", err)
		}

		excludeMatchers = append(excludeMatchers, gitDirMatcher)
	}

	return excludeMatchers, nil
}

func printLargeFileWarning(file filetree.File) {
	if len(file.Data) > warnFileSizeThreshold {
		largeFileMsg := fmt.Sprintf(
			"warning: %s is very large (%d bytes) and will degrade performance.\n",
			file.Path, len(file.Data))

		ui.PrintMessage(largeFileMsg, ui.MessageTypeWarning)
	}
}

func createContext(fileTree string, files []filetree.File) string {
	contextStr := "File tree:\n\n"
	contextStr += "```\n" + fileTree + "```\n\n"
	contextStr += "File contents:\n\n"

	for _, file := range files {
		// find extension by splitting on ".". if no extension, use
		contextStr += fmt.Sprintf("./%s\n```%s\n%s\n```\n\n", file.Path, file.Type, file.Data)
	}

	return contextStr
}
