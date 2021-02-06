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

# References

1. [Email uniqueness as an aggregate invariant](https://enterprisecraftsmanship.com/posts/email-uniqueness-as-aggregate-invariant/)
2. [Entity validation with visitors and extension methods](https://lostechies.com/jimmybogard/2007/10/24/entity-validation-with-visitors-and-extension-methods/)
