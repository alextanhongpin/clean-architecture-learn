# Constructor


Use constructor for creating valid only object. If they cannot be created as valid only object (may because type can be initialize without calling constructor), then they should have a validate method instead.


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

## Use setter to protect against invariant

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


## Constructor should not return error

In language like JavaScript, we can easily throw an error to protect against invariant:

```typescript
class User {
	constructor(private name: string, private email: string) {
		if (!name) throw new Error('UserError: name is required')
		if (!email) throw new Error('UserError: email is required')
	}
}
```

However, in Golang, the idiomatic way is to return error:

```go
type User struct {
	name string
	email string
}

func NewUser(name, email string) (*User, error) {
	u := &User{
		name: name,
		email: email,
	}
	if err := u.Validate(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) Validate() error {
	if u.name == "" {
		return errors.New("user: name is required")
	}
	if u.email == "" {
		return errors.New("user: email is required")
	}

	return nil
}
```


Let's think about the scenario where we may have an aggregate with nested entities that needs to be constructed:


```go
type User struct {
	name string
	email string
	address *Address
}

type Address struct {
	line1 string
	line2 string
	country string
	city string
	state string
	postalCode string
}
```

First of all, we can update the `User`'s `Validate` method to take the address validation into account:

```go
func (u *User) Validate() error {
	if u.name == "" {
		return errors.New("user: name is required")
	}
	if u.email == "" {
		return errors.New("user: email is required")
	}

	if u.address != nil {
		return u.address.Validate()
	}

	return nil
}
```

Instead of returning an error when constructing new user, we just return the pointer to the user. It is up to the end user to call the `Validate` method:

```go
func NewUser(name, email string, address *Address) *User {
	return &User{
		name: name,
		email: email,
		address: address,
	}
}
```


## Long constructors

Say if we have 80 attributes for a given entity, the constructor will have to take 80 fields. This may become problematic in a codebase when you have a lot of construction being called. Making change to the attribute requires changing all of them, and potentially causing errors if the order is incorrect.


## Private fields

See [entity](entity.md).



https://stackoverflow.com/questions/43456801/best-practices-for-long-constructors-in-javascript
