// Code generated by github.com/alextanhongpin/builder, DO NOT EDIT.
package main

import "fmt"

type UserBuilder struct {
	user      User
	fields    []string
	fieldsSet uint64
}

func NewUserBuilder(additionalFields ...string) *UserBuilder {
	for _, field := range additionalFields {
		if field == "" {
			panic("builder: empty string in constructor")
		}
	}
	exists := make(map[string]bool)
	fields := append([]string{"name", "age", "married", "hobbies"}, additionalFields...)
	for _, field := range fields {
		if exists[field] {
			panic(fmt.Errorf("builder: duplicate field %q", field))
		}
		exists[field] = true
	}
	return &UserBuilder{fields: fields}
}

// WithName sets name.
func (b UserBuilder) WithName(name string) UserBuilder {
	b.mustSet("name")
	b.user.name = name
	return b
}

// WithAge sets age.
func (b UserBuilder) WithAge(age int64) UserBuilder {
	b.mustSet("age")
	b.user.age = age
	return b
}

// WithMarried sets married.
func (b UserBuilder) WithMarried(married bool) UserBuilder {
	b.mustSet("married")
	b.user.married = married
	return b
}

// WithHobbies sets hobbies.
func (b UserBuilder) WithHobbies(hobbies []string) UserBuilder {
	b.mustSet("hobbies")
	b.user.hobbies = hobbies
	return b
}

// Build returns User.
func (b UserBuilder) Build() User {
	for i, field := range b.fields {
		if !b.isSet(i) {
			panic(fmt.Errorf("builder: %q not set", field))
		}
	}
	return b.user
}

// Build returns User.
func (b UserBuilder) BuildPartial() User {
	return b.user
}

func (b *UserBuilder) mustSet(field string) {
	i := b.indexOf(field)
	if b.isSet(i) {
		panic(fmt.Errorf("builder: set %q twice", field))
	}
	b.fieldsSet |= 1 << i
}

func (b UserBuilder) isSet(pos int) bool {
	return (b.fieldsSet & (1 << pos)) == (1 << pos)
}

func (b UserBuilder) indexOf(field string) int {
	for i, f := range b.fields {
		if f == field {
			return i
		}
	}
	panic(fmt.Errorf("builder: field: %q not found", field))
}
