package loader

import (
	"go/types"
)

// StructField for the example below.
//type Foo struct {
//  Name sql.NullString `json:"name"
//}
type StructField struct {
	// Name of the struct field.
	Name string `example:"Name"`

	// Useful when the output directory doesn't match the existing ones.
	PkgPath string `example:"github.com/alextanhongpin/go-codegen/test"`
	PkgName string `example:"test"`

	Exported bool `example:"true"`

	// Stores the original position of the field in the struct.
	Ordinal int

	Tag  string `example:"get:'Renamed'"` // To ignore getter.
	Type types.Type
}

func NewStructFields(structType *types.Struct) (map[string]StructField, error) {
	fields := make(map[string]StructField)
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		tag := structType.Tag(i)

		if skipField(tag) {
			continue
		}

		fields[field.Name()] = StructField{
			Name:     field.Name(),
			PkgPath:  field.Pkg().Path(),
			Exported: field.Exported(),
			Type:     field.Type(),
			Tag:      tag,
			Ordinal:  i,
		}
	}

	return fields, nil
}
