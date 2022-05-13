// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	// introduce setter getter struct field, and reflect to check all fields are set.
	var s Setter[int]
	s.Set(1)

	m := map[string]interface{}{
		"name": Setter[string]{},
		"age":  Setter[int]{},
	}
	m["age"].Set(1)
	fmt.Println("Hello, 世界", m, s)
}

// DeferSet
// SetBeforeGet
// nullType with generics
// Apply setter getter to struct
// How about combining reflect at runtime (once)
// And validating fields set?

type Setter[T any] struct {
	value T
	dirty bool
}

func (s *Setter[T]) Set(t T) {
	s.value = t
	s.dirty = true
}


-- set.go --
// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"reflect"
)

type S struct{}

func (S) Valid() bool {
	return true
}
func main() {
	var s S
	InvokeValid(s, "Valid")
	fmt.Println("Hello, 世界")
}
func InvokeValid(s any, name string) {
	t := reflect.Indirect(reflect.ValueOf(s))
	v := t.MethodByName(name).Call([]reflect.Value{})
	return v[0].Bool()
}
