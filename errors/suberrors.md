# Suberrors and Rules

Most of the time, we would have suberrors derived from a parent error. Say for a given `User` entity, we want to validate the `email` and `password` field if they are invalid. There may be more than one suberrors for each field, e.g. email cannot be empty, or email format must be valid. Suberrors are a more elegant way of handling growing errors, compared to using error codes.


```go
package main

import (
	"errors"
	"fmt"
	"strings"
)

const minPasswordLength = 8

var (
	// Error.
	ErrInvalidEmail = errors.New("user: invalid email")

	// Suberrors.
	ErrEmailRequired    = fmt.Errorf("%w: cannot be empty", ErrInvalidEmail)
	ErrEmailWrongFormat = fmt.Errorf("%w: wrong format", ErrInvalidEmail)

	ErrInvalidPassword  = errors.New("user: invalid password")
	ErrPasswordRequired = fmt.Errorf("%w: cannot be empty", ErrInvalidPassword)
	ErrPasswordTooShort = fmt.Errorf("%w: too short, min %d characters", ErrInvalidPassword, minPasswordLength)
)

type AggregateError struct {
	errors []error
}

func (a AggregateError) Error() string {
	messages := make([]string, len(a.errors))
	for i, err := range a.errors {
		messages[i] = err.Error()
	}
	return strings.Join(messages, "\n")
}

func (a AggregateError) Is(err error) bool {
	for _, e := range a.errors {
		if errors.Is(e, err) {
			return true
		}
	}
	return false
}

func (a *AggregateError) Add(err error) {
	a.errors = append(a.errors, err)
}

func NewAggregateError() *AggregateError {
	return &AggregateError{}
}

type User struct {
	Email             string
	EncryptedPassword string
}

func encrypt(plaintext string) string {
	return plaintext
}

func NewUser(email, password string) (User, error) {
	err := NewAggregateError()

	if email == "" {
		err.Add(ErrEmailRequired)
	} else if !strings.Contains(email, "@") {
		err.Add(ErrEmailWrongFormat)
	}

	if password == "" {
		err.Add(ErrPasswordRequired)
	} else if len(password) < 8 {
		err.Add(ErrPasswordTooShort)
	}

	return User{
		Email:             email,
		EncryptedPassword: encrypt(password),
	}, err
}

func main() {
	_, err := NewUser("john.doe", "5")
	fmt.Println(err)
	fmt.Println(errors.Is(err, ErrInvalidPassword))
	fmt.Println(errors.Is(err, ErrPasswordRequired))
	fmt.Println(errors.Is(err, ErrPasswordTooShort))
}
```
