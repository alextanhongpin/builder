package main

import (
	"flag"
	"fmt"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alextanhongpin/pkg/stringcase"
	"github.com/dave/jennifer/jen"
	. "github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

// Generator is the name of this tool.
const Generator = "builder"

// StructField for the example below.
//type Foo struct {
//  Name sql.NullString `json:"name"
//}
type StructField struct {
	Name string `example:"Name"`
	// Useful when the output directory doesn't match the existing ones.
	PkgPath   string `example:"github.com/alextanhongpin/go-codegen/test"`
	Exported  bool   `example:"true"`
	FieldType string `example:"NullString"`

	// When true, the FieldPkgPath is not empty.
	NamedField   bool   `example:"true"`
	FieldPkgPath string `example:"database/sql"`
	Tag          string `example:""`
	IsPointer    bool
}

func main() {
	structPtr := flag.String("type", "", "the target struct name")
	inPtr := flag.String("in", os.Getenv("GOFILE"), "the input file, defaults to the file with the go:generate comment")
	outPtr := flag.String("out", "", "the output directory")
	flag.Parse()

	in := fullPath(*inPtr)

	// Allows -type=Foo,Bar
	structNames := strings.Split(*structPtr, ",")
	for _, structName := range structNames {
		var out string
		if o := *outPtr; o == "" {
			// Foo becomes foo.go
			fileName := stringcase.SnakeCase(structName) + ".go"

			// foo.go becomes foo_gen.go
			genFileName := safeAddSuffixToFileName(fileName, "_gen")

			// path/to/main.go becomes path/to/foo_gen.go
			out = safeAddFileName(filepath.Dir(in), genFileName)
		} else {
			out = fullPath(o)
		}

		pkg := loadPackage(packagePath(in)) // github.com/your-github-username/your-pkg.
		pkgName := pkg.Name                 // Specify the config packages.NeedName to get this value.

		obj := pkg.Types.Scope().Lookup(structName)
		if obj == nil {
			log.Fatalf("struct %s not found", structName)
		}

		// Check if it is a declared typed.
		if _, ok := obj.(*types.TypeName); !ok {
			log.Fatalf("%v is not a named type", obj)
		}

		// Check if the type is a struct.
		structType, ok := obj.Type().Underlying().(*types.Struct)
		if !ok {
			log.Fatalf("%v is not a struct", obj)
		}

		fields := extractFields(structType)
		if err := generateStructFromFields(pkgName, out, structName, fields); err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("wrote file to %s\n", out)
	}
}

func loadPackage(path string) *packages.Package {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedImports,
	}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		log.Fatalf("failed to load package: %v", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}
	return pkgs[0]
}

func extractFields(structType *types.Struct) []StructField {
	fields := make([]StructField, structType.NumFields())
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		tag := structType.Tag(i)

		var (
			name         = field.Name()
			pkgPath      = field.Pkg().Path()
			exported     = field.Exported()
			namedField   = false
			fieldPkgPath = ""
			fieldType    = ""
			isPointer    = false
		)

		typ := field.Type()
		if ptr, ok := field.Type().(*types.Pointer); ok {
			isPointer = true
			typ = ptr.Elem()
		}

		switch t := typ.(type) {
		case *types.Named:
			obj := t.Obj()
			fieldPkgPath = obj.Pkg().Path()
			fieldType = obj.Name()
			namedField = true
		//case *types.Pointer:
		//// If it is a pointer, get the element.
		//typ := t.Elem()
		//isPointer  =true
		//switch t := typ.(type) {
		//case *types.Named:
		//obj := t.Obj()
		//fieldPkgPath = obj.Pkg().Path()
		//fieldType = obj.Name()
		//namedField = true
		//default:
		//fieldType = t.String()
		//namedField = false
		//}

		default:
			fieldType = t.String()
			namedField = false
		}
		fields[i] = StructField{
			Name:         name,
			PkgPath:      pkgPath,
			Exported:     exported,
			FieldType:    fieldType,
			NamedField:   namedField,
			FieldPkgPath: fieldPkgPath,
			Tag:          tag,
			IsPointer:    isPointer,
		}
	}
	return fields
}

func genBuilderName(structName string) string {
	builderName := fmt.Sprintf("%sBuilder", structName)
	return builderName
}

