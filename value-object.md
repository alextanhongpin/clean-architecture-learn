# Value Object


Primitive obsession is an anti-pattern, and value objects are highly recommended. However, some languages does not offer the required functionality to create a value object in its valid state.

- are DTOs value object? No, see [1]. DTO is a class representing some data with no logic in it.
- Value Object is a full member of your domain model
- The only difference between Value Object and Entity is that Value Object doesn’t have its own identity.
- Value Object are `comparable` by values
- Value Objects do contain logic and, typically, they are not used for transferring data between application boundaries.
- value object does not have identity, so no id field
- good contender of value object is stock quantity, age, email

# Example

Take for example `golang`, modelling __value object__ through type definition is suboptimal:

```go
package main

import "fmt"

type Password string

func NewPassword(s string) Password {
	return Password(s)
}

func main() {
	var p Password // Already in invalid state
	fmt.Println(p)
}
```

One workaround is to use interface:

```go
package main

import (
	"errors"
	"fmt"
	"log"
)

type Password interface {
	Password() string
}

func IsPassword(in interface{}) bool {
	_, ok := in.(Password)
	return ok
}

type password string

func NewPassword(p string) (Password, error) {
	if p == "" {
		return password(""), errors.New("password: cannot be empty")
	}
	return password(p), nil
}

func (p password) Password() string {
	return string(p)
}

func (p password) String() string {
	return "**REDACTED**"
}

func main() {
	p, err := NewPassword("secret")
	if err != nil {
		log.Fatal(err)
	}
	doStuff(p)
	fmt.Println(IsPassword(p))
}

func doStuff(p Password) {
	fmt.Println(p.Password(), p)
}
```

Or wrap the primitive in a struct:

```go
package main

import (
	"errors"
	"fmt"
	"log"
)

type Password struct {
	value string
}

func NewPassword(v string) (Password, error) {
	if len(v) == 0 {
		return Password{}, errors.New("password: cannot be empty")
	}
	if len(v) < 8 {
		return Password{}, errors.New("password: too short")
	}
	return Password{
		value: v,
	}, nil
}

func (p Password) Value() string {
	return p.value
}

func (p Password) String() string {
	return "**REDACTED**"
}

func main() {
	p, err := NewPassword("hello")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(p)
}
```

## The always valid

Given the `Email` value object:
```go
type Email struct {
	value string
}

func NewEmail(v string) (Email, error) {
	if len(v) == 0 {
		return Email{}, errors.New("email: cannot be empty")
	}
	if !isEmail(v) {
		return Email{}, errors.New("email: invalid")
	}
	return Email{
		value: v,
	}, nil
}
```

A valid email can only be build from the constructor. However, when loading from the database, the value can be empty too, so instead of returning pointer email, we return a null object pattern.

```go
// Not this
func NewEmail(v string) (*Email, error) {}

// Do this
func NewEmail(v string) (Email, error) {}
```

This allows us to skip the error when loading from the db, while creating a valid object:
```go
email, _ := NewEmail(emailFromDB)
```

The takeaway is, when reading value objects, they can be invalid (if they are not set). However, when writing, they have to be valid.

## My take on value object now

Given the same `Password` example above, I would have write it this way now:

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"errors"
	"fmt"
)

const MinPasswordLen = 8

var (
	ErrEmptyPassword    = errors.New("password is empty")
	ErrPasswordNotSet   = errors.New("password not set")
	ErrPasswordTooShort = errors.New("password too short")
)

type Password struct {
	// This is still preferred over "type Password string"
	// because the only way to set the "value" is through the constructor.
	value       string
	constructed bool
}

func (p *Password) Validate() error {
	// Not set is not the same as empty.
	if p == nil || !p.constructed {
		return ErrPasswordNotSet
	}

	if p.value == "" {
		return ErrEmptyPassword
	}

	if len(p.value) < MinPasswordLen {
		return ErrPasswordTooShort
	}

	return nil
}

func NewPassword(v string) (*Password, error) {
	p := &Password{value: v, constructed: true}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p Password) Value() (string, error) {
	return p.value, p.Validate()
}

func (p Password) String() string {
	if !p.constructed {
		return "NOT SET"
	}

	// This condition below is not possible.
	if p.value == "" {
		return "EMPTY"
	}

	return "**REDACTED**"
}

func main() {
	p, err := NewPassword("hello world")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p.Value())
	fmt.Println(p.Validate(), "// p")

	var p2 Password
	fmt.Println(p2.Validate(), p2, "// p2")

	var p3 *Password
	fmt.Println(p3.Validate(), p3, "// p3")

	p4, err := NewPassword("")
	if err != nil {
		fmt.Println(err, p4, "// p4")
	}
}
```

# References

1. [DTO vs Value Object vs POCO](https://enterprisecraftsmanship.com/posts/dto-vs-value-object-vs-poco/#:~:text=DTO%20is%20a%20class%20representing%20some%20data%20with%20no%20logic%20in%20it.&text=On%20the%20other%20hand%2C%20Value,t%20have%20its%20own%20identity.)
2. http://gorodinski.com/blog/2012/05/19/validation-in-domain-driven-design-ddd/
