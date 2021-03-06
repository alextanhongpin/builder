package main

import (
	"fmt"
	"go/types"
	"sort"

	"github.com/alextanhongpin/builder/loader"
	"github.com/dave/jennifer/jen"
	. "github.com/dave/jennifer/jen"
)

// Generator is the name of this tool.
const Generator = "github.com/alextanhongpin/builder"

func main() {
	if err := loader.New(generateStructFromFields); err != nil {
		panic(err)
	}
}

func generateBuilderName(structName string) string {
	builderName := fmt.Sprintf("%sBuilder", structName)
	return builderName
}

func generateStructFromFields(opt loader.Option) error {
	var (
		pkgName = opt.PkgName
		pkgPath = opt.PkgPath
		out     = opt.Out
		structs = opt.Items
	)

	// If there are only 1 struct, then generate the file name based on the
	// struct name. Otherwise, generate the file name based on the input file
	// name with _gen suffix.
	if len(structs) == 1 {
		out = structs[0].Path
	}

	f := NewFilePathName(pkgPath, pkgName) // e.g. main
	f.HeaderComment(fmt.Sprintf("Code generated by %s, DO NOT EDIT.", Generator))

	for _, str := range structs {
		structType := str.Type
		structName := str.Name

		fields, err := generateSortedStructFields(structType)
		if err != nil {
			return err
		}

		generateBuilder(f, structName)
		generateBuilderConstructor(f, structName, fields)

		for i, field := range fields {
			if isPointer(field.Type) {
				generateWitherPointer(f, pkgPath, structName, field, i)
			} else {
				generateWither(f, pkgPath, structName, field, i)
			}
		}

		generateBuildFunc(f, structName)
		generateBuildPartialFunc(f, structName)
		generateMustSetFunc(f, structName)
		generateIsSetFunc(f, structName)
		generateRegisterFunc(f, structName)
	}
	if err := f.Save(out); err != nil { // e.g. main_gen.go
		return err
	}
	fmt.Printf("success: generated %s\n", out)
	return nil
}

func generateBuilder(f *jen.File, structName string) {
	// Output:
	//
	// type FooBuilder struct {
	//   foo Foo
	//   fields map[string]int
	//   fieldsSet uint64
	// }

	f.Type().Id(generateBuilderName(structName)).Struct(
		Id(loader.LowerFirst(structName)).Id(structName),
		Id("fields").Map(String()).Int(),
		Id("fieldsSet").Uint64(),
	).Line()
}

func generateBuilderConstructor(f *jen.File, structName string, fields []loader.StructField) {
	/*
		 Output:

		 func NewFooBuilder() *FooBuilder {
			 fields := make(map[string]int)
				 for i, field := range []string{"name", "age"}{
					 fields[field] = i
				 }
			 }
			 return &FooBuilder{
				 fields: fields,
			 }
		 }
	*/

	values := make([]Code, len(fields))
	for i, field := range fields {
		values[i] = Lit(field.Name)
	}

	builderName := generateBuilderName(structName)
	f.Func().Id(fmt.Sprintf("New%s", builderName)).Params().Op("*").Id(builderName).Block(
		Id("fields").Op(":=").Make(Map(String()).Int()),
		For(List(Id("i"), Id("field")).Op(":=").Range().Index().String().Values(values...)).Block(
			Id("fields").Index(Id("field")).Op("=").Id("i"),
		),
		Return(Op("&").Id(builderName).Values(Dict{
			Id("fields"): Id("fields"),
		})),
	).Line()
}

func generateWither(f *jen.File, pkgPath, structName string, field loader.StructField, pos int) {
	// Output:
	// WithName sets name.
	// func (b FooBuilder) WithName(name string) FooBuilder {
	// 	 b.foo.name = name
	//   b.Set("name")
	//   return b
	// }

	builderName := generateBuilderName(structName)
	shortName := "b"
	funcName := fmt.Sprintf("With%s", loader.UpperCommonInitialism(field.Name))

	f.Comment(fmt.Sprintf("%s sets %s.", funcName, field.Name))
	f.Func().Params(
		Id(shortName).Id(builderName), // (b *FooBuilder)
	).Id(funcName). // WithName
			Params(Id(loader.LowerFirst(field.Name)).Add(generateType(field.Type))). // name string
			Id(builderName).                                                         // Return type: FooBuilder
			Block(
			Id(shortName).Dot(loader.LowerFirst(structName)).Dot(field.Name).Op("=").Id(loader.LowerFirst(field.Name)),
			Id("b").Dot("Set").Call(Lit(field.Name)),
			Return(Id(shortName)),
		).Line()
}

