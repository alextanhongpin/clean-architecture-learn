# Value Object

Primitive obsession is an anti-pattern, and value objects are highly recommended. However, some languages does not offer the required functionality to create a value object in its valid state.

Take for example `golang`:

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

The only workaround is to use interface:

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
