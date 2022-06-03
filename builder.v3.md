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


## Another variation, handles ignored fields

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type User struct {
	Name    string
	Age     int `construct:"-"`
	Married bool
	hobbies []string // Works with private fields
}

func NewUser(name string, age int, married bool, hobbies []string) *User {
	return &User{
		Name:    name,
		Age:     age,
		Married: married,
		hobbies: hobbies,
	}
}

var UserConstructor = ConstructorFactory(User{})

func main() {
	ctor := UserConstructor()
	u := User{
		Name:    Set("Name", "john", ctor),
		Age:     Set("Age", 10, ctor),
		Married: Set("Married", false, ctor),
		hobbies: Set("hobbies", []string{}, ctor),
	}
	if err := ctor.Validate(); err != nil {
		panic(err)
	}
	fmt.Println(u)

	ctor = UserConstructor()
	u2 := NewUser(
		Set("Name", "john", ctor),
		Set("Age", 10, ctor),
		Set("Married", false, ctor),
		Set("hobbies", []string{}, ctor),
	)
	if err := ctor.Validate(); err != nil {
		panic(err)
	}
	fmt.Println(u2)
}

type constructor interface {
	Set(name string) bool
}

func Set[T constructor, K any](name string, k K, t T) K {
	if !t.Set(name) {
		panic(fmt.Errorf("%w: %s", ErrUnknownField, name))
	}
	return k
}

func ConstructorFactory[T any](unk T) func() *Constructor[T] {
	t := reflect.Indirect(reflect.ValueOf(unk)).Type()

	fields := make(map[string]int, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Tag.Get("construct") == "-" {
			fields[f.Name] = -i
		} else {
			fields[f.Name] = i
		}
	}

	return func() *Constructor[T] {
		return NewConstructor[T](fields)
	}
}

var (
	ErrUnsetFields  = errors.New("constructor: unset fields")
	ErrUnknownField = errors.New("constructor: field not found")
)

type Constructor[T any] struct {
	set    int
	fields map[string]int
}

func NewConstructor[T any](fields map[string]int) *Constructor[T] {
	return &Constructor[T]{fields: fields}
}

func (c *Constructor[T]) Set(name string) bool {
	i, ok := c.fields[name]
	if !ok {
		return false
	}
	if i < 0 {
		return true
	}

	bit := 1 << i
	if c.isSet(bit) {
		return false
	}

	c.set |= bit
	return true
}

func (c *Constructor[T]) isSet(bit int) bool {
	return c.set&bit == bit
}

func (c *Constructor[T]) Validate() error {
	unsetFields := make([]string, 0, len(c.fields))
	for f, i := range c.fields {
		if i < 0 {
			continue
		}

		bit := 1 << i
		if c.isSet(bit) {
			continue
		}

		unsetFields = append(unsetFields, f)
	}

	if len(unsetFields) > 0 {
		return fmt.Errorf("%w: %s", ErrUnsetFields, strings.Join(unsetFields, ", "))
	}

	return nil
}
```


## Conclusion

After benchmarking, the difference between using this and just setting the fields individually shows huge performance difference.

So, avoid the approach above, write tests instead to check if the values are set, and use [go-cmp](https://github.com/google/go-cmp) instead. To check against bool, always set it to true.
