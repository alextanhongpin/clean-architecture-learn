# Enforcing unique constraints

Context:
- we want to enforce unique email upon user registration. Which layer do we handle that?

Explanation:
- For uniqueness, we cannot perform them at the domain service layer - the problem with this is the information of the unique email is only known to the database
- So we will apply that constraints at the database layer, and handle the database error upon constraint violations
- So the logic should belong to the application service layer

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

## Update

The repository should contain the list of possible errors that will propagate out of the system, so the implementation above should be:

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
# References

1. [Email uniqueness as an aggregate invariant](https://enterprisecraftsmanship.com/posts/email-uniqueness-as-aggregate-invariant/)
2. [Entity validation with visitors and extension methods](https://lostechies.com/jimmybogard/2007/10/24/entity-validation-with-visitors-and-extension-methods/)
