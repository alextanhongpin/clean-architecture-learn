# Enforcing unique constraints

Context:
- we want to enforce unique email upon user registration. Which layer do we handle that?

Explanation:
- For uniqueness, we cannot perform them at the domain service layer - the problem with this is the information of the unique email is only known to the database. 
- Why not check if the email exists first? Simple, the system may not be consistent yet. Which means, we may validate that the email does not exists, but there is always a chance that the email may be inserted into the database by another concurrent request before the current one is completed. We can say that the validation is deferred until the actual user is actually inserted into the database.


# Approach 1: Application Service throws the error

- Since the domain service should not call the repository for persistence, and the entity should not know about the database error, so the logic should belong to the application service layer
- So we will apply that constraints at the database layer, and handle the database error upon constraint violations
- the repository throws an exception on duplicate email error, and the application service handles it. This form of delegation lies in the application service, which has an disadvantage when the same validation logic needs to be applied many times.
- shouldn't the error be part of the entity? is there no way to enforce it?

```js
class ApplicationService {
  constructor(userRepository, userFactory) {
    this.userRepository = userRepository
    this.userFactory = userFactory
  }

  // Enforcing unique email?
  register(email, password) {
    // Build the user entity, performs hashing of plaintext password to encrypted password.
    const user = await this.userFactory.createWithCredentials(email, password)

    // Throw exception on duplicate email.
    try {
      const newUser = await this.userRepository.create(user)
      return newUser
    } catch (error) {
      // Since the error message is not part of domain service, but database, we are handling this logic at the application service.
      if (isDuplicateError(error)) {
        throw new Error('email exists')
      }
      throw error
    }
  }
}
```

## Approach 2: Repository throws domain error (Preferred)

The repository should contain the list of possible errors that will propagate out of the system, so the implementation above could be:

```js
class ApplicationService {
  constructor(userRepository, userFactory) {
    this.userRepository = userRepository
    this.userFactory = userFactory
  }

  // Enforcing unique email?
  register(email, password) {
    // Build the user entity, performs hashing of plaintext password to encrypted password.
    const user = await this.userFactory.createWithCredentials(email, password)

    // The userRepository is now responsible for propagating domain errors.
    const newUser = await this.userRepository.create(user)
    return newUser
  }
}

class UserRepository {
  constructor(db) {
    this.db = db
  }

  async create({
    name,
    email,
    encryptedPassword
  }) {
    try {
      const user = await this.db.insert({
        name,
        email,
        encryptedPassword
      })
      return user
    } catch (error) {
      if (isDuplicateError(error)) {
        // throw new ErrEmailNotUnique('email exists')
	user.rejectEmail()
      }
      throw error
    }
  }
}
```

- However, now the details of the errors leaks to the repository layer. Not true, let the entity throw error.
- Also, the implementation is very database specific, and now the user have to enforce checking for the database duplicate error and propagating them to the application service. Update: its fine, it is better than leaking the database internals up to application service - it means no swapping db.
- This example is fine, because it has only a single entry point (registration), and a single unique constraints. What if we have an entity with multiple unique constraints, and the validation logic needs to be repeated for updates? That is the main reason why business logic should be in the entity layer. Update: yes, let the entity throw the error,
- In a way, this is actually the preffered method. it has less complexity. entity cannot call repository, but the opposite holds true. Also, the errors could be domain errors, if we handle the db specific errors at the application layer, we are already losing the possibility to swap database. 

## Approach 3: Deferred validation

