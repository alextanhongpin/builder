package main

import "time"

//go:generate go run ../main.go -type ConfirmationToken -out=pointer_gen.go
type ConfirmationToken struct {
	expiresAt *time.Time
	valid     *bool
	reference *string
}
