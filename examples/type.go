package main

import (
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

//go:generate go run ../main.go -type Book -out=type_gen.go
type Book struct {
	id    uuid.UUID
	name  sql.NullString
	extra json.RawMessage
}
