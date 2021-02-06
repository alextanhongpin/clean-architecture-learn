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

## Approach 2: Repository throws domain error

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
        throw new ErrEmailNotUnique('email exists')
      }
      // Do not propagate database errors, replace with other domain errors instead (below is a poor choice).
      throw new ErrInternalServerError('internal server error')
    }
  }
}
```

- However, now the details of the errors leaks to the repository layer. 
- Also, the implementation is very database specific, and now the user have to enforce checking for the database duplicate error and propagating them to the application service
- This example is fine, because it has only a single entry point (registration), and a single unique constraints. What if we have an entity with multiple unique constraints, and the validation logic needs to be repeated for updates? That is the main reason why business logic should be in the entity layer.

## Approach 3: Deferred validation

Think of it as an `after` hook. This is important, because the repo in the hooks are bounded by the same transaction.
```go
package main

import (
	"context"
	"errors"
	"fmt"
)

var ErrEmailExists = errors.New("email exists")

func main() {
	fmt.Println("Hello, playground")
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

func (u *User) ConfirmEmail() {
	u.email = u.unconfirmedEmail
	u.unconfirmedEmail = ""
}

func (u *User) RejectEmail() error {
	u.unconfirmedEmail = ""
	return ErrEmailExists
}

// Application Service

type userCreator interface {
	CreateUser(ctx context.Context, user *User) error
}

type UserUsecase struct {
	repo userCreator
	validator
}

func (u *UserUsecase) CreateUser(ctx context.Context, email string) error {
	usr := NewUnconfirmedUser(email)
	v := NewDeferredCreateUserValidator(u.repo)
	if err := v.Validate(ctx, u); err != nil {
		return err
	}
	// Once the creation is successful, we set the unconfirmed email to confirmed.
	usr.ConfirmEmail()
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
	if err := v.repo.CreateUser(ctx, u); err != nil {
		// Check if the error is part of the unique constraints error.
		// If yes, return it.
		return u.RejectEmail()
	}
	return nil
}
```

# References

1. [Email uniqueness as an aggregate invariant](https://enterprisecraftsmanship.com/posts/email-uniqueness-as-aggregate-invariant/)
2. [Entity validation with visitors and extension methods](https://lostechies.com/jimmybogard/2007/10/24/entity-validation-with-visitors-and-extension-methods/)
