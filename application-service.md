# Application service

- aka usecase level/service layer
- not to be confused with _domain service_ (aka business layer), application service does not contain business logic
- common flow: find an entity, execute domain service, persist entity through repository
- Is controller an application service? See [1]. Nope, in _clean architecture_ controller is primary adapter. The controller calls the application service.
- Application service is the glue code for domain services and repository, domain service performs business logic, repository persist them. 
- Other adapters such as logging, auth, message bus may be called here 
- The outer layer only knows the inner layer one layer below
- The inner layer does not know the outer layer
- responsible for executing transactions, see [3]
- if a user domain service has a method to `create` user, then the corresponsing application service will have a method `createUser`, see [2]
- The Service Layer is usually constructed in terms of discrete operations that have to be supported for a client. See [4]
- Knowing this, you may realize that your Business Layer really is a Service Layer as well. At some point, the point from which you're asking this question being one such point, the distinction is mostly semantic.
- deals with calling repository to persist an entity/aggregate


## Example of Application Service

Find, execute, commit pattern:
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

# References

1. [StackOverflow: Domain Service, Application Service](https://stackoverflow.com/questions/2268699/domain-driven-design-domain-service-application-service)
2. [Design of Service Layer and Application Logic](https://emacsway.github.io/en/service-layer/)
3. [Framework Design Guidelines: Domain Logic Patterns](https://www.informit.com/articles/article.aspx?p=1398617&seqNum=4)
4. [StackOverflow: Service Layer vs Business Layer in architecting web applications?](https://stackoverflow.com/questions/4108824/service-layer-vs-business-layer-in-architecting-web-applications#:~:text=The%20Service%20Layer%20is%20usually,objects%20to%20be%20persisted%2C%20etc.)
