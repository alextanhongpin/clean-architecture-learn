# Finder Pattern

As mentioned in `persistence-validation`, we learn to distinguish between persistence validation, and local entity validation. For persistence validation, the entity rule belongs to the repository, since the repository knows about the entity aggregates. The example below demonstrates how to validate the logic, and coupling them to the entity:

- we want to find a user with valid confirmation token, and confirm the email
- we initialize a blank user, and attempt to confirm the token
- find the user with a valid confirmation token (persistence validation)
- the user state is then rehydated
- validate the token has not yet expired (local validation)

Improvement:
- the current persistence validator accepts a User pointer. Should we be more specific in what fields it should validate? Would it lead to an explosion of interfaces?

```go
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

const ConfirmWithin = 5 * time.Minute

var (
	// Entity error.
	ErrUser = errors.New("User")

	// General.
	ErrUserNotFound = fmt.Errorf("%w: not found", ErrUser)

	// Field errors.
	ErrInvalidConfirmationToken = fmt.Errorf("%w: invalid confirmation token", ErrUser)

	// Field suberrors.
	ErrConfirmationTokenRequired = fmt.Errorf("%w: empty", ErrInvalidConfirmationToken)
	ErrConfirmationTokenExpired  = fmt.Errorf("%w: expired after %s", ErrInvalidConfirmationToken, ConfirmWithin)
)

func main() {
	repo := new(UserRepository)
	usecase := &UserUsecase{repo: repo}
	if err := usecase.ConfirmEmail(context.Background(), "valid"); err != nil {
		log.Println(err)
	}
	if err := usecase.ConfirmEmail(context.Background(), "invalid"); err != nil {
		log.Println(err)
	}
}

type User struct {
	email              string
	unconfirmedEmail   string
	confirmationSentAt *time.Time
	confirmationToken  string
	confirmedAt        *time.Time
}

type userPersistenceValidator interface {
	Validate(*User) error
}

func (u *User) ConfirmEmail(token string, v userPersistenceValidator) error {
	if token == "" {
		return ErrConfirmationTokenRequired
	}
	u.confirmationToken = token
	// This only validates that the user with the token exists, a.k.a persistence validation.
	// Responsibility of the token validity belongs to the domain entity.
	if err := v.Validate(u); err != nil {
		return err
	}
	
	// This is local validation and should be performed by the entity.
	if time.Since(*u.confirmationSentAt) > ConfirmWithin {
		return ErrConfirmationTokenExpired
	}
	now := time.Now()
	u.confirmationToken = ""
	u.confirmationSentAt = nil
	u.confirmedAt = &now
	u.email = u.unconfirmedEmail
	u.unconfirmedEmail = ""
	return nil
}

type ConfirmEmailValidator struct {
	ctx  context.Context // Request-bounded.
	repo userRepository
}

func NewConfirmEmailValidator(ctx context.Context, repo userRepository) *ConfirmEmailValidator {
	return &ConfirmEmailValidator{
		ctx:  ctx, // We don't want to expose context to User entity.
		repo: repo, // The repo may be transaction bound, hence we initialize the validator here.
	}
}

func (v *ConfirmEmailValidator) Validate(u *User) error {
	usr, err := v.repo.FindUserWithConfirmationToken(v.ctx, u.confirmationToken)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return ErrUserNotFound
		}
		return err
	}
	*u = usr
	return nil
}

type userRepository interface {
	FindUserWithConfirmationToken(ctx context.Context, token string) (User, error)
}

// ApplicationService
type UserUsecase struct {
	repo userRepository
}

func (uc *UserUsecase) ConfirmEmail(ctx context.Context, token string) error {
	// The user does not exists yet.
	u := new(User)
	v := NewConfirmEmailValidator(ctx, uc.repo)
	
	// Rehydrate the user if the confirmation token is valid.
	return u.ConfirmEmail(token, v)
}

// Repository.
type UserRepository struct{}

func (r *UserRepository) FindUserWithConfirmationToken(ctx context.Context, token string) (User, error) {
	if token == "valid" {
		now := time.Now().Add(-6 * time.Minute)
		return User{confirmationSentAt: &now}, nil
	}
	return User{}, sql.ErrNoRows
}

```
