# Avoid Builder

Not to be confused with the GoF builder pattern, we are referring to the builder pattern similar to Java that uses method chaining to build object instead of constructor.

You may be tempted to do the following in golang:

```go
package main

import "fmt"

func main() {
	user := NewUserBuilder().
		SetID(1).
		SetName("John Appleseed").
		SetAge(13).
		Build()
	fmt.Println(user)
}

type User struct {
	id   int
	name string
	age  int
}

type UserBuilder struct {
	user User
}

func (u *UserBuilder) SetID(id int) *UserBuilder {
	u.user.id = id
	return u
}
func (u *UserBuilder) SetName(name string) *UserBuilder {
	u.user.name = name
	return u
}

func (u *UserBuilder) SetAge(age int) *UserBuilder {
	u.user.age = age
	return u
}

func (u *UserBuilder) Build() *User {
	return &u.user
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}
```


Prefer simple constructor instead:

```go
package main

import "fmt"

func main() {
	user := NewUser(1, "John Appleseed", 13)
	fmt.Println(user)
}

type User struct {
	id   int
	name string
	age  int
}

func NewUser(id int, name string, age int) *User {
	return &User{
		id:   id,
		name: name,
		age:  age,
	}
}
```

It is shorter and simpler, and makes every field _mandatory_.

This example is simple, and we are skipping a lot of steps such as:

- using value object instead of primitives
- validation to ensure always valid object created (e.g. age should not be negative, and above 13)
- growing number of fields (requires modification to all the places calling the constructor)

However, in most scenarios, sticking to plain constructor keep things maintainable.

Also, if your struct contains mostly public fields, you could introduce a mapping layer instead to map types from one layer to another.

Using public fields avoids the mega constructor, where you have to define all private fields for the construction of a large struct.

```go
package main

import (
	"fmt"

	"play.ground/repository"
)

func main() {
	user := repository.NewUser(repository.User{
		ID:   1,
		Name: "John Appleseed",
		Age:  13,
	})
	fmt.Println(user)
}
-- go.mod --
module play.ground
-- domain/user.go --
package domain

type User struct {
	ID   int
	Name string
	Age  int64
}
-- repository/user_repository.go --
package repository

type User struct {
	ID   int
	Name string
	Age  int64
}
-- repository/user.go --
// Package user.go contains mapping from repository.User to domain.User
package repository

import "play.ground/domain"

func NewUser(u User) domain.User {
	return domain.User{
		ID:   u.ID,
		Name: u.Name,
		Age:  u.Age,
	}
}
```
