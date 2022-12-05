# Constructor 


Use constructor for creating valid only object. If they cannot be created as valid only object (may because type can be initialize without calling constructor), then they should have a validate method instead.
https://stackoverflow.com/questions/43456801/best-practices-for-long-constructors-in-javascript


If you have validation logic that needs to be shared between creation and update, then place the logic in a separate function or method of the class.


## Example: New User 

Consider the implementation below, we have two validation logic here, one for `name` and another for `age`. During construction, we ensure that these fields are always valid. However, we did not do the same for update:

```go
var (
	ErrNameRequired = errors.New("name is required")
	ErrAgeIllegal   = errors.New("must be at least 13 years old")
)

type User struct {
	name string
	age  int
}

func NewUser(name string, age int) (*User, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if age < 13 {
		return nil, ErrAgeIllegal
	}
	return &User{
		name: name,
		age: age,
	}, nil
}
```

We can improve it by using custom setter and adding the validation there:
```go
var (
	ErrNameRequired = errors.New("name is required")
	ErrAgeIllegal   = errors.New("must be at least 13 years old")
)

type User struct {
	name string
	age  int
}

func NewUser(name string, age int) (*User, error) {
	u := &User{}
	if err := u.SetName(name); err != nil {
		return nil, err
	}
	if err := u.SetAge(age); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) SetName(name string) error {
	if name == "" {
		return ErrNameRequired
	}
	return nil
}
func (u *User) SetAge(age int) error {
	if age < 13 {
		return ErrAgeIllegal
	}
	return nil
}
```

This is also another preferred way, if you need a way to constantly validate if the user is valid.

```go
var (
	ErrNameRequired = errors.New("name is required")
	ErrAgeIllegal   = errors.New("must be at least 13 years old")
)

type User struct {
	name string
	age  int
}

func NewUser(name string, age int) (*User, error) {
	u := &User{
		name: name,
		age:  age,
	}
	if err := u.Validate(name); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) Validate() error {
	if u.name == "" {
		return ErrNameRequired
	}
	if age < 13 {
		return ErrAgeIllegal
	}
	return nil
}
```

Another disadvantage of having the `Validate` method is that it is easier to chain them if you have deeply nested association.
