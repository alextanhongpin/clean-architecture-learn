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
- The repository queries the user entity, but they have conditionals
- The repository know have business logic (?)
- There are other example queries, especially those dealing with a collection of entities, or the aggregated result (e.g., `findRegisteredUserCountsSinceDate`, `findUsersThatRegisteredAfter` etc)


# With Conditionals/enums

The idea is just to keep the public interface layer small, by exposing as little operations as possible. This however will introduce conditionals in a repository, which itself could be a bad idea.
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


# References

1. [Applied Domain-Driven Design (DDD), Part 3 - Specification Pattern](http://www.zankavtaskin.com/2013/10/applied-domain-driven-design-ddd-part-3.html)
2. [SPECIFICATIONS, EXPRESSION TREES, AND NHIBERNATE](https://davefancher.com/2012/07/03/specifications-expression-trees-and-nhibernate/)
