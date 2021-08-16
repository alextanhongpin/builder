# builder

Generate golang builders.


## Installation

```bash
$ go install github.com/alextanhongpin/builder
```

## Usage

```go
//go:generate go run ../main.go -type Simple
type Simple struct {
	name string
	age  int `build:"-"`
}
```

Output:

https://github.com/alextanhongpin/builder/blob/037167612b0a546e481560e45d75406345dbff4d/examples/simple_gen.go#L1-L70

## Ignoring and setting custom setter.

Sometimes you want a custom type, but also need to take advantage of the `Build()` which panics if not all the fields are set
```go
// Extend simple builder and check if the field is set.
func (s SimpleBuilder) WithCustomAge(age int) SimpleBuilder {
	s.mustSet("age")
	s.simple.age = age
	return s
}
```

## Build


```go
func main() {
	builder := NewSimpleBuilder("age")                              // Pass custom fields that needs to be set.
	log.Println(builder.BuildPartial())                             // Allows the entity to be build partially.
	log.Println(builder)                                            // None of the values are set yet.
	log.Println(builder.WithName("john"))                           // name is set to true
	log.Println(builder.WithName("john").WithCustomAge(10).Build()) // name and age set and build success
	log.Println(builder)                                            // Every instance is immutable and they don't share state.
	log.Println(builder.Build())                                    // This will panic, since "name" and "age" is not set yet.
}
```
