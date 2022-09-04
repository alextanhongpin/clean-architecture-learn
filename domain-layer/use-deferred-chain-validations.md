# Use deferred chain validation

This topic covers the following question
- how to validate domain model
- how to create always valid domain model
- where do validation belong
- how to deal with annoying error checking
- hot to construct and validate value object correctly

In golang, the type system does not allow creation of always valid domain model. Any attempt at implementing them is futile, and hence we will stick with deferred chain validation instead.


To first illustrate the problem, let us create an age `value object`:


```go
package main

import "fmt"

func main() {
	age, err := NewAge(10)
	fmt.Println(age, err)

	var badAge Age
	fmt.Println(badAge)
}

const LegalAge = 13

var ErrIllegalAge = fmt.Errorf("age: must be at least %d years old", LegalAge)

// Value object using primitives.
type Age int

func NewAge(age int) (Age, error) {
	if age < 13 {
		return 0, ErrIllegalAge
	}
	return Age(age), nil
}
```

We can see that the variable `badAge` itself could be initialized without going through the constructor `NewAge`, regardless of how it is designed:

```go
package main

import "fmt"

func main() {
	age, err := NewAge(10)
	fmt.Println(age, err)

	// Still cannot be avoided.
	var badAge Age
	fmt.Println(badAge)
}

// Value object using struct.
type Age struct {
	n int
}

func NewAge(age int) (*Age, error) {
	if age < 13 {
		return nil, ErrIllegalAge
	}
	return &Age{n: age}, nil
}
```

The above implementation is just another variation. However, the author prefers the primitive version since they can be serialized/deserialized from JSON data easily versus the struct version.

In reality, declaration of variables without going through constructor might not be a huge concern, because most of the time you will be constructing domain model from
- API layer, through the request object
- Database layer, when hydrating an entity from the database

However, the `NewAge` constructor may still pose a problem when you start using value objects for all your fields:


```go
package main

import (
	"errors"
	"fmt"
	"unicode"
)

type User struct {
	Age  Age
	Name Name
}

func NewUser(name string, age int) (*User, error) {
	n, err := NewName(name)
	if err != nil {
		return nil, err
	}
	a, err := NewAge(age)
	if err != nil {
		return nil, err
	}

	return &User{
		Name: n,
		Age:  a,
	}, nil
}

const LegalAge = 13

var ErrIllegalAge = fmt.Errorf("age: must be at least %d years old", LegalAge)

type Age int

func NewAge(age int) (Age, error) {
	if age < 13 {
		return 0, ErrIllegalAge
	}
	return Age(age), nil
}

var ErrBadNameFormat = errors.New("name: can only contain alphabets and/or spaces")

type Name string

func NewName(name string) (Name, error) {
	for _, r := range name {
		valid := unicode.IsLetter(r) || unicode.IsSpace(r)
		if !valid {
			return "", ErrBadNameFormat
		}
	}

	return Name(name), nil
}
```

There is a lot of error handling, and we do not even have an option to re-validate the name and age too.

We can design better:

```go
package main

import (
	"errors"
	"fmt"
	"unicode"
)

func main() {
	fmt.Println(NewUser("john", 13))
}

type User struct {
	Age  Age
	Name Name
}

func (u *User) Validate() error {
	if err := u.Age.Validate(); err != nil {
		return err
	}
	return u.Name.Validate()
}

func NewUser(name string, age int) (*User, error) {
	u := &User{
		Name: Name(name),
		Age:  Age(age),
	}
	if err := u.Validate(); err != nil {
		return nil, err
	}
	return u, nil
}

const LegalAge = 13

var ErrIllegalAge = fmt.Errorf("age: must be at least %d years old", LegalAge)

type Age int

func (age Age) Validate() error {
	if age < 13 {
		return ErrIllegalAge
	}
	return nil
}

var ErrBadNameFormat = errors.New("name: can only contain alphabets and/or spaces")

type Name string

func (name Name) Validate() error {
	for _, r := range name {
		valid := unicode.IsLetter(r) || unicode.IsSpace(r)
		if !valid {
			return ErrBadNameFormat
		}
	}
	return nil
}
```

With the solution above, we can now:
- validate each field independently by just calling `Validate`
- chain the validation in a single place, which is the `user.Validate()`
- we can still use constructor that calls the `Validate` method
- avoid individual constructors with errors (`NewName` and `NewAge` is not required at all, since the User will be responsible for checking if all fields are valid). Also, constructor that returns errors might not be idiomatic.
- while many people have rely on constructor to perform validation, itpose an issuesince the validation cannot be reused elsewhere, or may become hard to test if there are simply too many branches. Also, the goal of constructor should purely be constructing, not validating.
