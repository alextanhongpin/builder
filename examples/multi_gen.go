// Code generated by github.com/alextanhongpin/builder, DO NOT EDIT.

package main

import (
	"fmt"
	uuid "github.com/google/uuid"
)

type HelloBuilder struct {
	hello     Hello
	fields    []string
	fieldsSet uint64
}

func NewHelloBuilder(additionalFields ...string) *HelloBuilder {
	for _, field := range additionalFields {
		if field == "" {
			panic("builder: empty string in constructor")
		}
	}
	exists := make(map[string]bool)
	fields := append([]string{"id"}, additionalFields...)
	for _, field := range fields {
		if exists[field] {
			panic(fmt.Errorf("builder: duplicate field %q", field))
		}
		exists[field] = true
	}
	return &HelloBuilder{fields: fields}
}

// WithID sets id.
func (b HelloBuilder) WithID(id uuid.UUID) HelloBuilder {
	b.mustSet("id")
	b.hello.id = id
	return b
}

// Build returns Hello.
func (b HelloBuilder) Build() Hello {
	for i, field := range b.fields {
		if !b.isSet(i) {
			panic(fmt.Errorf("builder: %q not set", field))
		}
	}
	return b.hello
}

// Build returns Hello.
func (b HelloBuilder) BuildPartial() Hello {
	return b.hello
}

func (b *HelloBuilder) mustSet(field string) {
	i := b.indexOf(field)
	if b.isSet(i) {
		panic(fmt.Errorf("builder: set %q twice", field))
	}
	b.fieldsSet |= 1 << i
}

func (b HelloBuilder) isSet(pos int) bool {
	return (b.fieldsSet & (1 << pos)) == (1 << pos)
}

func (b HelloBuilder) indexOf(field string) int {
	for i, f := range b.fields {
		if f == field {
			return i
		}
	}
	panic(fmt.Errorf("builder: field: %q not found", field))
}

type WorldBuilder struct {
	world     World
	fields    []string
	fieldsSet uint64
}

func NewWorldBuilder(additionalFields ...string) *WorldBuilder {
	for _, field := range additionalFields {
		if field == "" {
			panic("builder: empty string in constructor")
		}
	}
	exists := make(map[string]bool)
	fields := append([]string{"id"}, additionalFields...)
	for _, field := range fields {
		if exists[field] {
			panic(fmt.Errorf("builder: duplicate field %q", field))
		}
		exists[field] = true
	}
	return &WorldBuilder{fields: fields}
}

// WithID sets id.
func (b WorldBuilder) WithID(id uuid.UUID) WorldBuilder {
	b.mustSet("id")
	b.world.id = id
	return b
}

// Build returns World.
func (b WorldBuilder) Build() World {
	for i, field := range b.fields {
		if !b.isSet(i) {
			panic(fmt.Errorf("builder: %q not set", field))
		}
	}
	return b.world
}

// Build returns World.
func (b WorldBuilder) BuildPartial() World {
	return b.world
}

func (b *WorldBuilder) mustSet(field string) {
	i := b.indexOf(field)
	if b.isSet(i) {
		panic(fmt.Errorf("builder: set %q twice", field))
	}
	b.fieldsSet |= 1 << i
}

func (b WorldBuilder) isSet(pos int) bool {
	return (b.fieldsSet & (1 << pos)) == (1 << pos)
}

func (b WorldBuilder) indexOf(field string) int {
	for i, f := range b.fields {
		if f == field {
			return i
		}
	}
	panic(fmt.Errorf("builder: field: %q not found", field))
}
