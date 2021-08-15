package gen

import (
	"os"
	"path/filepath"
	"strings"
)

// SkipCurrentPackagePath returns "" if the package is referring to an entity
// that belongs to itself, e.g. package foo calling foo.Foo.
func SkipCurrentPackagePath(pkgPath string, field StructField) string {
	if pkgPath == field.FieldPkgPath {
		return ""
	}
	return field.FieldPkgPath
}

// packagePath returns the github package path from any given path,
// e.g. path/to/github.com/your-repo/your-pkg returns github.com/your-repo/your-pkg
// If your package is not hosted on github, you may need to override $PKG to
// set the prefix of your package.
func packagePath(path string) string {
	if ext := filepath.Ext(path); ext != "" {
		base := filepath.Base(path)
		path = path[:len(path)-len(base)]
	}
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	pkg := os.Getenv("PKG")
	if pkg == "" {
		pkg = "github.com"
	}
	idx := strings.Index(path, pkg)
	return path[idx:]
}

// packageName returns the base package name.
func packageName(path string) string {
	return filepath.Base(packagePath(path))
}
