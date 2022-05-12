# Validation vs invariant

__Validation__ is the process of approving a given object state, while __invariant__ enforcement happens before that state has even been reached [^1]


## Kinds of validation

- input validation (normally at presentation/controller layer)
- domain validation (at model layer)
- always valid model where validation is carried during construction of model
- deferred validation where validation can only be done through external means outside the domain layer (e.g. checking db for unique emails, calling API for OTP verification)


## Preventing Invariants

The best way to prevent invariants is to 
1) create always valid objects, which is only allow creating a type through constructor
2) keep the fields private and only expose them through setter and getter. When setting the new value, ensure that validation is carried out.

However, that does not apply to all language, e.g. golang, where you can declare a new variable without enforcing constructor. One possible solution is to add deferred validation, through a `Validate` method, and wrap the types in interface containers that requires the `Validate` method to be called.

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
	ErrUserNotFound  = errors.New("user not found")
	ErrEmailNotFound = errors.New("email not found")
	ErrInvalidEmail  = errors.New("invalid email")
)

func main() {
	u := &User{
		email:       NewEmail("john.doe@mail.com"),
		PublicEmail: NewValidator(NewEmail("john.doe@mail.com")),
	}
	fmt.Println(u.Validate(), "// calling u.Validate()")

	email, err := u.Email()
	if err != nil {
		panic(err)
	}

	printEmail(NewValidator(email))
	printEmail(u.PublicEmail)

	fmt.Println(email.Value())
}

// Wrapping the Email in ToValidate enforces that Validate method be called before performing an action.
func printEmail(v ToValidate[*Email]) {
	if v == nil {
		panic("unknown validator")
	}
	email, err := v.Validate()
	if err != nil {
		panic(err)
	}
	fmt.Println("email", email)
}

type validatable interface {
	Validate() error
}

// ToValidate wraps a type that requires deferred validation.
type ToValidate[T validatable] interface {
	Validate() (T, error)
}

type Validator[T validatable] struct {
	value T
}

func NewValidator[T validatable](t T) *Validator[T] {
	return &Validator[T]{t}
}

func (r *Validator[T]) Validate() (T, error) {
	return r.value, r.value.Validate()
}

// Email value object.
type Email struct {
	value string
}

func NewEmail(email string) *Email {
	return &Email{email}
}

func (e *Email) Validate() error {
	if e == nil {
		return ErrEmailNotFound
	}
	if !strings.Contains(e.value, "@") {
		return ErrInvalidEmail
	}
	return nil
}

func (e *Email) Value() (string, error) {
	return e.value, e.Validate()
}

type User struct {
	// For private field, we can control the visibility with methods and enforce that .Validate is called.
	email *Email

	// For public field, we can "decorate" the field with another type.
	PublicEmail ToValidate[*Email]
}

func (u *User) Validate() error {
	if u == nil {
		return ErrUserNotFound
	}

	if err := u.email.Validate(); err != nil {
		return err
	}
	return nil
}

func (u *User) Email() (*Email, error) {
	return u.email, u.email.Validate()
}
```


# References
[^1]: [StackOverflow: What is the difference between invariants and validation rules](https://stackoverflow.com/questions/30190302/what-is-the-difference-between-invariants-and-validation-rules)

Where does validation happen?

https://softwareengineering.stackexchange.com/questions/270607/where-should-i-place-business-logic-validations

What is the difference between validation and business rule?
https://stackoverflow.com/questions/6631280/what-is-the-difference-between-a-validation-rule-and-a-business-rule
