package main

import "database/sql"

//go:generate go run ../main.go -type Foo
type Foo struct {
	id      int64
	name    string
	age     sql.NullInt64
	valid   *bool
	url     string
	realAge *int64
}
