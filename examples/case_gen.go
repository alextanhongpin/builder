// Code generated by github.com/alextanhongpin/builder, DO NOT EDIT.

package main

import "fmt"

type AccountBuilder struct {
	account   Account
	fields    map[string]int
	fieldsSet uint64
}

func NewAccountBuilder() *AccountBuilder {
	fields := make(map[string]int)
	for i, field := range []string{"ID", "Typ", "Name"} {
		fields[field] = i
	}
	return &AccountBuilder{fields: fields}
}

// WithID sets ID.
func (b AccountBuilder) WithID(iD int64) AccountBuilder {
	b.account.ID = iD
	b.Set("ID")
	return b
}

// WithTyp sets Typ.
func (b AccountBuilder) WithTyp(typ string) AccountBuilder {
	b.account.Typ = typ
	b.Set("Typ")
	return b
}

// WithName sets Name.
func (b AccountBuilder) WithName(name string) AccountBuilder {
	b.account.Name = name
	b.Set("Name")
	return b
}

// Build returns Account.
func (b AccountBuilder) Build() Account {
	for field := range b.fields {
		if !b.IsSet(field) {
			panic(fmt.Errorf("builder: %q not set", field))
		}
	}
	return b.account
}

// Build returns Account.
func (b AccountBuilder) BuildPartial() Account {
	return b.account
}

func (b *AccountBuilder) Set(field string) bool {
	n, ok := b.fields[field]
	if !ok {
		return false
	}
	b.fieldsSet |= 1 << n
	return true

}

func (b AccountBuilder) IsSet(field string) bool {
	pos := b.fields[field]
	return (b.fieldsSet & (1 << pos)) == (1 << pos)
}

func (b *AccountBuilder) Register(field string) error {
	if _, ok := b.fields[field]; ok {
		return fmt.Errorf("field %q already registered", field)
	}
	b.fields[field] = len(b.fields)
	return nil
}
