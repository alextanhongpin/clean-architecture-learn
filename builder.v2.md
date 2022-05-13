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

var (
	TypeString reflect.Type = reflect.TypeOf("")
	TypeInt    reflect.Type = reflect.TypeOf(0)
	TypeFloat  reflect.Type = reflect.TypeOf(0.0)
	TypeBool   reflect.Type = reflect.TypeOf(false)
)

type User struct {
	Name    Name
	Age     Age
	Hobbies []Hobby
	Height  float32
	Address *Address
	Skills  []string
	Married bool
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

// Declare it once.
var NewUserBuilder = NewBuilder(User{})

func main() {
	// Cons: Only work with public fields at the moment, due to json unmarshaling skipping private fields.

	user := NewUserBuilder().
		SetString("Name", "John").
		SetInt("Age", 10).
		SetBool("Married", true).
		SetFloat("Height", 167.0).
		Set("Hobbies", []string{"dancing"}).
		Set("Address", (*Address)(nil)). // This is different from new(Address).
		Set("Skills", []string{"JavaScript"}).
		Build()

	fmt.Println("user:", user)

	user2 := NewUserBuilder().SetMap(
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
	name   string
	values map[string]interface{}
	meta   map[string]reflect.Type
}

func NewBuilder[T any](t T) func() *Builder[T] {
	meta := make(map[string]reflect.Type)

	val := reflect.Indirect(reflect.ValueOf(t))
	for i := 0; i < val.NumField(); i++ {
		f := val.Type().Field(i)
		meta[f.Name] = f.Type
	}
	name := reflect.TypeOf(t).Name()
	return func() *Builder[T] {
		return &Builder[T]{
			name:   name,
			values: make(map[string]interface{}),
			meta:   meta,
		}
	}
}

func (b *Builder[T]) SetMap(kv map[string]interface{}) *Builder[T] {
	for k, v := range kv {
		b.Set(k, v)
	}
	return b
}

func (b *Builder[T]) isMatchingType(src, tgt reflect.Type) bool {
	return src == tgt || tgt.ConvertibleTo(src)
}

func (b *Builder[T]) setterType(key string) reflect.Type {
	t, ok := b.meta[key]
	if !ok {
		panic(fmt.Errorf("key does not exist: %s", key))
	}
	if _, isSet := b.values[key]; isSet {
		panic(fmt.Errorf("key has been set: %s", key))
	}
	return t
}

func (b *Builder[T]) SetString(key, value string) *Builder[T] {
	t := b.setterType(key)
	if !b.isMatchingType(t, TypeString) {
		panic(fmt.Errorf("type does not match: expected %s, got %s", t, TypeString))
	}

	b.values[key] = value
	return b
}

func (b *Builder[T]) SetInt(key string, value int) *Builder[T] {
	t := b.setterType(key)
	if !b.isMatchingType(t, TypeInt) {
		panic(fmt.Errorf("type does not match: expected %s, got %s", t, TypeInt))
	}

	b.values[key] = value
	return b
}

func (b *Builder[T]) SetFloat(key string, value float64) *Builder[T] {
	t := b.setterType(key)
	if !b.isMatchingType(t, TypeFloat) {
		panic(fmt.Errorf("type does not match: expected %s, got %s", t, TypeFloat))
	}

	b.values[key] = value
	return b
}

func (b *Builder[T]) SetBool(key string, value bool) *Builder[T] {
	t := b.setterType(key)
	if !b.isMatchingType(t, TypeBool) {
		panic(fmt.Errorf("type does not match: expected %s, got %s", t, TypeBool))
	}

	b.values[key] = value
	return b
}

func (b *Builder[T]) Set(key string, value interface{}) *Builder[T] {
	t := b.setterType(key)
	tt := reflect.TypeOf(value)

	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		t = t.Elem()
		tt = tt.Elem()
	}

	if !b.isMatchingType(t, tt) {
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
		panic(fmt.Errorf("%s fields not set: %s", b.name, strings.Join(fieldsNotSet, ", ")))
	}

	return b.BuildPartial()
}
```
