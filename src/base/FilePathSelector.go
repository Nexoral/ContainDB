package base

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

// SelectFilePath provides a file selection prompt with path completion
func SelectFilePath(label string, defaultPath string, extension string) (string, error) {
	currentPath := defaultPath

	for {
		// If the path is a directory, append a separator
		pathInfo, err := os.Stat(currentPath)
		if err == nil && pathInfo.IsDir() && !strings.HasSuffix(currentPath, string(os.PathSeparator)) {
			currentPath += string(os.PathSeparator)
		}

		prompt := promptui.Prompt{
			Label:   label,
			Default: currentPath,
		}

		result, err := prompt.Run()
		if err != nil {
			return "", err
		}

		// Check if the result is a valid file that exists
		fileInfo, err := os.Stat(result)
		if err == nil {
			// If it's a directory, explore it
			if fileInfo.IsDir() {
				currentPath = result
				if !strings.HasSuffix(currentPath, string(os.PathSeparator)) {
					currentPath += string(os.PathSeparator)
				}

				// List directory contents
				files, err := os.ReadDir(currentPath)
				if err != nil {
					fmt.Printf("Error reading directory: %v\n", err)
					continue
				}

				fmt.Println("\nDirectory contents:")
				for _, file := range files {
					name := file.Name()
					if file.IsDir() {
						name += "/"
					}
					fmt.Println("  " + name)
				}

				continue
			} else {
				// It's a file and it exists, check if it has the right extension
				if extension == "" || strings.HasSuffix(result, extension) {
					return result, nil
				} else {
					fmt.Printf("File doesn't have the expected extension (%s). Continue anyway? (y/n): ", extension)
					var response string
					fmt.Scanln(&response)
					if strings.ToLower(response) == "y" {
						return result, nil
					}
					continue
				}
			}
		}

		// Path doesn't exist, check if it's a tab completion request
		dir := filepath.Dir(result)
		if dir != "." {
			// Check if directory exists
			if _, err := os.Stat(dir); err == nil {
				// List files in directory for possible completions
				files, err := os.ReadDir(dir)
				if err == nil {
					prefix := filepath.Base(result)
					matches := []string{}

					for _, file := range files {
						if strings.HasPrefix(file.Name(), prefix) {
							matches = append(matches, filepath.Join(dir, file.Name()))
						}
					}

					if len(matches) > 0 {
						fmt.Println("\nPossible completions:")
						for _, match := range matches {
							info, _ := os.Stat(match)
							if info.IsDir() {
								fmt.Println("  " + match + "/")
							} else {
								fmt.Println("  " + match)
							}
						}
						if len(matches) == 1 {
							currentPath = matches[0]
							continue
						}
					}
				}
			}
		}

		// Confirm creating new or using non-existent file
		fmt.Println("File doesn't exist. Use this path anyway? (y/n):")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) == "y" {
			return result, nil
		}

		// Start over with the current path
	}
}
