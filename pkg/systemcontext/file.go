package systemcontext

import (
	"fmt"

	"github.com/intility/cwc/pkg/config"
	"github.com/intility/cwc/pkg/errors"
	"github.com/intility/cwc/pkg/filetree"
	"github.com/intility/cwc/pkg/pathmatcher"
	"github.com/intility/cwc/pkg/ui"
)

type FileContextRetriever struct {
	ui             ui.UI
	cfgProvider    config.Provider
	includePattern string
	excludePattern string
	searchScopes   []string
	contextPrinter func(fileTree string, files []filetree.File)
}

type FileContextRetrieverOptions struct {
	CfgProvider    config.Provider
	IncludePattern string
	ExcludePattern string
	SearchScopes   []string
	ContextPrinter func(fileTree string, files []filetree.File)
}

func NewFileContextRetriever(opts FileContextRetrieverOptions) *FileContextRetriever {
	return &FileContextRetriever{
		ui:             ui.NewUI(),
		cfgProvider:    opts.CfgProvider,
		includePattern: opts.IncludePattern,
		excludePattern: opts.ExcludePattern,
		searchScopes:   opts.SearchScopes,
		contextPrinter: opts.ContextPrinter,
	}
}

func (r *FileContextRetriever) RetrieveContext() (string, error) {
	files, rootNode, err := r.gatherContext()
	if err != nil {
		return "", fmt.Errorf("error gathering context: %w", err)
	}

	fileTree := filetree.GenerateFileTree(rootNode, "", true)

	if r.contextPrinter != nil {
		r.contextPrinter(fileTree, files)
	}

	ctx := r.createContext(fileTree, files)

	return ctx, nil
}

func (r *FileContextRetriever) createContext(fileTree string, files []filetree.File) string {
	contextStr := "File tree:\n\n"
	contextStr += "```\n" + fileTree + "```\n\n"
	contextStr += "File contents:\n\n"

	for _, file := range files {
		// find extension by splitting on ".". if no extension, use
		contextStr += fmt.Sprintf("./%s\n```%s\n%s\n```\n\n", file.Path, file.Type, file.Data)
	}

	return contextStr
}

func (r *FileContextRetriever) gatherContext() ([]filetree.File, *filetree.FileNode, error) {
	var excludeMatchers []pathmatcher.PathMatcher

	// add exclude flag to excludeMatchers
	if r.excludePattern != "" {
		excludeMatcher, err := pathmatcher.NewRegexPathMatcher(r.excludePattern)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating exclude matcher: %w", err)
		}

		excludeMatchers = append(excludeMatchers, excludeMatcher)
	}

	excludeMatchersFromConfig, err := r.excludeMatchersFromConfig()
	if err != nil {
		return nil, nil, err
	}

	excludeMatchers = append(excludeMatchers, excludeMatchersFromConfig...)

	excludeMatcher := pathmatcher.NewCompoundPathMatcher(excludeMatchers...)

	includeMatcher, err := pathmatcher.NewRegexPathMatcher(r.includePattern)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating include matcher: %w", err)
	}

	files, rootNode, err := filetree.GatherFiles(&filetree.FileGatherOptions{
		IncludeMatcher: includeMatcher,
		ExcludeMatcher: excludeMatcher,
		PathScopes:     r.searchScopes,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error gathering files: %w", err)
	}

	return files, rootNode, nil
}

func (r *FileContextRetriever) excludeMatchersFromConfig() ([]pathmatcher.PathMatcher, error) {
	var excludeMatchers []pathmatcher.PathMatcher

	cfg, err := r.cfgProvider.GetConfig()
	if err != nil {
		return excludeMatchers, fmt.Errorf("error loading config: %w", err)
	}

	if cfg.UseGitignore {
		gitignoreMatcher, err := pathmatcher.NewGitignorePathMatcher()
		if err != nil {
			switch {
			case errors.IsGitNotInstalledError(err):
				r.ui.PrintMessage("warning: git not found in PATH, skipping .gitignore\n", ui.MessageTypeWarning)
			case errors.IsNotAGitRepositoryError(err):
				r.ui.PrintMessage("warning: not a git repository, skipping .gitignore\n", ui.MessageTypeWarning)
			default:
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
