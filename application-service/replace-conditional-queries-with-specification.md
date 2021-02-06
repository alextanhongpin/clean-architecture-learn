# Replace Conditional Queries with Specification

As usual, any business usecase that related to access to repositories can make things hard for us.

Let's say you have a repository with the following interface:

```js
class UserRepository {
  findUser(id) {}
  findUserWithEmail(email) {}
  findUserWithConfirmationToken(token) {}
  findUserWithResetPasswordToken(token) {}
}
```

So what can be infer from this?
- The repository that queries the user entity are _explicit_ (which is good)
- The repository now have business logic (are filter logic business logic?) They can be part of the business rule, but the _domain service_ should still validate the entities returned by the repository. That is why the methods are called `findUserWithConfirmationToken`, and not `findUserWithValidConfirmationToken`. See below.
- There are other example queries, especially those dealing with a collection of entities, or the aggregated result (e.g., `findRegisteredUserCountsSinceDate`, `findUsersThatRegisteredAfter` etc)
- This could lead _explosion of methods_. However, this is much better than having a _generic repository_, e.g. `findBy(query)`. Also testing is much easier.
- At the moment, this is my go-to method.


# With Conditionals/enums

The idea is just to keep the public interface layer small, by exposing as little operations as possible. This however will introduce conditionals in a repository, which itself could be a bad idea. Also, this only works for a single field, but this is what we aim to achieve, by reducing the complexity of the repository query.
```js
const UserFields = {
  ID: 'id',
  EMAIL: 'email',
  CONFIRMATION_EMAIL_TOKEN: 'confirmation_email_token',
  RESET_PASSWORD_TOKEN: 'reset_password_token'
}

class UserRepository {
  findUser(value, key = UserFields.ID) {}
}

// Usage:
userRepository.find(id) // Defaults to findByID
userRepository.find(email, UserFields.EMAIL)
userRepository.find(email, UserFields.CONFIRMATION_EMAIL_TOKEN)
```

# With Specification Pattern
Can we simplify it to `findUser(UserSpecification)`?

## Example of Repository

Good:
- repository returns the entity
- domain service performs the business logic for checking token expiry, and throwing the error

```js
class ApplicationService {
  constructor(userRepository, userService) {
    this.userRepository = userRepository
    this.userService = userService
  }

  async confirmEmail({
    token
  }) {
    const user = await this.userRepository.findUserWithConfirmationToken(token)
    await this.userService.validateConfirmationTokenNotExpired(user)

    user.confirmEmail()
    await this.userRepository.updateEmail(user)
  }
}

const Time = {
  SECOND: 1000,
  MINUTE: 1000 * 60,
  HOUR: 1000 * 60 * 60,
  DAY: 1000 * 60 * 60 * 24
}

class UserService {
  constructor(confirmWith = 5 * Time.MINUTE) {
    this.confirmWithin = confirmWith
  }
  validateConfirmationTokenNotExpired(user) {
    if (Date.now() - user.confirmationSentAt.getTime() > this.confirmWith) {
      throw new ConfirmationTokenExpiredError('token expired')
    }
  }
}
```

Bad:
- the repository is performing a business logic
- errors are not obvious

```sql
class ApplicationService {
  constructor(userRepository, userService) {
    this.userRepository = userRepository
    this.userService = userService
  }

  async confirmEmail({
    token
  }) {
    const user = await this.userRepository.findUserWithValidConfirmationToken(token)

    user.confirmEmail()
    await this.userRepository.updateEmail(user)
  }
}

class UserRepository {
  constructor(db) {
    this.db = db
  }
  async findUserWithValidConfirmationToken(token) {
    const user = await this.db.query(`
			SELECT * 
			FROM user 
			WHERE token = $1 
			AND now() - confirmed_sent_at > interval '5 minute'
		`, token)
    // This will return empty if the query does not match, hence we 
    // are unable to discern if the user does not exists, or if the token expired.
    return user
  }
}
```


# References

1. [Applied Domain-Driven Design (DDD), Part 3 - Specification Pattern](http://www.zankavtaskin.com/2013/10/applied-domain-driven-design-ddd-part-3.html)
2. [SPECIFICATIONS, EXPRESSION TREES, AND NHIBERNATE](https://davefancher.com/2012/07/03/specifications-expression-trees-and-nhibernate/)
