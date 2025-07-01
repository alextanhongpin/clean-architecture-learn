# Use getter/setter (or wither)

In golang, getters and setters is not a language feature [^1].

Getters and setters may be necessary if your struct fields are mostly private. However, exposing a getter/setter for every field may not be idiomatic:


```go
type User struct {
	id   int
	name string
	age  int64
}

func (u *User) ID() int             { return u.id }
func (u *User) SetID(id int)        { u.id = id }

func (u *User) Name() string        { return u.name }
func (u *User) SetName(name string) { u.name = name }

func (u *User) Age() int64          { return u.age }
func (u *User) SetAge(age int64)    { u.age = age }
```

It would have been easier to just use public fields:

```go
type User struct {
	ID   int
	Name string
	Age  int64
}
```


However, if you have need to prevent invariant when setting a value, then the setter would have serve a purpose:


```go
package main

import (
	"fmt"
)

func main() {
	user := new(User)
	if err := user.SetAge(10); err != nil {
		panic(err)
	}
	fmt.Println(user)
}

const LegalAge = 13

var ErrIllegalAge = fmt.Errorf("user: must be at least %d years old", LegalAge)

type User struct {
	id   int
	name string
	age  int64
}

func (u *User) SetAge(age int64) error {
	if age < LegalAge {
		return ErrIllegalAge
	}
	u.age = age

	return nil
}
```

> Avoid pure setters without behaviour. Use setter to prevent invariants.

However, an even better approach is to use _value object_ instead to prevent invariants.

In the example below, we go back to exposing all fields as public fields, and introduce the age `value object`.

```go
package main

import (
	"fmt"
)

func main() {
	user := new(User)
	user.Age = Age(10)
	fmt.Println(user.Validate())
}

const LegalAge = 13

var ErrIllegalAge = fmt.Errorf("user: must be at least %d years old", LegalAge)

type User struct {
	ID   int
	Name string
	Age  Age
}

func (u *User) Validate() error {
	return u.Age.Validate()
}

type Age int64

func (a Age) Validate() error {
	if a < LegalAge {
		return ErrIllegalAge
	}

	return nil
}
```

Some may question the method `Validate()`, as it _defers_ the validation and go against the _always valid domain model_ [^2]. This will be addressed in a separate chapter by itself, related to errors handling.

[^1]: https://go.dev/doc/effective_go#Getters
[^2]: https://enterprisecraftsmanship.com/posts/always-valid-domain-model/
