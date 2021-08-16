package main

import (
	"database/sql"
	"log"
)

type Bar string

//go:generate go run ../main.go -type Foo
type Foo struct {
	id                  int64
	name                string
	age                 sql.NullInt64
	valid               *bool
	url                 string
	realAge             *int64
	bar                 Bar
	skip                string `build:"-"`
	bars                []Bar
	barPtrs             []*Bar
	barByString         map[string]Bar
	stringByBar         map[Bar]string
	sliceBarByString    map[string][]Bar
	sliceBarPtrByString map[string][]*Bar
	barPtrByString      map[string]*Bar
	stringByBarPtr      map[*Bar]string
}

//go:generate go run ../main.go -type Simple
type Simple struct {
	name string
	age  int `build:"-"`
}

// Extend simple builder and check if the field is set.
func (s SimpleBuilder) WithCustomAge(age int) SimpleBuilder {
	s.mustSet("age")
	s.simple.age = age
	return s
}

func main() {
	builder := NewSimpleBuilder("age")
	log.Println(builder)                  // None of the values are set yet.
	log.Println(builder.WithName("john")) // name is set to true
	//log.Println(builder.WithName("john").WithName("jessie"))        // panic on setting twice.
	log.Println(builder.WithName("john").WithCustomAge(10))         // name and age is set to true.
	log.Println(builder.WithName("john").WithCustomAge(10).Build()) // Build succeeds when all fields are set.
	log.Println(builder.BuildPartial())                             // Build succeeds even when not all fields are set.
	log.Println(builder)                                            // Every instance is immutable and they don't share state.
	log.Println(builder.Build())                                    // This will panic, since "name" and "age" is not set yet.
}
