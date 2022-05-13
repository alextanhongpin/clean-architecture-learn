## Builder Pattern with go generics 1.18

```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type User struct {
	Name    Name
	Age     Age
	Hobbies []Hobby
	Address *Address
	Skills  []string
}

type Name string

type Age int

type Hobby string

type Address struct {
	Line1    string
	Line2    string
	City     string
	State    string
	Postcode string
	Country  string
}

func main() {
	// Cons: Only work with public fields at the moment, due to json unmarshaling skipping private fields.
	builder := NewBuilder(User{})
	user := builder.
		Set("Name", "john").
		Set("Age", 10).
		Set("Hobbies", []string{"dancing"}).
		Set("Address", (*Address)(nil)). // This is different from new(Address).
		Set("Skills", []string{"JavaScript"}).
		Build()

	fmt.Println("user:", user)

	user2 := NewBuilder(User{}).SetMap(
		map[string]interface{}{
			"Name":    "Jane",
			"Age":     100,
			"Hobbies": []string{"walking"},
			// "Address": new(Address),
			"Skills": []string{"Python"},
		},
	).BuildPartial()
	fmt.Println("user2:", user2)
}

type Builder[T any] struct {
	values map[string]interface{}
	meta   map[string]reflect.Type
}

func NewBuilder[T any](t T) *Builder[T] {
	meta := make(map[string]reflect.Type)

	val := reflect.Indirect(reflect.ValueOf(t))
	for i := 0; i < val.NumField(); i++ {
		f := val.Type().Field(i)
		meta[f.Name] = f.Type
	}

	return &Builder[T]{
		values: make(map[string]interface{}),
		meta:   meta,
	}
}

func (b *Builder[T]) SetMap(kv map[string]interface{}) *Builder[T] {
	for k, v := range kv {
		b.Set(k, v)
	}
	return b
}

func (b *Builder[T]) Set(key string, value interface{}) *Builder[T] {
	t, ok := b.meta[key]
	if !ok {
		panic(fmt.Errorf("key does not exist: %s", key))
	}
	if _, isSet := b.values[key]; isSet {
		panic(fmt.Errorf("key has been set: %s", key))
	}

	tt := reflect.TypeOf(value)

	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		t = t.Elem()
		tt = tt.Elem()
	}
	valid := t == tt || tt.ConvertibleTo(t)
	if !valid {
		panic(fmt.Errorf("type does not match: expected %s, got %s", t, tt))
	}

	b.values[key] = value
	return b
}

// BuildPartial build with sane defaults.
func (b *Builder[T]) BuildPartial() T {
	var t T
	by, err := json.Marshal(b.values)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(by, &t); err != nil {
		panic(err)
	}
	return t
}

// Build requires all fields to be set, otherwise it will panic. To build with sane defaults, use BuildPartial.
func (b *Builder[T]) Build() T {
	if len(b.meta) != len(b.values) {
		var fieldsNotSet []string
		for k := range b.meta {
			if _, ok := b.values[k]; !ok {
				fieldsNotSet = append(fieldsNotSet, k)

			}
		}
		panic(fmt.Errorf("fields not set: %s", strings.Join(fieldsNotSet, ", ")))
	}

	return b.BuildPartial()
}
```
