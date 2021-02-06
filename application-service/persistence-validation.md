# Persistence Validation

Validations should ideally be part of _Entity_, and there are different variations to it:

During mutation:
1. update entity state locally first, then validate, before persisting
2. validate locally first before updating entity state, then persist
3. update at persistence layer directly and let the errors propagate

During creation:
1. Create in-memory first, then validate before persisting
2. Validate before create, ensuring the entity will always be in it's valid state. Typically done using factory/constructor.
3. Create and let the persistence layer warn on errors. Typically for uniqueness constraints (or check constraints in db)

During deletion:
1. Validate if can be deleted
2. Delete and let persistence layer throw error if fails

Above we notice some patterns
- the entity knows the state required
- the persistence layer knows the list of aggregate states. However, the errors thrown by persistence layer does not translate to domain errors.


How do we solve this?

## Persistence validation with visitor pattern

This example may be contrived, because uniqueness constraint are defined in the database. And after validation and before the user is inserted, another concurrent user with that email might have been inserted into the database. For unique email/slug (or other unique compound constraints), what we want is _deferred validation_, that is validating the errors once it is inserted into the database. See [2]. 

```go
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrEmailExists = errors.New("email exists")

type uniqueEmailValidator interface {
	Validate(email string) error
}

func main() {
	fmt.Println("Hello, playground")
}

// Entity.
type User struct {
	email string
}

func (u *User) ChangeEmail(email string, validator uniqueEmailValidator) error {
	if err := validator.Validate(email); err != nil {
		return err
	}
	u.Email = email
	return nil
}

// Repository interface.
type userFinder interface {
	FindUserByEmail(email string) (User, error)
}

// Validator.
type UniqueEmailValidator struct {
	repo userFinder
}

func NewUniqueEmailValidator(repo userRepository) *UniqueEmailValidator {
	return &UniqueEmailValidator{
		repo: repo,
	}
}

func (v *UniqueEmailValidator) Validate(email string) error {
	// User with email exists.
	_, err := v.repo.FindUserByEmail(email)
	if err == nil {
		return ErrEmailExists
	}

	// No such email.
	if errors.Is(sql.ErrNoRows, err) {
		return nil
	}

	return err
}

// Application Service

type userRepository interface {
	CreateUser(ctx context.Context, user *User) error
}

type UserUsecase struct {
	validator uniqueEmailValidator
	repo      userRepository
}

func (u *UserUsecase) CreateUser(ctx context.Context, email, password string) error {
	usr := new(User)
	if err := usr.ChangeEmail(email, u.validator); err != nil {
		return err
	}
	// Perform other task...
	// usr.SetPassword(password)
	return u.repo.CreateUser(ctx, usr)
}
```

# Reference

1. [Validation and DDD](https://enterprisecraftsmanship.com/posts/validation-and-ddd/)
2. [StackOverflow: How exactly should a CQRS Command be validated and transformed to a domain object?](https://softwareengineering.stackexchange.com/questions/348337/how-exactly-should-a-cqrs-command-be-validated-and-transformed-to-a-domain-objec)
