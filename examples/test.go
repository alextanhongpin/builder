package main

import (
	"database/sql"
	"log"
)

type Bar string

//go:generate go run ../main.go -type Foo,Simple
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

type Simple struct {
	name string
}

func main() {
	builder := NewSimpleBuilder()
	log.Println(builder.BuildPartial())
	log.Println(builder)
	log.Println(builder.WithName("john"))
	log.Println(builder.WithName("john").Build())
	log.Println(builder)
	log.Println(builder.Build())
}
