package summarizer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

var DefaultIgnorePatterns = []string{
	".git/",
	"summary.md",
	".gitignore",
}

type fileInfo struct {
	RelativePath string
	Content      []byte
	IsDir        bool
}

func SummarizeProject(rootDir string, outputFileName string, additionalIgnores []string) error {
	ignorePatterns := make([]string, 0, len(DefaultIgnorePatterns)+1+len(additionalIgnores))
	ignorePatterns = append(ignorePatterns, DefaultIgnorePatterns...)
	ignorePatterns = append(ignorePatterns, outputFileName)

	if len(additionalIgnores) > 0 {
		ignorePatterns = append(ignorePatterns, additionalIgnores...)
	}

	ignoreMatcher := gitignore.CompileIgnoreLines(ignorePatterns...)

	gitignorePath := filepath.Join(rootDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		if _, err = gitignore.CompileIgnoreFile(gitignorePath); err == nil {
			ignoreMatcher = gitignore.CompileIgnoreLines(func() []string {
				allPatterns := make([]string, len(ignorePatterns))
				copy(allPatterns, ignorePatterns)

				content, errFile := os.ReadFile(gitignorePath)
				if errFile == nil {
					lines := strings.Split(string(content), "\n")
					for _, line := range lines {
						line = strings.TrimSpace(line)
						if line != "" && !strings.HasPrefix(line, "#") {
							allPatterns = append(allPatterns, line)
						}
					}
				}
				return allPatterns
			}()...)
		}
	}

	includedFiles := make([]fileInfo, 0)
	filesForStructure := make([]string, 0)

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			_ = fmt.Errorf("Error accessing %q: %v\n", path, err)
			return err
		}
		relativePath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return fmt.Errorf("could not get relative path for %s: %w", path, err)
		}
		if relativePath == "." {
			return nil
		}

		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}

		checkPath := filepath.ToSlash(relativePath)
		if ignoreMatcher.MatchesPath(checkPath) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		filesForStructure = append(filesForStructure, relativePath)

		if !d.IsDir() {
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				return nil
			}
			includedFiles = append(includedFiles, fileInfo{
				RelativePath: relativePath,
				Content:      content,
				IsDir:        false,
			})
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error during directory traversal %s: %w", rootDir, err)
	}

	sort.Strings(filesForStructure)
	sort.Slice(includedFiles, func(i, j int) bool {
		return includedFiles[i].RelativePath < includedFiles[j].RelativePath
	})

	var summaryContent strings.Builder
	summaryContent.WriteString("# Project Structure\n\n```\n")
	treeString, err := generateProjectTree(filesForStructure)
	if err != nil {
		return fmt.Errorf("error generating project tree: %w", err)
	}
	summaryContent.WriteString(treeString)
	summaryContent.WriteString("```\n\n\n\n# Project Files\n\n")

	for _, file := range includedFiles {
		displayPath := filepath.ToSlash(file.RelativePath)
		summaryContent.WriteString(fmt.Sprintf("◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦\nfilepath: ./%s\n◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦\n", displayPath))
		summaryContent.WriteString(string(file.Content))
		if len(file.Content) > 0 && file.Content[len(file.Content)-1] != '\n' {
			summaryContent.WriteString("\n")
		}
		summaryContent.WriteString("◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦◦\n\n\n\n")
	}

	outputPath := filepath.Join(rootDir, outputFileName)
	err = os.WriteFile(outputPath, []byte(summaryContent.String()), 0o644)
	if err != nil {
		return fmt.Errorf("could not write the file %s: %w", outputFileName, err)
	}

	return nil
}

func generateProjectTree(paths []string) (string, error) {
	var tree strings.Builder
	lastTopLevelIndex := -1
	for i := len(paths) - 1; i >= 0; i-- {
		if !strings.Contains(paths[i], string(filepath.Separator)) {
			lastTopLevelIndex = i
			break
		}
	}

	for i, path := range paths {
		parts := strings.Split(path, string(filepath.Separator))
		level := len(parts) - 1
		var prefix strings.Builder
		for j := 0; j < level; j++ {
			parentPath := filepath.Join(parts[:j+1]...)
			isParentLast := isLastSibling(parentPath, paths)
			if isParentLast {
				prefix.WriteString("   ")
			} else {
				prefix.WriteString("│  ")
			}
		}
		isLast := isLastSibling(path, paths)
		if level == 0 {
			if i == lastTopLevelIndex {
				prefix.WriteString("└─ ")
			} else {
				prefix.WriteString("├─ ")
			}
		} else {
			if isLast {
				prefix.WriteString("└─ ")
			} else {
				prefix.WriteString("├─ ")
			}
		}
		tree.WriteString(prefix.String())
		tree.WriteString(parts[level])
		tree.WriteString("\n")
	}
	return tree.String(), nil
}

func isLastSibling(currentPath string, allPaths []string) bool {
	parentDir := filepath.Dir(currentPath)
	var effectiveParentDir string
	if parentDir == "." {
		effectiveParentDir = ""
	} else {
		effectiveParentDir = parentDir
	}

	lastSiblingInDir := ""
	for i := len(allPaths) - 1; i >= 0; i-- {
		p := allPaths[i]
		pDir := filepath.Dir(p)
		var effectivePDir string
		if pDir == "." {
			effectivePDir = ""
		} else {
			effectivePDir = pDir
		}

		if effectivePDir == effectiveParentDir {
			lastSiblingInDir = p
			break
		}
	}
	return currentPath == lastSiblingInDir
}