Think of it as an `after` hook. This is important, because the repo in the hooks are bounded by the same transaction.
```go
package main

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrEmailExists              = errors.New("email exists")
	ErrUnconfirmedEmailRequired = errors.New("unconfirmed email is required")

	ErrDBUniqueConstraintViolation = errors.New("unique violation: email")
)

func main() {
	repo := new(UserRepository)
	usecase := &UserUsecase{repo: repo}
	fmt.Println(usecase.CreateUser(context.Background(), "john.doe@mail.com"))
	fmt.Println(usecase.CreateUser(context.Background(), "john"))
}

// Entity.
type User struct {
	email            string // Can't set before creation! That's violating the domain rule.
	unconfirmedEmail string // Use this instead.
}

func NewUnconfirmedUser(email string) *User {
	return &User{
		unconfirmedEmail: email,
	}
}

func (u *User) ConfirmEmail() error {
	if u.unconfirmedEmail == "" {
		return ErrUnconfirmedEmailRequired
	}
	u.email = u.unconfirmedEmail
	u.unconfirmedEmail = ""
	return nil
}

func (u *User) RejectEmail() error {
	u.unconfirmedEmail = ""
	return ErrEmailExists
}

// Application Service

type userCreator interface {
	CreateUser(ctx context.Context, user *User) error
}

type deferredCreateUserValidator interface {
	Validate(ctx context.Context, u *User) error
}

type UserUsecase struct {
	repo userCreator
}

func (u *UserUsecase) CreateUser(ctx context.Context, email string) error {
	usr := NewUnconfirmedUser(email)

	// Why not place this in the struct? Because the operation might be transaction-bounded.
	v := NewDeferredCreateUserValidator(u.repo)
	if err := v.Validate(ctx, usr); err != nil {
		return err
	}

	// Once the creation is successful, we set the unconfirmed email to confirmed.
	return usr.ConfirmEmail()
}

// This is another layer in between the repository and application service.
// Unlike domain services, this can be stateful, but more of a decorator to the existing repository.
// You can think of it as an "after" hook too.
type DeferredCreateUserValidator struct {
	repo userCreator
}

func NewDeferredCreateUserValidator(repo userCreator) *DeferredCreateUserValidator {
	return &DeferredCreateUserValidator{repo: repo}
}

// Validate validates the creation of the user with unique email - this is
// because the uniqueness constraints are set by the database.
func (v *DeferredCreateUserValidator) Validate(ctx context.Context, u *User) error {
	if err := v.repo.CreateUser(ctx, u); err != nil {
		// Check if the error is part of the unique constraints error.
		// If yes, return it.
		if errors.Is(ErrDBUniqueConstraintViolation, err) {
			return u.RejectEmail()
		}
		return err
	}
	return nil
}

type UserRepository struct{}

func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
	if user.unconfirmedEmail == "john.doe@mail.com" {
		return ErrDBUniqueConstraintViolation
	}
	return nil
}
```

Similar implementation for JavaScript:
```js
class User {
  #unconfirmedEmail = ''
  #email = ''

  set unconfirmedEmail(value) {
    if (!value) {
      throw new Error("setUnconfirmedEmailError: unconfirmed email cannot be empty")
    }
    this.#unconfirmedEmail = value
  }

  confirmEmail() {
    if (!this.#unconfirmedEmail) {
      throw new Error("no email to confirm")
    }
    this.#email = this.#unconfirmedEmail
    this.#unconfirmedEmail = ''
  }

  rejectEmail() {
    this.#unconfirmedEmail = ''
    throw new Error('email exists')
  }
}

class ApplicationService {
  constructor(userRepository, userFactory) {
    this.userRepository = userRepository
    this.userFactory = userFactory
  }

  register(email, password) {
    // Factory creates the user entity with unconfirmed email, as well as performing hashing of plaintext password to encrypted password.
    const user = await this.userFactory.createUnconfirmedUserWithCredentials(email, password)

    const validator = new DeferredCreateUserValidator(this.userRepository)
    await validator.validate(user)

    user.confirmEmail()
  }
}

class DeferredCreateUserValidator {
  constructor(userRepository) {
    this.userRepository = userRepository
  }
  // validate throws exception on duplicate email.
  async validate(user) {
    try {
      const newUser = await this.userRepository.create(user)
      return newUser
    } catch (error) {
      // The user is responsible for throwing domain errors.
      if (isDuplicateError(error)) {
        user.rejectEmail()
      }
      throw error
    }
  }
}
```

This is unnecessary complexity.

Improvement by ensuring the `confirmEmail` or `rejectEmail` is called:

