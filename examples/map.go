package main

import "github.com/google/uuid"

//go:generate go run ../main.go -type Cache -out=map_gen.go
type Cache struct {
	userByID        map[uuid.UUID]*User
	booksByAuthorID map[int64][]Book
	userByIDByName  map[string]map[int64]User
}
