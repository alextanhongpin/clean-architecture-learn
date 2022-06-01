# Using builder with generics

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"reflect"
	"strings"
)

type User struct {
	name    string
	age     int
	married bool `builder:"-"`
	hobbies []string
}

var UserBuilder = StructBuilder(&User{})

func main() {
	b := UserBuilder()

	u := User{
		name:    Set("name", "john", b),
		age:     Set("age", 10, b),
		hobbies: Set("hobbies", []string{"cycling"}, b),
	}
	fmt.Println(u)
	fmt.Println(b)
	fmt.Println(b.Error())
}

func StructBuilder[T any](t T) func() *Builder[T] {
	typ := reflect.Indirect(reflect.ValueOf(t)).Type()
	fields := make([]string, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Tag.Get("builder") == "-" {
			continue
		}
		fields[i] = f.Name
	}

	return func() *Builder[T] {
		return NewBuilder[T](fields)
	}
}

type Builder[T any] struct {
	fields []string
	set    int
}

func NewBuilder[T any](fields []string) *Builder[T] {
	if len(fields) == 0 {
		panic("builder: no fields provided in constructor")
	}
	return &Builder[T]{
		fields: fields,
	}
}

func (b *Builder[T]) index(name string) (idx int, found bool) {
	for i, field := range b.fields {
		if field == "" {
			continue
		}
		if field == name {
			idx = i
			found = true
			break
		}
	}
	return
}

func (b *Builder[T]) Set(name string) bool {
	i, found := b.index(name)
	if !found {
		return false
	}

	n := 1 << i
	if b.set&n == n {
		return false
	}

	b.set |= n
	return true
}

func (b *Builder[T]) Error() error {
	fields := make([]string, 0, len(b.fields))
	for i, field := range b.fields {
		if field == "" {
			continue
		}
		n := 1 << i
		if b.set&n != n {
			fields = append(fields, field)
		}
	}
	if len(fields) > 0 {
		return fmt.Errorf("builder: %s not set", strings.Join(fields, ", "))
	}
	return nil
}

func Set[T any](name string, t T, setter interface{ Set(name string) bool }) T {
	if !setter.Set(name) {
		panic(fmt.Errorf("field not set: %s", name))
	}
	return t
}
```
