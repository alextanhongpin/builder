// Code generated by builder, DO NOT EDIT.
package main

import (
	"database/sql"
	"fmt"
)

type FooBuilder struct {
	foo    Foo
	fields map[string]bool
}

func NewFooBuilder() *FooBuilder {
	return &FooBuilder{fields: map[string]bool{
		"age":                 false,
		"bar":                 false,
		"barByString":         false,
		"barPtrByString":      false,
		"barPtrs":             false,
		"bars":                false,
		"id":                  false,
		"name":                false,
		"realAge":             false,
		"sliceBarByString":    false,
		"sliceBarPtrByString": false,
		"stringByBar":         false,
		"stringByBarPtr":      false,
		"url":                 false,
		"valid":               false,
	}}
}

// WithID sets id.
func (b FooBuilder) WithID(id int64) FooBuilder {
	b.setOrPanic("id")
	b.foo.id = id
	return b
}

// WithName sets name.
func (b FooBuilder) WithName(name string) FooBuilder {
	b.setOrPanic("name")
	b.foo.name = name
	return b
}

// WithAge sets age.
func (b FooBuilder) WithAge(age sql.NullInt64) FooBuilder {
	b.setOrPanic("age")
	b.foo.age = age
	return b
}

// WithValid sets valid.
func (b FooBuilder) WithValid(valid bool, valid1 bool) FooBuilder {
	b.setOrPanic("valid")
	if valid1 {
		b.foo.valid = &valid
	}
	return b
}

// WithURL sets url.
func (b FooBuilder) WithURL(url string) FooBuilder {
	b.setOrPanic("url")
	b.foo.url = url
	return b
}

// WithRealAge sets realAge.
func (b FooBuilder) WithRealAge(realAge int64, valid bool) FooBuilder {
	b.setOrPanic("realAge")
	if valid {
		b.foo.realAge = &realAge
	}
	return b
}

// WithBar sets bar.
func (b FooBuilder) WithBar(bar Bar) FooBuilder {
	b.setOrPanic("bar")
	b.foo.bar = bar
	return b
}

// WithBars sets bars.
func (b FooBuilder) WithBars(bars []Bar) FooBuilder {
	b.setOrPanic("bars")
	b.foo.bars = bars
	return b
}

// WithBarPtrs sets barPtrs.
func (b FooBuilder) WithBarPtrs(barPtrs []*Bar, valid bool) FooBuilder {
	b.setOrPanic("barPtrs")
	if valid {
		b.foo.barPtrs = barPtrs
	}
	return b
}

// WithBarByString sets barByString.
func (b FooBuilder) WithBarByString(barByString map[string]Bar) FooBuilder {
	b.setOrPanic("barByString")
	b.foo.barByString = barByString
	return b
}

// WithStringByBar sets stringByBar.
func (b FooBuilder) WithStringByBar(stringByBar map[Bar]string) FooBuilder {
	b.setOrPanic("stringByBar")
	b.foo.stringByBar = stringByBar
	return b
}

// WithSliceBarByString sets sliceBarByString.
func (b FooBuilder) WithSliceBarByString(sliceBarByString map[string][]Bar) FooBuilder {
	b.setOrPanic("sliceBarByString")
	b.foo.sliceBarByString = sliceBarByString
	return b
}

// WithSliceBarPtrByString sets sliceBarPtrByString.
func (b FooBuilder) WithSliceBarPtrByString(sliceBarPtrByString map[string][]*Bar) FooBuilder {
	b.setOrPanic("sliceBarPtrByString")
	b.foo.sliceBarPtrByString = sliceBarPtrByString
	return b
}

// WithBarPtrByString sets barPtrByString.
func (b FooBuilder) WithBarPtrByString(barPtrByString map[string]*Bar) FooBuilder {
	b.setOrPanic("barPtrByString")
	b.foo.barPtrByString = barPtrByString
	return b
}

// WithStringByBarPtr sets stringByBarPtr.
func (b FooBuilder) WithStringByBarPtr(stringByBarPtr map[*Bar]string) FooBuilder {
	b.setOrPanic("stringByBarPtr")
	b.foo.stringByBarPtr = stringByBarPtr
	return b
}

// Build returns Foo.
func (b FooBuilder) Build() Foo {
	for field, isSet := range b.fields {
		if !isSet {
			panic(fmt.Sprintf("builder.BuildErr: %q not set", field))
		}
	}
	return b.foo
}

// Build returns Foo.
func (b FooBuilder) BuildPartial() Foo {
	return b.foo
}

// setOrPanic sets the fields only if it has not yet been set. It will panic when calling it twice.
func (b *FooBuilder) setOrPanic(field string) {
	c := b.cloneFields()
	if c[field] {
		panic(fmt.Sprintf("builder.BuildErr: cannot set %q twice", field))
	}
	c[field] = true
	b.fields = c
}

// cloneFields clone the fields to avoid mutation
func (b FooBuilder) cloneFields() map[string]bool {
	result := make(map[string]bool)
	for k, v := range b.fields {
		result[k] = v
	}
	return result
}
