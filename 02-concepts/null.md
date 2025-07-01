How to deal with null values in DDD?

How can values be null? They have not been created yet, especially for aggregate root entities.

Use null object pattern. https://blog.ndepend.com/null-evil/


Alternatively, we can introduce a valid method (r.g. isZero for golang), but this defeats the purpose of always valid object. Entity should be created in its valid state. Alternatively, just set it to null and return a Boolean.


## Null 

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"play.ground/nulls"
)

type T struct {
	*time.Time
}
type User struct {
	Name     string                 `json:"name"`
	Birthday *nulls.Null[time.Time] `json:"birthday"`
}

func main() {
	var u User
	u.Name = "john"
	u.Birthday = nulls.NewNull(time.Now(), false)
	b, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))

	var u2 User
	if err := json.Unmarshal(b, &u2); err != nil {
		panic(err)
	}
	fmt.Println(u2.Birthday.Valid())
	fmt.Println(u2.Birthday.Value().String())
}
-- go.mod --
module play.ground
-- nulls/nulls.go --
package nulls

import (
	"bytes"
	"encoding/json"
	"errors"
)

var ErrNotSet = errors.New("not set")

type Null[T any] struct {
	value       T
	valid       bool
	constructed bool
}

func NewNull[T any](t T, valid bool) *Null[T] {
	return &Null[T]{
		value:       t,
		valid:       valid,
		constructed: true,
	}
}

func (n *Null[T]) Validate() error {
	if n == nil || !n.constructed {
		return ErrNotSet
	}
	return nil
}

func (n *Null[T]) Value() T {
	if err := n.Validate(); err != nil {
		var t T
		return t
	}
	return n.value
}

func (n *Null[T]) Valid() bool {
	if err := n.Validate(); err != nil {
		return false
	}
	return n.valid
}
func (n *Null[T]) MarshalJSON() ([]byte, error) {
	if n.Valid() {
		return json.Marshal(n.value)
	}
	return []byte("null"), nil
}

func (n *Null[T]) UnmarshalJSON(raw []byte) error {
	if bytes.Equal(raw, []byte("null")) {
		return nil
	}
	if err := json.Unmarshal(raw, &n.value); err != nil {
		return err
	}
	n.constructed = true
	n.valid = true
	return nil
}
```
