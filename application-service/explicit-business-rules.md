# Strong rules

I used to make business rules dynamic, even when they are not needed. Here's an example of a user email confirmation service:

```js
// Enums are value objects.
const Time = {
  SECOND: 1000,
  MINUTE: 1000 * 60,
  HOUR: 1000 * 60 * 60,
  DAY: 1000 * 60 * 60 * 24
}

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

The integration test:

```js
const repository = new UserRepository(db)
const service = new UserService(0) // Hurray, I can test the scenario that the token expired.
const usecase = new ApplicationService(repository, service)
usecase.confirmEmail({ token })
```

Here's a list of mistakes I made:
- the domain service takes a `confirmWithin` parameter, the reason is that I want to mock the behaviour later (we will see later why it's wrong)
- that makes the domain service _stateful_. Domain services should be constructable without factories, we can identify a bad domain service if they are not pure functions. The reason for creating them as class methods is purely for namespacing.
- The integratin test is mocking the `confirmWithin` duration... in a sense, it's not integration testing anymore. Putting the duration to 0 is actually validating the wrong behaviour.
- The business rule is not explicit enough, in fact, it's uncertain.

Here's how to correct the implementation:

```js
// Explicit business rule. This is not hardcoding.
const CONFIRM_WITHIN = 5 * Time.MINUTE

class UserService {
  // Constructor is no longer required.
  validateConfirmationTokenNotExpired(user) {
    if (Date.now() - user.confirmationSentAt.getTime() > CONFIRM_WITHIN) {
      throw new ConfirmationTokenExpiredError('token expired after %s', CONFIRM_WITHIN) // Explicit rules makes explicit errors.
    }
  }
}


const repository = new UserRepository(db)
const service = new UserService()
const usecase = new ApplicationService(repository, service)

// Turn back time. Now we are testing actual behaviour!
repository.overrideConfirmationSentAt(-6 * Time.MINUTE)
usecase.confirmEmail({ token })
```
