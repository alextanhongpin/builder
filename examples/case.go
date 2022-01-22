package main

//go:generate go run ../main.go -type Account -out=case_gen.go
type Account struct {
	ID   int64
	Typ  string
	Name string
}
