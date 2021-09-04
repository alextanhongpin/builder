package main

type Year int64

//go:generate go run ../main.go -type Author -out=slice_gen.go
type Author struct {
	user     *User
	books    []Book
	releases []*Year
}
