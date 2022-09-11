# Use value object

Value object does not have identity.

Value object is used as an alternative to primitive obsession.

## Value object implementation

```go
package main

import (
	"errors"
	"fmt"
	"strings"
)

func main() {
	var e Email
	fmt.Println(e.Validate())
}

var ErrEmailBadFormat = errors.New("email: bad format")

type Email string

func MakeEmail(email string) Email {
	return Email(email)
}

func (e Email) String() string {
	return string(e)
}

func (e Email) Valid() bool {
	return strings.Contains(e.String(), "@")
}

func (e Email) Validate() error {
	if !e.Valid() {
		return ErrEmailBadFormat
	}

	return nil
}
```
