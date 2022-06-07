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
package main_test

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

var _ = os.Setenv("DEBUG_CTOR", "true")

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

var user User

func BenchmarkCtor(b *testing.B) {
	var u User
	for i := 0; i < b.N; i++ {
		ctor := UserConstructor()
		u = User{
			Name:    Set("Name", "john", ctor),
			Age:     Set("Age", 10, ctor),
			Married: Set("Married", false, ctor),
			hobbies: Set("hobbies", []string{}, ctor),
		}
		if err := ctor.Validate(); err != nil {
			panic(err)
		}
	}
	user = u
}

var user2 User

func BenchmarkCtorNullObject(b *testing.B) {
	var u User
	for i := 0; i < b.N; i++ {
		var ctor *Constructor[User]
		u = User{
			Name:    Set("Name", "john", ctor),
			Age:     Set("Age", 10, ctor),
			Married: Set("Married", false, ctor),
			hobbies: Set("hobbies", []string{}, ctor),
		}
		if err := ctor.Validate(); err != nil {
			panic(err)
		}
	}
	user2 = u
}

var user3 User

func BenchmarkNormal(b *testing.B) {
	var u User
	for i := 0; i < b.N; i++ {
		u = User{
			Name:    "john",
			Age:     10,
			Married: false,
			hobbies: []string{},
		}
	}
	user3 = u
}

var user4 User

func BenchmarkConstructor(b *testing.B) {
	var u *User
	for i := 0; i < b.N; i++ {
		u = NewUser("john", 10, false, []string{})
	}

	user4 = *u
}

var user5 User

func BenchmarkConstructorCtor(b *testing.B) {
	var u *User
	for i := 0; i < b.N; i++ {
		ctor := UserConstructor()
		u = NewUser(
			Set("Name", "john", ctor),
			Set("Age", 10, ctor),
			Set("Married", false, ctor),
			Set("hobbies", []string{}, ctor),
		)
		if err := ctor.Validate(); err != nil {
			panic(err)
		}
	}

	user5 = *u
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

func envBool(name string) bool {
	ok, _ := strconv.ParseBool(os.Getenv(name))
	return ok
}

func ConstructorFactory[T any](unk T) func() *Constructor[T] {
	debug := envBool("DEBUG_CTOR")
	if !debug {
		return func() *Constructor[T] {
			return nil
		}
	}

	t := reflect.Indirect(reflect.ValueOf(unk)).Type()

	var ignore int
	fields := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {

		f := t.Field(i)
		if f.Tag.Get("construct") == "-" {
			ignore |= 1 << i
		}

		fields[i] = f.Name
	}

	return func() *Constructor[T] {
		return NewConstructor[T](fields, ignore)
	}
}

var (
	ErrUnsetFields  = errors.New("constructor: unset fields")
	ErrUnknownField = errors.New("constructor: field not found")
)

type Constructor[T any] struct {
	fields []string
	set    int
	ignore int
}

func NewConstructor[T any](fields []string, ignore int) *Constructor[T] {
	return &Constructor[T]{
		fields: fields,
		ignore: ignore,
	}
}

func (c *Constructor[T]) Set(name string) bool {
	if c == nil {
		return true
	}

	i, ok := c.indexOf(name)
	if !ok {
		return false
	}

	bit := 1 << i
	if c.isSet(bit) || c.isIgnore(bit) {
		return c.isIgnore(bit)
	}

	c.set |= bit
	return true
}

func (c *Constructor[T]) Validate() error {
	if c == nil {
		return nil
	}

	unsetFields := make([]string, 0, len(c.fields))
	for i, f := range c.fields {
		bit := 1 << i
		if c.isSet(bit) || c.isIgnore(bit) {
			continue
		}

		unsetFields = append(unsetFields, f)
	}

	if len(unsetFields) > 0 {
		return fmt.Errorf("%w: %s", ErrUnsetFields, strings.Join(unsetFields, ", "))
	}

	return nil
}

func (c *Constructor[T]) isSet(bit int) bool {
	return c.set&bit == bit
}

func (c *Constructor[T]) isIgnore(bit int) bool {
	return c.ignore&bit == bit
}

func (c *Constructor[T]) indexOf(name string) (int, bool) {
	for i, f := range c.fields {
		if f == name {
			return i, true
		}
	}

	return 0, false
}
```

Benchmark result:

```
goos: darwin
goarch: amd64
pkg: github.com/alextanhongpin/ctor
cpu: Intel(R) Core(TM) i5-6267U CPU @ 2.90GHz
BenchmarkCtor-4                  5465547               217.1 ns/op
BenchmarkCtorNullObject-4       45910035                24.57 ns/op
BenchmarkNormal-4               412362188                2.946 ns/op
BenchmarkConstructor-4          24110187                48.74 ns/op
BenchmarkConstructorCtor-4       4541757               263.9 ns/op
PASS
ok      github.com/alextanhongpin/ctor  7.205s
```

## Conclusion

After benchmarking, the difference between using this and just setting the fields individually shows huge performance difference.

So, avoid the approach above, write tests instead to check if the values are set, and use [go-cmp](https://github.com/google/go-cmp) instead. To check against bool, always set it to true.
