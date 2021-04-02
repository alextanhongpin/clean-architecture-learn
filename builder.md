# Fluent Builder Pattern

Why not just use a constructor?
- entities attributes can grow
- order in constructor is important

Why separate builder from entity?
- entity should have mostly read only fields
- setters should have invariants, but normally when loading from db, they can be empty first

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

## Improvements
```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, playground")
}

type personValidator interface {
	Validate(p *Person) error
}

// PersonBuilder with fluent interface. An advantage is we can build a valid person without passing all arguments to the constructor.
type PersonBuilder struct {
	p      *Person
	errors []error
	validator personValidator
}

func (b *PersonBuilder) WithAge(age int) *PersonBuilder {
	// Handle validation per field. (Or better, use value object).
	if age > 150 || age < 0 {
		b.errors = append(b.errors, ErrInvalidAge)
		return b
	}
	b.p.age = age
	return b
}
func (b *PersonBuilder) WithName(name string) *PersonBuilder {
	// Convert primitives into value object/domain primitives.
	username, err := NewUsername(name)
	if err != nil {
		b.errors = append(b.errors, err)
		return b
	}
	b.p.name = username
	return b
}

// Alternatively, we can pass a validator to validate the struct as a whole.
func (b *PersonBuilder) WithValidator(validator personValidator) *PersonBuilder {
	b.validator = validator
	return b
}

// Build builds and validate the person entity.
func (b *PersonBuilder) Build() (Person, error) {
	if err := b.validator.Validate(b.p); err != nil {
		return Person{}, err
	}
	// REDACTED
	return p, nil
}


// Another alternative is to always pass a validation during build time.
func (b *PersonBuilder) Build(validator personValidator) (Person, error) {
	if err := validator.Validate(b.p); err != nil {
		return Person{}, err
	}
	// REDACTED
	return p, nil
}
```
## Building nested entity

```go
package main

import (
	"fmt"
)

type Address struct {
	street string
}

type AddressBuilder struct {
	address     *Address
	userBuilder *UserBuilder
}

func (b *AddressBuilder) WithStreet(street string) *AddressBuilder {
	b.address.street = street
	return b
}

func (b *AddressBuilder) Build() *Address {
	a := b.address
	b.address = nil
	return a
}

func (b *AddressBuilder) BuildUserAddress() *UserBuilder {
	a := b.Build()
	b.userBuilder.user.address = *a
	return b.userBuilder
}

type User struct {
	name    string
	address Address
}
type UserBuilder struct {
	user *User
}

func (b *UserBuilder) WithName(name string) *UserBuilder {
	b.user.name = name
	return b
}

func (b *UserBuilder) AddAddress() *AddressBuilder {
	return &AddressBuilder{
		address:     new(Address),
		userBuilder: b,
	}
}

func (b *UserBuilder) Build() *User {
	user := b.user
	b.user = nil
	return user
}

func main() {
	b := &UserBuilder{
		user: new(User),
	}
	user := b.WithName("john doe").
		AddAddress().
		WithStreet("street 1").
		BuildUserAddress().
		Build()
	fmt.Println(user)
}
```
# References
https://www.ssw.com.au/rules/rules-to-better-clean-architecture


