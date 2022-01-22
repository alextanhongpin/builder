package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBasic(t *testing.T) {
	user := NewUserBuilder().
		WithName("john").
		WithAge(10).
		WithMarried(false).
		WithHobbies([]string{"hello", "world"}).
		Build()

	if diff := cmp.Diff(User{
		name:    "john",
		age:     10,
		married: false,
		hobbies: []string{"hello", "world"},
	}, user, cmp.AllowUnexported(User{})); diff != "" {
		t.Fatal(diff)
	}
}

func TestPartialBasic(t *testing.T) {
	user := NewUserBuilder().
		WithName("john").
		WithAge(10).
		WithMarried(false).
		BuildPartial()

	if diff := cmp.Diff(User{
		name:    "john",
		age:     10,
		married: false,
	}, user, cmp.AllowUnexported(User{})); diff != "" {
		t.Fatal(diff)
	}

	_ = NewUserBuilder().BuildPartial()
}

func TestPartialErrorBasic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic when field `hobbies` is not set")
		}
	}()

	_ = NewUserBuilder().
		WithName("john").
		WithAge(10).
		WithMarried(false).
		Build()
}

func (u UserBuilder) WithCustomField(remarks string, valid bool) UserBuilder {
	u.Set("remarks")
	if valid {
		u.user.remarks = &remarks
	}
	return u
}

func TestCustomFieldBasic(t *testing.T) {
	// Supply a new additional fields, `remarks`.
	ub := NewUserBuilder()
	if err := ub.Register("remarks"); err != nil {
		t.Fatalf("failed to register field `remarks`: %v", err)
	}

	user := ub.
		WithName("john").
		WithAge(10).
		WithMarried(false).
		WithHobbies([]string{"hello", "world"}).
		WithCustomField("hello world", true).
		Build()

	remarks := "hello world"
	if diff := cmp.Diff(User{
		name:    "john",
		age:     10,
		married: false,
		hobbies: []string{"hello", "world"},
		remarks: &remarks,
	}, user, cmp.AllowUnexported(User{})); diff != "" {
		t.Fatal(diff)
	}
}

func TestCustomFieldErrorBasic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic when field `remarks` is not set")
		}
	}()

	ub := NewUserBuilder()
	ub.Register("remarks")
	_ = ub.
		WithName("john").
		WithAge(10).
		WithMarried(false).
		WithHobbies([]string{"hello", "world"}).
		Build()
}

func TestCustomFieldDuplicateBasic(t *testing.T) {
	ub := NewUserBuilder()
	err := ub.Register("name")
	if err == nil {
		t.Fatalf("expected error when registering field `name` twice")
	}
	t.Log(err)
}

func TestSetTwiceBasic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic when field `name` is set twice")
		} else {
			t.Log(err)
		}
	}()

	_ = NewUserBuilder().
		WithName("john").
		WithName("john").
		Build()
}
