
```go
package main

import (
	"fmt"
)

func main() {
	u := user{}
	u2 := u.WithName("john")
	printName(u)
	printName(u2)
}

type user struct {
	name string
}

func (u user) Name() string {
	return u.name
}

// Immutable.
func (u user) WithName(name string) user {
	u.name = name
	return u
}

type IName interface {
	Name() string
}

// Using interface for reusability.
func printName(u IName) {
	fmt.Println(u.Name())
}
```
