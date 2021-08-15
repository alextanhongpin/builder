// Code generated by builder, DO NOT EDIT.
package main

import (
	"database/sql"
	"fmt"
)

type FooBuilder struct {
	foo       Foo
	fields    []string
	fieldsSet uint64
}

func NewFooBuilder() *FooBuilder {
	return &FooBuilder{fields: []string{"id", "name", "age", "valid", "url", "realAge", "bar", "bars", "barByString", "stringByBar", "sliceBarByString"}}
}

// WithID sets id.
func (b FooBuilder) WithID(id int64) FooBuilder {
	b.fieldsSet |= 1 << 0
	b.foo.id = id
	return b
}

// WithName sets name.
func (b FooBuilder) WithName(name string) FooBuilder {
	b.fieldsSet |= 1 << 1
	b.foo.name = name
	return b
}

// WithAge sets age.
func (b FooBuilder) WithAge(age sql.NullInt64) FooBuilder {
	b.fieldsSet |= 1 << 2
	b.foo.age = age
	return b
}

// WithValid sets valid.
func (b FooBuilder) WithValid(valid bool, valid1 bool) FooBuilder {
	b.fieldsSet |= 1 << 3
	if valid1 {
		b.foo.valid = &valid
	}
	return b
}

// WithURL sets url.
func (b FooBuilder) WithURL(url string) FooBuilder {
	b.fieldsSet |= 1 << 4
	b.foo.url = url
	return b
}

// WithRealAge sets realAge.
func (b FooBuilder) WithRealAge(realAge int64, valid bool) FooBuilder {
	b.fieldsSet |= 1 << 5
	if valid {
		b.foo.realAge = &realAge
	}
	return b
}

// WithBar sets bar.
func (b FooBuilder) WithBar(bar Bar) FooBuilder {
	b.fieldsSet |= 1 << 6
	b.foo.bar = bar
	return b
}

// WithBars sets bars.
func (b FooBuilder) WithBars(bars []Bar) FooBuilder {
	b.fieldsSet |= 1 << 7
	b.foo.bars = bars
	return b
}

// WithBarByString sets barByString.
func (b FooBuilder) WithBarByString(barByString map[string]Bar) FooBuilder {
	b.fieldsSet |= 1 << 8
	b.foo.barByString = barByString
	return b
}

// WithStringByBar sets stringByBar.
func (b FooBuilder) WithStringByBar(stringByBar map[Bar]string) FooBuilder {
	b.fieldsSet |= 1 << 9
	b.foo.stringByBar = stringByBar
	return b
}

// WithSliceBarByString sets sliceBarByString.
func (b FooBuilder) WithSliceBarByString(sliceBarByString map[string][]Bar) FooBuilder {
	b.fieldsSet |= 1 << 10
	b.foo.sliceBarByString = sliceBarByString
	return b
}

// Build returns Foo.
func (b FooBuilder) Build() Foo {
	for i := 0; i < len(b.fields); i++ {
		if (b.fieldsSet & 1 << i) != 1<<i {
			panic(fmt.Sprintf("builder.BuildErr: %q not set", b.fields[i]))
		}
	}
	return b.foo
}

// Build returns Foo.
func (b FooBuilder) BuildPartial() Foo {
	return b.foo
}