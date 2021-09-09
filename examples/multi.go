package main

import "github.com/google/uuid"

//go:generate go run ../main.go -type Hello,World
type Hello struct {
	id uuid.UUID
}

type World struct {
	id uuid.UUID
}
