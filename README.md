# clean-architecture-learn
Learnings from clean architecture


Domain layer
- Domain events
- Entities
- Services
- Repositories
- Value objects
- Aggregates

Application layer
- Application service, transaction happens here
- Repositories
- Message brokers, adapters
- Workers
- Authorization

Outermost
- controllers
- Authentication 
- Serialiser, logging, tracing metrics


## Find, execute, commit pattern
```js
// ApplicationService are also known as usecase layer. They do not contain business logic.
class ApplicationService {
  constructor(userRepository, userService) {
    this.userRepository = userRepository
    this.userService = userService
  }
  
  // Usecase to request confirmation email.
  async requestConfirmationEmail(email) {
    // 1. Repository: Find entity.
    const user = await this.userRepository.find(email)
    
    // 2. Domain service: Execute business logic.
    await this.userService.validateNotYetConfirmed(user) // Throws on error.
    
    // 3. Domain service: Update state of entity in-memory.
    const userWithConfirmationToken = await this.userService.createConfirmationToken(user)
    
    // 4. Repository: Persist entity state.
    const token = await this.userRepository.updateConfirmationToken(user)
    
    // Application service should not return entity. Either define a custom DTO, or return primitives.
    return token
  }
}
```

## References

1. [StackOverflow: UseCase-Drive vs Domain-Driven](https://stackoverflow.com/questions/3173070/design-methodology-use-case-driven-vs-domain-driven)
