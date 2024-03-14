package filetree

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	pm "github.com/emilkje/cwc/pkg/pathmatcher"
	"github.com/emilkje/cwc/pkg/ui"
)

type FileNode struct {
	Name     string
	IsDir    bool
	Children []*FileNode
}

type File struct {
	Path string
	Data []byte
	Type string
}

type FileGatherOptions struct {
	IncludeMatcher pm.PathMatcher
	ExcludeMatcher pm.PathMatcher
	PathScopes     []string
}

func GatherFiles(opts *FileGatherOptions) ([]File, *FileNode, error) { //nolint:funlen,gocognit,cyclop
	includeMatcher := opts.IncludeMatcher
	excludeMatcher := opts.ExcludeMatcher
	pathScopes := opts.PathScopes

	var files []File

	knownLanguage := cachedLanguageChecker()

	rootNode := &FileNode{Name: "/", IsDir: true, Children: []*FileNode{}}

	for _, path := range pathScopes {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// start by skipping the .git directory
			if strings.HasPrefix(path, ".git/") {
				return nil
			}

			if !includeMatcher.Match(path) || info.IsDir() || excludeMatcher.Match(path) {
				return nil
			}

			fileType, ok := knownLanguage(path)

			if !ok {
				ui.PrintMessage("skipping unknown file type: "+path+"\n", ui.MessageTypeWarning)
				return nil
			}

			file := &File{
				Path: path,
				Type: fileType,
				Data: []byte{},
			}

			codeFile, err := os.OpenFile(path, os.O_RDONLY, 0) // #nosec
			if err != nil {
				return fmt.Errorf("error opening codeFile: %w", err)
			}

			defer func() {
				err = codeFile.Close()
				if err != nil {
					ui.PrintMessage(fmt.Sprintf("error closing codeFile: %s\n", err), ui.MessageTypeError)
				}
			}()

			file.Data, err = os.ReadFile(path) // #nosec

			if err != nil {
				return fmt.Errorf("error reading codeFile: %w", err)
			}

			files = append(files, *file)

			// Construct the codeFile tree
			parts := strings.Split(path, string(os.PathSeparator))
			current := rootNode

			for _, part := range parts[:len(parts)-1] { // Exclude the last part which is the codeFile itself
				found := false

				for _, child := range current.Children {
					if child.Name == part && child.IsDir {
						current = child
						found = true

						break
					}
				}

				if !found {
					newNode := &FileNode{Name: part, IsDir: true, Children: []*FileNode{}}
					current.Children = append(current.Children, newNode)
					current = newNode
				}
			}

			current.Children = append(current.Children,
				&FileNode{Name: parts[len(parts)-1], IsDir: false, Children: []*FileNode{}})

			return nil
		})
		if err != nil {
			return nil, nil, fmt.Errorf("error walking the path: %w", err)
		}
	}

	// Sort the files for consistent output
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	return files, rootNode, nil
}

type languageCheckerCache struct {
	cache     map[string]string
	cacheHits int
}

func (l *languageCheckerCache) Get(ext string) (string, bool) {
	val, ok := l.cache[ext]
	return val, ok
}

func (l *languageCheckerCache) Set(ext, lang string) {
	l.cache[ext] = lang
}

func cachedLanguageChecker() func(string) (string, bool) {
	cache := &languageCheckerCache{cache: make(map[string]string), cacheHits: 0}

	return func(path string) (string, bool) {
		if ext, ok := cache.Get(filepath.Ext(path)); ok {
			cache.cacheHits++
			return ext, true
		}
		// .md should be interpreted as markdown, not lisp
		if filepath.Ext(path) == ".md" {
			cache.Set(filepath.Ext(path), "markdown")
			return "markdown", true
		}

		for _, lang := range languages {
			if slices.Contains(lang.Extensions, filepath.Ext(path)) {
				// cache the extension for faster lookup
				cache.Set(filepath.Ext(path), lang.AceMode)
				return lang.AceMode, true
			}

			if slices.Contains(lang.Filenames, filepath.Base(path)) {
				// cache the filename for faster lookup
				cache.Set(filepath.Base(path), lang.AceMode)
				return lang.AceMode, true
			}
		}

		return "", false
	}
}

func GenerateFileTree(node *FileNode, indent string, isLast bool) string {
	// Handle the case for the root node differently
	var tree strings.Builder

	if node.Name == "/" && node.IsDir {
		tree.WriteString(".\n")
	} else {
		// Choose the appropriate prefix
		prefix := "├── "
		if isLast {
			prefix = "└── "
		}

		// Print the name of the current node with the correct indentation
		tree.WriteString(indent + prefix + node.Name + "\n")

		if node.IsDir && !isLast {
			indent += "│   "
		} else {
			indent += "    "
		}
	}

	if node.IsDir {
		// Sort node's children for consistent output
		sort.Slice(node.Children, func(i, j int) bool {
			return node.Children[i].Name < node.Children[j].Name
		})
		// Recursively call `printFileTree` for each child
		for i, child := range node.Children {
			// Check if the child node is the last in the list
			isLastChild := i == len(node.Children)-1
			tree.WriteString(GenerateFileTree(child, indent, isLastChild))
		}
	}

	return tree.String()
}
