package filetree

import (
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

func GatherFiles(includeMatcher pm.PathMatcher, excludeMatcher pm.PathMatcher, pathScopes []string) ([]File, *FileNode, error) {
	var files []File

	cache := make(map[string]string)

	// TODO: Move this outside the file gathering function
	knownLanguage := func(path string) bool {
		if _, ok := cache[filepath.Ext(path)]; ok {
			return true
		}
		// .md should be interpreted as markdown, not lisp
		if filepath.Ext(path) == ".md" {
			cache[filepath.Ext(path)] = "markdown"
			return true
		}

		for _, lang := range languages {
			if slices.Contains(lang.Extensions, filepath.Ext(path)) {
				// cache the extension for faster lookup
				cache[filepath.Ext(path)] = lang.AceMode
				return true
			}
			if slices.Contains(lang.Filenames, filepath.Base(path)) {
				// cache the filename for faster lookup
				cache[filepath.Base(path)] = lang.AceMode
				return true
			}
		}

		ui.PrintMessage("skipping invalid source file: "+path+"\n", ui.MessageTypeWarning)
		return false
	}

	rootNode := &FileNode{Name: "/", IsDir: true}

	for _, path := range pathScopes {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			// start by skipping the .git directory
			if strings.HasPrefix(path, ".git/") {
				return nil
			}

			if !includeMatcher.Match(path) || info.IsDir() || excludeMatcher.Match(path) || !knownLanguage(path) {
				return nil
			}

			f := &File{
				Path: path,
				Type: cache[filepath.Ext(path)],
			}

			file, err := os.OpenFile(path, os.O_RDONLY, 0) // #nosec
			if err != nil {
				return err
			}

			defer file.Close()

			f.Data, err = os.ReadFile(path) // #nosec

			if err != nil {
				return err
			}

			files = append(files, *f)

			// Construct the file tree
			parts := strings.Split(path, string(os.PathSeparator))
			current := rootNode
			for _, part := range parts[:len(parts)-1] { // Exclude the last part which is the file itself
				found := false
				for _, child := range current.Children {
					if child.Name == part && child.IsDir {
						current = child
						found = true
						break
					}
				}
				if !found {
					newNode := &FileNode{Name: part, IsDir: true}
					current.Children = append(current.Children, newNode)
					current = newNode
				}
			}
			current.Children = append(current.Children, &FileNode{Name: parts[len(parts)-1], IsDir: false})

			return nil
		})

		if err != nil {
			return nil, nil, err
		}
	}

	// Sort the files for consistent output
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})
	return files, rootNode, nil
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
