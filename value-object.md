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

## 2022: My take on value object now

Given the same `Password` example above, I would have write it this way now below.

The major changes are:
- a `Validate` method is added to encapsulate the validation method
- the `Validate` method is now called during construction, as well as when returning the `Value`, this ensures that the value will always be valid
- a new private field `constructed` is added, to check if the values are constructed through the constructor. This prevents declaration of the variables or variable pointer since by default the boolean value will always be false


```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordNotSet   = errors.New("password not set")
	ErrPasswordTooShort = errors.New("password too short")
)

func main() {
	hash, err := NewPassword("hello world")
	if err != nil {
		return
	}
	fmt.Println(hash)
	fmt.Println(hash.Match("hello world"))
	hashValue, err := hash.Value()
	if err != nil {
		panic(err)
	}

	pwd := NewPasswordFromHash(hashValue)
	fmt.Println(pwd.Match("hello world"))

	var pwd2 *Password
	fmt.Println(pwd2.Value())
}

type Password struct {
	// This ensures that structs must be initialized with keys
	// e.g `&Password{constructed: true}` is allowed,
	// but `&Password{true}` is not allowed
	_           struct{}
	hash        string
	constructed bool
}

func NewPasswordFromHash(hash string) *Password {
	return &Password{
		constructed: true,
		hash:        hash,
	}
}

func NewPassword(password string) (*Password, error) {
	pwd := &Password{constructed: true}
	pwd.encrypt(password)
	return pwd, pwd.Validate()
}

func (p *Password) Validate() error {
	/* We need to check p == nil to guard against variable pointer declaration from panicking, e.g
	var pwd *Password
	pwd.Value() // will panic
	*/
	if p == nil || !p.constructed {
		return ErrPasswordNotSet
	}

	return nil
}

func (p *Password) validate(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}

func (p *Password) encrypt(password string) error {
	if err := p.validate(password); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.hash = string(hash)

	return nil
}

func (p *Password) Match(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(password))
	return err == nil
}

func (p *Password) String() string {
	if !p.constructed {
		return "NOT SET"
	}

	return "**REDACTED**"
}

func (p *Password) Value() (string, error) {
	// Early validation to capture nil pointer scenario.
	if err := p.Validate(); err != nil {
		return "", err
	}

	return p.hash, nil
}
```


And the email example, which handles scenario where other developers can just declare the variables:

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmailNotSet   = errors.New("email not set")
	ErrEmailRequired = errors.New("email is required")
	ErrEmailInvalid  = errors.New("email is invalid")
)

func main() {
	email, err := NewEmail("john.doe@mail.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(email.Value())

	var email2 *Email
	fmt.Println(email2.Value())
	fmt.Println(email2.Validate())

	var email3 Email
	fmt.Println(email3.Value())
	fmt.Println(email3.Validate())
}

type Email struct {
	value       string
	constructed bool
}

/* NOTE: This will cause panic, since the e is nil
func (e *Email) Value() (string, error) {
	return e.value, e.Validate()
}
*/

func (e *Email) Value() (string, error) {
	if err := e.Validate(); err != nil {
		return "", err
	}
	return e.value, nil
}

func (e *Email) Validate() error {
	if e == nil || !e.constructed {
		return ErrEmailNotSet
	}
	if e.value == "" {
		return ErrEmailRequired
	}

	// Naive checking - don't do this in production, this is only for demonstration purpose.
	if !strings.Contains(e.value, "@") {
		return ErrEmailInvalid
	}
	return nil
}

func NewEmail(v string) (*Email, error) {
	email := &Email{value: v, constructed: true}
	if err := email.Validate(); err != nil {
		return nil, err
	}
	return email, nil
}
```

## Defining value object

For the value object created, we might just want to expose the primitive type and not caring about the actual implementation.

```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	u := User{
		email: &Email{value: "john@mail.com"},
	}
	fmt.Println(u)
	fmt.Println("Hello, 世界")
}

type Validatable interface {
	Validate() error
}

type Value[T any] interface {
	Value() (T, error)
	Validatable
}

type User struct {
	email Value[string] // An value object of type string.
}

type Email struct {
	value string
}

func (e *Email) Value() (string, error) {
	return e.value, nil
}

func (e *Email) Validate() error {
	return nil
}
```

Or perhaps we want to declare the actual type alongside the primitive?

```go
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	u := User{
		email: &Email{value: "john@mail.com"},
	}
	fmt.Println(u)
	fmt.Println("Hello, 世界")
}

type Validatable interface {
	Validate() error
}

type Value[E any, T comparable] interface {
	Value() (T, error)
	Self() E
	Validatable
}

type User struct {
	email Value[*Email, string]
}

type Email struct {
	value string
}

func (e *Email) Value() (string, error) {
	return e.value, nil
}

func (e *Email) Validate() error {
	return nil
}

func (e *Email) Self() *Email {
	return e
}
```

# References

1. [DTO vs Value Object vs POCO](https://enterprisecraftsmanship.com/posts/dto-vs-value-object-vs-poco/#:~:text=DTO%20is%20a%20class%20representing%20some%20data%20with%20no%20logic%20in%20it.&text=On%20the%20other%20hand%2C%20Value,t%20have%20its%20own%20identity.)
2. http://gorodinski.com/blog/2012/05/19/validation-in-domain-driven-design-ddd/
