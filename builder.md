# Fluent Builder Pattern

```go
package main

import (
	"fmt"
)

// https://levelup.gitconnected.com/building-immutable-data-structures-in-go-56a1068c76b2
func main() {
	fmt.Println(NewPerson("jane"))
	pb := NewPersonBuilder().
		WithName("john").
		WithFavouriteColors("red", "blue")

	person := pb.Build()
	fmt.Println(person)
	fmt.Println(person.FavouriteColorAt(0))
	pb.Build()
}

type Person struct {
	name            string
	favouriteColors []string
}

// When only select arguments are required.
func NewPerson(name string) Person {
	return NewPersonBuilder().
		WithName(name).
		Build()
}

func (p Person) Name() string {
	return p.name
}

func (p Person) FavouriteColorAt(i int) string {
	return p.favouriteColors[i]
}

// PersonBuilder with fluent interface. An advantage is we can build a valid person without passing all arguments to the constructor.
type PersonBuilder struct {
	p *Person
}

func NewPersonBuilder() *PersonBuilder {
	return &PersonBuilder{p: new(Person)}
}

func (b *PersonBuilder) WithName(name string) *PersonBuilder {
	b.p.name = name
	return b
}

func (b *PersonBuilder) WithFavouriteColors(colors ...string) *PersonBuilder {
	b.p.favouriteColors = colors
	return b
}

// Build returns a new Person (optionally an error).
func (b *PersonBuilder) Build() Person {
	// Create a copy.
	p := *b.p

	// TODO: Validate Person invariant, at least the required fields.
	if p.name == "" {
		panic("person name is required")
	}

	// We have a few options here.
	// 1. Consume after used, so that calling the second time will panic.
	// 2. Return a copy, so that we can reuse the builder to build another person after modifying other withers.
	b.p = nil
	return p
}
```
