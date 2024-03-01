package filetree

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type FileNode struct {
	Name     string
	IsDir    bool
	Children []*FileNode
}

func GatherFiles(re *regexp.Regexp, paths []string, ignorePatterns []*regexp.Regexp) (map[string][]byte, []string, *FileNode, error) {
	fileMap := make(map[string][]byte)
	var sortedPaths []string

	ignoreMatcher := func(path string) bool {
		for _, pattern := range ignorePatterns {
			if pattern.MatchString(path) {
				return true
			}
		}
		return false
	}

	rootNode := &FileNode{Name: "/", IsDir: true}

	for _, path := range paths {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if !re.MatchString(path) || info.IsDir() || ignoreMatcher(path) {
				return nil
			}

			var file []byte
			file, err = os.ReadFile(path) // #nosec

			if err != nil {
				return err
			}

			fileMap[path] = file

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
			return nil, nil, nil, err
		}
	}

	for path := range fileMap {
		sortedPaths = append(sortedPaths, path)
	}
	sort.Strings(sortedPaths)

	return fileMap, sortedPaths, rootNode, nil
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
