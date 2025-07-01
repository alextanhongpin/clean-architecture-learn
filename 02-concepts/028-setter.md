# Setter, getter and reflect

- How about combining reflect at runtime (once) and validating fields set?

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"

	"play.golang/structs"
)

type User struct {
	Name structs.Getter[string]
	Age  structs.Getter[int]
}

func main() {
	fmt.Println(structs.AllFieldsSet(User{
		Name: structs.NewField("john"),
		Age:  structs.NewField(10),
	}))
	fmt.Println("Hello, 世界")
}
-- go.mod --
module play.golang
-- structs/structs.go --
package structs

import (
	"errors"
	"fmt"
	"reflect"
)

type Changeable interface {
	Dirty() bool
}

type Setter[T any] interface {
	Set(T)
	Changeable
}

type Getter[T any] interface {
	Get() (T, bool)
	MustGet() T
	Changeable
}

var changeableType = reflect.TypeOf((*Changeable)(nil)).Elem()

func AllFieldsSet(in interface{}) bool {
	v := reflect.Indirect(reflect.ValueOf(in))
	if v.Kind() != reflect.Struct {
		return false
	}
	
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Type().Implements(changeableType) {
			if f.IsNil() {
				return false
			}
			responses := f.MethodByName("Dirty").Call([]reflect.Value{})
			if !responses[0].Bool() {
				return false
			}
		}

	}
	return true
}

type Field[T any] struct {
	value T
	dirty bool
}

func NewField[T any](t T) *Field[T] {
	return &Field[T]{
		value: t,
		dirty: true,
	}
}

func (f *Field[T]) Validate() error {
	if f == nil {
		return errors.New("not set")
	}
	return nil
}

func (f *Field[T]) MustGet() T {
	v, ok := f.Get()
	if !ok {
		panic("not set")
	}
	return v
}

func (f *Field[T]) Get() (T, bool) {
	if err := f.Validate(); err != nil {
		var t T
		return t, false
	}
	return f.value, f.dirty
}

func (f *Field[T]) Set(t T) {
	f.value = t
	f.dirty = true
}

func (f *Field[T]) Dirty() bool {
	return f.dirty
}
// Alternative
// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	var count int
	incr := func() {
		count++
	}
	fmt.Println(RetFunc(1, incr))
	fmt.Println(RetFunc("hello", incr))
	fmt.Println(count)
	
	req := NewRequired("name")
	fmt.Println(req)
	Require("name", "John", req.Set)
	fmt.Println(req, req.Valid())
}

func RetFunc[T any](t T, fn func()) T {
	fn()
	return t
}

func Require[T any](name string, t T, fn func(string)) T {
	fn(name)
	return t
}

type Required struct {
	fields map[string]bool
}

func NewRequired(field string, fields ...string) *Required {
	fields = append(fields, field)
	req := &Required{make(map[string]bool)}
	for _, f := range fields {
		req.Add(f)
	}
	return req
}
func (r *Required) Add(name string) {
	if _, ok := r.fields[name]; ok {
		panic("already add")
	}
	r.fields[name] = false
}

func (r *Required) Set(name string) {
	if ok, exists := r.fields[name]; ok {
		panic("already set")
	} else if !exists {
		panic("not exists")
	}
	r.fields[name] = true
}

func (r Required) Valid() bool {
	var count int
	for _, v := range r.fields {
		if v {
			count++
		}
	}
	return count == len(r.fields)
}
