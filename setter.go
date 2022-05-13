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

type Setter[T any] struct {
	value T
	dirty bool
}

func (s *Setter[T]) Set(t T) {
	s.value = t
	s.dirty = true
}