```go
type confirmUserCreation interface {
	RejectEmail() error
	ConfirmEmail() error
}

// Validate validates the creation of the user with unique email - this is
// because the uniqueness constraints are set by the database.
func (v *DeferredCreateUserValidator) Validate(ctx context.Context, u confirmUserCreation) error {
	user, ok := u.(*User)
	if !ok {
		return errors.New("invalid user")
	}
	if err := v.repo.CreateUser(ctx, user); err != nil {
		// Check if the error is part of the unique constraints error.
		// If yes, return it.
		if errors.Is(ErrDBUniqueConstraintViolation, err) {
			return u.RejectEmail()
		}
		return err
	}
	return u.ConfirmEmail()
}
```

Using `visitor` to enforce the entity to perform the validation:
```go
package main

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrEmailExists              = errors.New("email exists")
	ErrUnconfirmedEmailRequired = errors.New("unconfirmed email is required")

	ErrDBUniqueConstraintViolation = errors.New("unique violation: email")
)

func main() {
	repo := new(UserRepository)
	usecase := &UserUsecase{repo: repo}
	fmt.Println(usecase.CreateUser(context.Background(), "john.doe@mail.com"))
	fmt.Println(usecase.CreateUser(context.Background(), "john"))
}

// Entity.
type User struct {
	email            string // Can't set before creation! That's violating the domain rule.
	unconfirmedEmail string // Use this instead.
}

type userCreationValidator interface {
	Validate(ctx context.Context, u *User) error
}

func (u *User) PlaceEmail(email string) error {
	if u.email != "" {
		return ErrEmailExists
	}
	u.unconfirmedEmail = email
	return nil
}

func (u *User) RejectEmail() error {
	u.unconfirmedEmail = ""
	return ErrEmailExists
}

// ConfirmEmailPlacement verifies that the email has been claimed upon successful creation.
func (u *User) ConfirmEmailPlacement(ctx context.Context, v userCreationValidator) error {
	if err := v.Validate(ctx, u); err != nil {
		// Just to ensure the contract is fulfilled.
		if errors.Is(ErrEmailExists, err) {
			// This would be called in the validator too.
			return u.RejectEmail()
		}
		return err
	}

	if u.unconfirmedEmail == "" {
		return ErrUnconfirmedEmailRequired
	}
	u.email = u.unconfirmedEmail
	u.unconfirmedEmail = ""
	return nil
}

// Application Service

type userCreator interface {
	CreateUser(ctx context.Context, user *User) error
}

type deferredCreateUserValidator interface {
	Validate(ctx context.Context, u *User) error
}

type UserUsecase struct {
	repo userCreator
}

func (u *UserUsecase) CreateUser(ctx context.Context, email string) error {
	usr := new(User)
	if err := usr.PlaceEmail(email); err != nil {
		return err
	}

	// Why not place this in the struct? Because the operation might be transaction-bounded.
	v := NewDeferredCreateUserValidator(u.repo)
	if err := usr.ConfirmEmailPlacement(ctx, v); err != nil {
		return err
	}

	return nil
}

// This is another layer in between the repository and application service.
// Unlike domain services, this can be stateful, but more of a decorator to the existing repository.
// You can think of it as an "after" hook too.
type DeferredCreateUserValidator struct {
	repo userCreator
}

func NewDeferredCreateUserValidator(repo userCreator) *DeferredCreateUserValidator {
	return &DeferredCreateUserValidator{repo: repo}
}

// Validate validates the creation of the user with unique email - this is
// because the uniqueness constraints are set by the database.
func (v *DeferredCreateUserValidator) Validate(ctx context.Context, u *User) error {
	err := v.repo.CreateUser(ctx, u)

	// Check if the error is part of the unique constraints error.
	// If yes, return it.
	if errors.Is(ErrDBUniqueConstraintViolation, err) {
		return u.RejectEmail()
	}
	return err
}

type UserRepository struct{}

func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
	if user.unconfirmedEmail == "john.doe@mail.com" {
		return ErrDBUniqueConstraintViolation
	}
	return nil
}
```

# Thoughts

- why not handle them as events, `EmailRejectedEvent`?

# References

1. [Email uniqueness as an aggregate invariant](https://enterprisecraftsmanship.com/posts/email-uniqueness-as-aggregate-invariant/)
2. [Entity validation with visitors and extension methods](https://lostechies.com/jimmybogard/2007/10/24/entity-validation-with-visitors-and-extension-methods/)
