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
- the place where dependencies are resolved
- interface to your domain - used by external consumers to talk to your system
- all logic concerning application workflow lies here
- it’s method is a use case, single flow
- handles the flow of usecases, including any additional concerns needed on top of the domains
- should generally have simple flow, complex application service flow indicates that domain logic has leaked out the domain
- coordinates application flow and infrastructure, but do not execute business logic rules or invariants
- it is common to see calls to repositories, unit of works (database transactions), message bus, cache
- initialize and oversees interaction between the domain objects and services
    - get domain object(s) from repository/check object exists
    - execute an action
    - put them back in repository/or not
- sits above domain model and coordinates application activity
- does not contain business logic, does not hold state for any entities (however it can store the state of a business workflow transaction)
- validate and save entity only in application service through repository
- performs applicative level logic as user interaction, input validation, logic not related to business but to other concerns (authentication, security, emailing)
- accepts and returns service contract objects or request/response objects, e.g. dto. So if the domain has a _User_ entity, then a simplified _UserDto_ can be returned. The adapter/transformer/converter to convert between entity and dto can lie in the application service layer.
- does not accept or return domain entities or value objects. The reason for this is that we do not want to expose the domain models in the outer layer. Any operations by the entity can only happen in the application service and the layers below. The output should be immutable, and that is why DTO is often recommended as the toutput value. Some people design application services as command handlers, which does not return anything, which makes me wonder how the testing is done.
- note, an alternative design is to return DPO, aka domain payload object https://vaadin.com/blog/ddd-part-3-domain-driven-design-and-the-hexagonal-architecture
- the above is mostly for anemic data model (or immutable)
- can an application service update multiple entity at the same time? Only if they belong to the same aggregate root. What if we need to update an unrelated entity? We do so by publishing an event.
- anstraction hides implementation, so avoid abstracting things you want to make it clear. Abstract things like server implementation/database transaction setup, but dont abtract things related to business logic.
- usecase accepts a request, and returns a response. Although it should not return domain objects, most of the time, it is simpler as it reduces mapping. Suitable if your domain model is anemic with only getters but no setters. An alternative is to separate mutation domain logic from the entity (read only or compute methods is fine) and put it in service layer (some language like julia dont have classes and the struct cant have methods). This works well with go especially when you want to keep the domain fields public since without it, you will end up with super long constructor to init the private fields.
- conclusion: return getter domain objects. move the mutation to service layer to reduce types conversion. this reduces the chance of the presentatioj layer from invoking methods from the domain layer that mutates its state. 
- treat usecase errors as part of the response. So have a dedicated error per usecase. Example of login usecase will have LoginError that could be union of UserNotFound or WrongPassword.

## Validation in Application Service


- _Input validation_ should be done at presentation layer, not application service.
- Application service should only return __ONE__ domain error to represent the error, not a list of validation input errors.
- Input validation refers to checking if a field exists, or is non-empty, or required etc.
- Application service should throw an exception when any of the params are invalid

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

## Event sourcing where to place event handler

https://softwareengineering.stackexchange.com/questions/325996/ddd-where-to-place-domain-event-handlers

https://softwareengineering.stackexchange.com/questions/168481/how-to-choose-between-using-a-domain-event-or-letting-the-application-layer-orc

## DTO assembler

https://buildplease.com/pages/repositories-dto/


## Avoiding Business logic in applications services


How to detect signs of business logic leaking in the application service layer?
- conditionals when handling operations (suggestion: refactor to strategy pattern, or polymorphism)
- usecase calling another usecase to orchestrate a flow (suggestion: use domain events, or outbox pattern)


## Mastery checklist 

- basic crud
- calling repositories 
- using unit of work
- validation 
- testing
- sending task to background 
- sending task to queue
- concurrency
# References

1. [StackOverflow: Domain Service, Application Service](https://stackoverflow.com/questions/2268699/domain-driven-design-domain-service-application-service)
2. [Design of Service Layer and Application Logic](https://emacsway.github.io/en/service-layer/)
3. [Framework Design Guidelines: Domain Logic Patterns](https://www.informit.com/articles/article.aspx?p=1398617&seqNum=4)
4. [StackOverflow: Service Layer vs Business Layer in architecting web applications?](https://stackoverflow.com/questions/4108824/service-layer-vs-business-layer-in-architecting-web-applications#:~:text=The%20Service%20Layer%20is%20usually,objects%20to%20be%20persisted%2C%20etc.)
5. [Application Services - 10 common doubts answered](https://blog.arkency.com/application-service-ruby-rails-ddd/)