func generateStructFromFields(pkgName, out, structName string, fields []StructField) error {
	f := NewFile(pkgName) // e.g. main
	f.PackageComment(fmt.Sprintf("Code generated by %s, DO NOT EDIT.", Generator))

	generateBuilder(f, structName)
	generateBuilderConstructor(f, structName)

	for _, field := range fields {
		if field.IsPointer {
			generateWitherPointer(f, structName, field)
		} else {
			generateWither(f, structName, field)
		}
	}

	generateBuildFunc(f, structName)

	return f.Save(out) // e.g. main_gen.go
}

func generateBuilder(f *jen.File, structName string) {
	// Output:
	//type FooBuilder struct {
	//  foo Foo
	//}

	f.Type().Id(genBuilderName(structName)).Struct(
		Id(lowerFirst(structName)).Id(structName),
	).Line()
}

func generateBuilderConstructor(f *jen.File, structName string) {
	// Output:
	//func NewFooBuilder() *FooBuilder {
	//  return &FooBuilder{}
	//}

	builderName := genBuilderName(structName)
	f.Func().Id(fmt.Sprintf("New%s", builderName)).Params().Op("*").Id(builderName).Block(
		Return(Op("&").Id(builderName).Values()),
	).Line()
}

func generateWither(f *jen.File, structName string, field StructField) {
	// Output:
	// WithName sets name.
	// func (b FooBuilder) WithName(name string) FooBuilder {
	// 	b.foo.name = name
	//  return b
	// }

	builderName := genBuilderName(structName)
	shortName := "b"
	funcName := fmt.Sprintf("With%s", upperCommonInitialism(field.Name))
	f.Comment(fmt.Sprintf("%s sets %s.", funcName, field.Name))
	f.Func().Params(
		Id(shortName).Id(builderName), // (b *FooBuilder)
	).Id(funcName). // WithName
			Params(
			// name string
			Id(lowerFirst(field.Name)).Qual(field.FieldPkgPath, field.FieldType),
		).
		Id(builderName). // Return type: FooBuilder
		Block(
			Id(shortName).Dot(lowerFirst(structName)).Dot(field.Name).Op("=").Id(lowerFirst(field.Name)),
			Return(Id(shortName)),
		).Line()
}

func generateWitherPointer(f *jen.File, structName string, field StructField) {
	// Output:
	// WithName sets name.
	// func (b FooBuilder) WithName(name string, valid bool) FooBuilder {
	//  if valid {
	//    b.foo.name = name
	//  }
	//  return b
	// }

	builderName := genBuilderName(structName)
	shortName := "b"
	// Avoid conflict in naming if the variable is already named "valid".
	validVar := "valid"
	if field.Name == validVar {
		validVar += "1"
	}
	funcName := fmt.Sprintf("With%s", upperCommonInitialism(field.Name))
	f.Comment(fmt.Sprintf("%s sets %s.", funcName, field.Name))
	f.Func().Params(
		Id(shortName).Id(builderName), // (b *FooBuilder)
	).Id(funcName). // WithName
			Params(
			// name string
			Id(lowerFirst(field.Name)).Qual(field.FieldPkgPath, field.FieldType),
			Id(validVar).Bool(),
		).
		Id(builderName). // Return type: FooBuilder
		Block(
			If(Id(validVar)).Block(
				Id(shortName).Dot(lowerFirst(structName)).Dot(field.Name).Op("=").Op("&").Id(lowerFirst(field.Name)),
			),
			Return(Id(shortName)),
		).Line()
}

func generateBuildFunc(f *jen.File, structName string) {
	// Output:
	// Build returns Foo.
	// func (b FooBuilder) Build() Foo {
	//  return b.foo
	// }

	builderName := genBuilderName(structName)
	shortName := "b"
	f.Comment(fmt.Sprintf("Build returns %s.", structName))
	f.Func().Params(
		Id(shortName).Id(builderName), // (b *FooBuilder)
	).Id("Build").Params(). // Build()
				Id(structName). // Return type: Foo
				Block(
			Return(Id(shortName).Dot(lowerFirst(structName))),
		).Line()
}

func upperFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func lowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// fullPath returns the full path to the package, relative to the caller.
func fullPath(rel string) string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get package directory: %v", err)
	}
	path = filepath.Join(path, rel)
	return path
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

// https://github.com/golang/lint/blob/83fdc39ff7b56453e3793356bcff3070b9b96445/lint.go#L770-L809
// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

func upperCommonInitialism(s string) string {
	isCommonInitialisms := commonInitialisms[strings.ToUpper(s)]
	if isCommonInitialisms {
		return strings.ToUpper(s)
	}
	return upperFirst(s)
}