func generateWitherPointer(f *jen.File, pkgPath, structName string, field loader.StructField, pos int) {
	// Output:
	// WithName sets name.
	// func (b FooBuilder) WithName(name string, valid bool) FooBuilder {
	//   if valid {
	//     b.foo.name = name
	//   }
	//   b.Set("name")
	//   return b
	// }

	builderName := generateBuilderName(structName)
	shortName := "b"
	// Avoid conflict in naming if the variable is already named "valid".
	validVar := "valid"
	if field.Name == validVar {
		validVar += "1"
	}

	funcName := fmt.Sprintf("With%s", loader.UpperCommonInitialism(field.Name))
	f.Comment(fmt.Sprintf("%s sets %s.", funcName, field.Name))
	f.Func().
		Params(Id(shortName).Id(builderName)). // (b *FooBuilder)
		Id(funcName).                          // WithName
		Params(
			Id(loader.LowerFirst(field.Name)).Add(generateType(field.Type)), // name string
			Id(validVar).Bool(), // valid bool
		).
		Id(builderName). // Return type: FooBuilder
		Block(
			If(Id(validVar)).Block(
				Id(shortName).Dot(loader.LowerFirst(structName)).Dot(field.Name).Op("=").Do(func(s *Statement) {
					if isPointer(field.Type) {
						s.Add(Op("&"))
					}
				}).Id(loader.LowerFirst(field.Name)),
			),
			Id("b").Dot("Set").Call(Lit(field.Name)),
			Return(Id(shortName)),
		).Line()
}

func generateBuildFunc(f *jen.File, structName string) {
	// Output:
	//
	// // Build returns Foo.
	// func (b FooBuilder) Build() Foo {
	//   for field := range b.fields {
	//     if !b.isSet(i) {
	//       panic(fmt.Sprintf("builder: %q not set", field)
	//     }
	//   }
	//   return b.foo
	// }

	builderName := generateBuilderName(structName)
	shortName := "b"
	f.Comment(fmt.Sprintf("Build returns %s.", structName))
	f.Func().Params(
		Id(shortName).Id(builderName), // (b *FooBuilder)
	).Id("Build").Params(). // Build()
				Id(structName). // Return type: Foo
				Block(
			For(
				Id("field").Op(":=").Range().Id("b").Dot("fields"),
			).Block(
				If(Op("!").Id("b").Dot("IsSet").Call(Id("field"))).Block(
					Panic(Qual("fmt", "Errorf").Call(Lit("builder: %q not set"), Id("field"))),
				),
			),
			Return(Id(shortName).Dot(loader.LowerFirst(structName))),
		).Line()
}

func generateBuildPartialFunc(f *jen.File, structName string) {
	// Output:
	//
	// // Build returns Foo.
	// func (b FooBuilder) BuildPartial() Foo {
	//   return b.foo
	// }

	builderName := generateBuilderName(structName)
	shortName := "b"
	f.Comment(fmt.Sprintf("Build returns %s.", structName))
	f.Func().
		Params(Id(shortName).Id(builderName)). // (b FooBuilder)
		Id("BuildPartial").Params().           // Build()
		Id(structName).                        // Return type: Foo
		Block(
			Return(Id(shortName).Dot(loader.LowerFirst(structName))),
		).Line()
}

func generateMustSetFunc(f *jen.File, structName string) {
	/*
		 Output:

		 func (b *FooBuilder) Set(field string) bool {
			 n, ok := b.fields[field]
			 if !ok {
				 return false
			 }
			 b.fieldsSet |= 1 << n
			 return true
		 }
	*/

	builderName := generateBuilderName(structName)
	shortName := "b"
	f.Func().
		Params(Id(shortName).Op("*").Id(builderName)). // (b FooBuilder)
		Id("Set").Params(Id("field").String()).        // mustSet(field string)
		Id("bool").                                    // Return type: bool
		Block(
			List(Id("n"), Id("ok")).Op(":=").Id("b").Dot("fields").Index(Id("field")),
			If(Op("!").Id("ok").Block(
				Return(Lit(false)),
			),
				Id("b").Dot("fieldsSet").Op("|=").Lit(1).Op("<<").Id("n"),
				Return(Lit(true)),
			).Line(),
		).Line()
}

