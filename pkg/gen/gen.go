package gen

import (
	"flag"
	"fmt"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alextanhongpin/pkg/stringcase"
	"golang.org/x/tools/go/packages"
)

// StructField for the example below.
//type Foo struct {
//  Name sql.NullString `json:"name"
//}
type StructField struct {
	Name string `example:"Name"`
	// Useful when the output directory doesn't match the existing ones.
	PkgPath   string `example:"github.com/alextanhongpin/go-codegen/test"`
	PkgName   string `example:"test"`
	Exported  bool   `example:"true"`
	FieldType string `example:"NullString"`

	// When true, the FieldPkgPath is not empty.
	NamedField   bool   `example:"true"`
	FieldPkgPath string `example:"database/sql"`
	Tag          string `example:"build:'-'"` // To ignore builder.
	IsPointer    bool
}

type Option struct {
	In         string
	Out        string
	PkgName    string
	PkgPath    string
	StructName string
	Fields     []StructField
}

type Generator func(opt Option) error

func New(fn Generator) {
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
		pkgPath := pkg.PkgPath              // Specify the config packages.NeedName to get this value.
		pkgName := pkg.Name                 // main

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
		if err := fn(Option{
			PkgName:    pkgName,
			PkgPath:    pkgPath,
			Out:        out,
			In:         in,
			StructName: structName,
			Fields:     fields,
		}); err != nil {
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
