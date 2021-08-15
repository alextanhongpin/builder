package gen

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// fullPath returns the full path to the package, relative to the caller.
func fullPath(rel string) string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get package directory: %v", err)
	}
	path = filepath.Join(path, rel)
	return path
}

// addSuffixToFileName adds a suffix to the filename, before the extension, to allow main.go -> main_gen.go
func addSuffixToFileName(path, suffix string) string {
	ext := filepath.Ext(path)
	path = path[:len(path)-len(ext)]
	return path + suffix + ext
}

// safeAddSuffixToFileName only adds the suffix if the user generated name does not already contains the suffix.
func safeAddSuffixToFileName(path, suffix string) string {
	if strings.Contains(path, suffix) {
		return path
	}
	return addSuffixToFileName(path, suffix)
}

func isFile(path string) bool {
	return filepath.Ext(path) != ""
}

func safeAddFileName(path, fileName string) string {
	if isFile(path) {
		return path
	}
	return filepath.Join(path, fileName)
}