func generateIsSetFunc(f *jen.File, structName string) {
	// Output:
	//
	// func (b SimpleBuilder) IsSet(field string) bool {
	//   pos := b.fields[field]
	//   return (b.fieldsSet & (1 << pos)) == (1 << pos)
	// }

	builderName := generateBuilderName(structName)
	shortName := "b"

	f.Func().Params(
		Id(shortName).Id(builderName), // (b FooBuilder)
	).Id("IsSet").Params(Id("field").String()). // IsSet(field string)
							Bool(). // bool
							Block(
			Id("pos").Op(":=").Id("b").Dot("fields").Index(Id("field")),
			Return(Parens(Id("b").Dot("fieldsSet").Op("&").Parens(Lit(1).Op("<<").Id("pos"))).Op("==").Parens(Lit(1).Op("<<").Id("pos"))),
		).Line()
}

func generateRegisterFunc(f *jen.File, structName string) {
	/*
		Output:

		func (b *FooBuilder) Register(field string) error {
			if _, ok := b.fields[field]; ok {
				return fmt.Errorf("field %q already registered", field)
			}
			b.fields[field] = len(b.fields)
			return nil
		}
	*/

	builderName := generateBuilderName(structName)
	shortName := "b"
	f.Func().Params(
		Id(shortName).Op("*").Id(builderName), // (b *FooBuilder)
	).Id("Register").Params(Id("field").String()). // Register(field string)
							Id("error"). // error
							Block(
			If(
				List(Id("_"), Id("ok")).Op(":=").Id("b").Dot("fields").Index(Id("field")),
				Id("ok"),
			).Block(
				Return(Qual("fmt", "Errorf").Call(Lit("field %q already registered"), Id("field"))),
			),
			Id("b").Dot("fields").Index(Id("field")).Op("=").Len(Id("b").Dot("fields")),
			Return(Nil()),
		).Line()
}

func generateSortedStructFields(T types.Type) ([]loader.StructField, error) {
	fields, err := loader.NewStructFields((T).(*types.Struct))
	if err != nil {
		return nil, err
	}

	var result []loader.StructField
	for _, field := range fields {
		result = append(result, field)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Ordinal < result[j].Ordinal
	})

	return result, nil
}

type pointerVisitor struct {
	isPointer bool
	depth     int
}

func newPointerVisitor() *pointerVisitor {
	return &pointerVisitor{}
}

func (v *pointerVisitor) Visit(T types.Type) bool {
	switch T.(type) {
	case *types.Pointer:
		v.isPointer = v.depth == 0
		return false
	default:
		v.depth++
		return true
	}
}

func isPointer(T types.Type) bool {
	v := newPointerVisitor()
	_ = loader.Walk(v, T)
	return v.isPointer
}

type returnTypeVisitor struct {
	code         *Statement
	allowPointer bool
}

func newReturnTypeVisitor(allowPointer bool) *returnTypeVisitor {
	return &returnTypeVisitor{
		code:         Null(),
		allowPointer: allowPointer,
	}
}

func (v *returnTypeVisitor) Visit(T types.Type) bool {
	switch u := T.(type) {
	case *types.Pointer:
		if v.allowPointer {
			v.code = v.code.Op("*")
		}
		v.allowPointer = true
	case *types.Slice:
		v.code = v.code.Index()
		v.allowPointer = true
	case *types.Array:
		v.code = v.code.Index(Lit(u.Len()))
		v.allowPointer = true
	case *types.Map:
		iv := newReturnTypeVisitor(true)
		_ = loader.Walk(iv, u.Key())
		v.code = v.code.Map(iv.code)
		v.allowPointer = true
	case *types.Named:
		o := u.Obj()
		p := o.Pkg()
		v.code = v.code.Qual(p.Path(), o.Name())
		v.allowPointer = true
		return false
	default:
		v.code = v.code.Id(u.String())
		v.allowPointer = true
	}
	return true
}

func generateType(T types.Type) *Statement {
	v := newReturnTypeVisitor(false)
	_ = loader.Walk(v, T)
	return v.code
}
