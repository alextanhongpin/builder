package loader

import (
	"flag"
	"fmt"
	"go/types"
	"os"
	"strings"
)

type Visitor interface {
	Visit(T types.Type) bool
}

func Walk(visitor Visitor, T types.Type) bool {
	if next := visitor.Visit(T); !next {
		return next
	}

	switch u := T.Underlying().(type) {
	case *types.Named:
		return Walk(visitor, u.Underlying())
	case *types.Pointer:
		return Walk(visitor, u.Elem())
	case *types.Array:
		return Walk(visitor, u.Elem())
	case *types.Slice:
		return Walk(visitor, u.Elem())
	case *types.Map:
		return Walk(visitor, u.Elem())
	default:
		return types.IdenticalIgnoreTags(T, u)
	}
}

type Option struct {
	In         string
	Out        string
	PkgName    string
	PkgPath    string
	Prefix     string
	StructName string
	Type       types.Type
}

type Generator func(opt Option) error

func New(fn Generator) error {
	inp := flag.String("in", os.Getenv("GOFILE"), "the input file, defaults to the file with the go:generate comment")
	outp := flag.String("out", "", "the output directory")
	pkgp := flag.String("pkg", "github.com", "the package or package prefix path")
	prefixp := flag.String("prefix", "", "the generated field name prefix, e.g. Get")
	prunep := flag.Bool("prune", true, "whether to prune the old file before generating a new one")
	structp := flag.String("type", "", "the target struct name")
	flag.Parse()

	in := fullPath(*inp)
	pkg := loadPackage(packagePath(*pkgp, in)) // github.com/your-github-username/your-pkg.
	pkgPath := pkg.PkgPath                     // Specify the config packages.NeedName to get this value.
	pkgName := pkg.Name                        // main

	// Allows -type=Foo,Bar
	structNames := strings.Split(*structp, ",")
	for _, structName := range structNames {
		out := FileNameFromTypeName(*inp, *outp, structName)
		if *prunep {
			if err := os.Remove(out); err != nil && !os.IsNotExist(err) {
				fmt.Printf("failed to remove %s: %s\n", out, err)
			}
		}

		obj := pkg.Types.Scope().Lookup(structName)
		if obj == nil {
			return fmt.Errorf("struct %s not found", structName)
		}

		// Check if it is a declared typed.
		if _, ok := obj.(*types.TypeName); !ok {
			return fmt.Errorf("%v is not a named type", obj)
		}

		// Check if the type is a struct.
		_, ok := obj.Type().Underlying().(*types.Struct)
		if !ok {
			return fmt.Errorf("%v is not a struct", obj)
		}

		if err := fn(Option{
			PkgName:    pkgName,
			PkgPath:    pkgPath,
			Prefix:     *prefixp,
			Out:        out,
			In:         in,
			StructName: structName,
			Type:       obj.Type().Underlying(),
		}); err != nil {
			return err
		}
	}
	return nil
}
