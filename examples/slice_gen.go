// Code generated by github.com/alextanhongpin/builder, DO NOT EDIT.

package main

import "fmt"

type AuthorBuilder struct {
	author    Author
	fields    map[string]int
	fieldsSet uint64
}

func NewAuthorBuilder() *AuthorBuilder {
	fields := make(map[string]int)
	for i, field := range []string{"user", "books", "releases"} {
		fields[field] = i
	}
	return &AuthorBuilder{fields: fields}
}

// WithUser sets user.
func (b AuthorBuilder) WithUser(user User, valid bool) AuthorBuilder {
	if valid {
		b.author.user = &user
	}
	b.Set("user")
	return b
}

// WithBooks sets books.
func (b AuthorBuilder) WithBooks(books []Book) AuthorBuilder {
	b.author.books = books
	b.Set("books")
	return b
}

// WithReleases sets releases.
func (b AuthorBuilder) WithReleases(releases []*Year) AuthorBuilder {
	b.author.releases = releases
	b.Set("releases")
	return b
}

// Build returns Author.
func (b AuthorBuilder) Build() Author {
	for field := range b.fields {
		if !b.IsSet(field) {
			panic(fmt.Errorf("builder: %q not set", field))
		}
	}
	return b.author
}

// Build returns Author.
func (b AuthorBuilder) BuildPartial() Author {
	return b.author
}

func (b *AuthorBuilder) Set(field string) bool {
	n, ok := b.fields[field]
	if !ok {
		return false
	}
	b.fieldsSet |= 1 << n
	return true

}

func (b AuthorBuilder) IsSet(field string) bool {
	pos := b.fields[field]
	return (b.fieldsSet & (1 << pos)) == (1 << pos)
}

func (b *AuthorBuilder) Register(field string) error {
	if _, ok := b.fields[field]; ok {
		return fmt.Errorf("field %q already registered", field)
	}
	b.fields[field] = len(b.fields)
	return nil
}
