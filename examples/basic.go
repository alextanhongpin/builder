package main

//go:generate go run ../main.go -type User -out=basic_gen.go
type User struct {
	name    string
	age     int64
	married bool
	remarks *string `build:"-"`
	hobbies []string
}
