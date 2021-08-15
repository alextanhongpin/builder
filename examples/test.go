package main

import (
	"database/sql"
)

type Bar string

//go:generate go run ../main.go -type Foo
type Foo struct {
	id      int64
	name    string
	age     sql.NullInt64
	valid   *bool
	url     string
	realAge *int64
	bar     Bar
	skip    string `build:"-"`
}
