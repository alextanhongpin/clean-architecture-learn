## Example of domain layer with golang


```go
// You can edit this code!
// Click here and start typing.
package main

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"fmt"
)

var ErrNotFound = errors.New("not found")

func main() {
	p := Password("hello world")
	fmt.Println(string(p), p)

	creds, err := CreateUserCredentials("john", p)
	if err != nil {
		panic(err)
	}
	fmt.Println(creds.Login(p))

	var fakeCreds *UserCredentials
	fmt.Println(fakeCreds.Validate())
}

/*
create table users (
	id 			uuid default gen_random_uuid(),
	name 			text not null,
	email 			text not null,
	encrypted_password 	text not null,

	primary key (id),
	unique (email)
)

create table user_addresses (
	user_id uuid,
	address_id uuid,

	primary key (user_id, address_id)
)

create table addresses (
	id 		uuid default gen_random_uuid(),
	line_1 		text,
	line_2 		text,
	city 		text,
	state 		text,
	country 	text,
	postal_code 	text,

	primary key (id)
)
*/

type User struct {
	ID   string
	Name string

	// UserCredentials should not be part of the User entity - it is only for the process of authenticating user
	// Load only what is needed
	Address *Address // Value object
}

const MinPasswordLen = 8

var (
	ErrPasswordNotSet   = errors.New("password: not set")
	ErrPasswordTooShort = errors.New("password: too short")
)

// Value object has value receivers, since they are immutable.
// They are non-valid on creation, and has to be validated.
type Password string

func (p Password) String() string {
	if p == "" {
		return "**EMPTY PASSSWORD**"
	}
	return "**REDACTED PASSWORD**"
}

// Wrap similar operations together - if password needs to be validated before encrypting, place them together in a method.
func (p Password) Encrypt() (EncryptedPassword, error) {
	if err := p.Validate(); err != nil {
		return "", err
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(string(p)), 14)
	if err != nil {
		return "", err
	}
	return EncryptedPassword(string(bytes)), nil
}

func (p Password) Validate() error {
	n := len(p)
	switch {
	case n == 0:
		return ErrPasswordNotSet
	case n < MinPasswordLen:
		return ErrPasswordTooShort
	default:
		return nil
	}
}

var ErrEncryptedPasswordInvalid = errors.New("encrypted password: invalid")

type EncryptedPassword string

func (e EncryptedPassword) Decrypt(p Password) error {
	if err := p.Validate(); err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(string(e)), []byte(string(p)))
}

func (e EncryptedPassword) Validate() error {
	if len(e) == 0 {
		return ErrEncryptedPasswordInvalid
	}

	return nil
}

var UserCredentialsNotFound = fmt.Errorf("user credentials: %w", ErrNotFound)

type UserCredentials struct {
	Username          string
	EncryptedPassword EncryptedPassword
}

func (c *UserCredentials) Login(password Password) bool {
	err := c.EncryptedPassword.Decrypt(password)
	return err == nil
}

func (c *UserCredentials) Validate() error {
	// Treat every nil or non-initialized struct as not found error.
	if c == nil || (*c == UserCredentials{}) {
		return UserCredentialsNotFound
	}

	return nil
}

// Use Create- when preparing the data to insert into the persistence layer.
// If the construction requires other external layer, e.g. service layer, then use a factory instead.
func CreateUserCredentials(username string, password Password) (*UserCredentials, error) {
	enc, err := password.Encrypt()
	if err != nil {
		return nil, err
	}

	return &UserCredentials{
		EncryptedPassword: enc,
		Username:          username,
	}, nil
}

// Use New- when hydrating (aka recreating) the entity from the persistence layer.
// The difference between Create- and New- is New- is only responsible for setting the values/converting the
// persistence layer data into value object or domain objects, while Create- is responsible for transforming
// the data that is required for the initial persistence of the entity.
// In the example above, CreateUserCredentials has the responsibility to encrypt the plaintext password to encrypted password
// to be stored in the users table. However, when loading the data, the application only loads the encrypted password.
// The reason why UserCredentials is created as a value object is because on update/patch operations, we want to re-create the
// UserCredentials rather than treating it as updating a single field.
// In short, values that change together should best be represented as value object, since they are immutable to begin with.
// Also, we usually use New- constructor to create entity when the struct has private fields, since they cannot be set from other packages.
func NewUserCredentials() *UserCredentials {
	return nil
}

// Another concept is data mapping - so New- usually accepts a list of primitive args (string, int ...), but gets unwieldy when the number of fields is huge.
// Instead, we use mapper, which takes in a struct and returns another struct

type Address struct {
	Line1      string
	Line2      string
	City       string
	State      string
	PostalCode string
	Country    string
}
```
