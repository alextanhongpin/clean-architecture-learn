# Repository


> the layer that requests external data and returns a domain entity 


external data is not necessarily related to persistence. it can be external api call, cache or storage object. repository connect multiple storage and client layers and returns domain entity.

## Why Repository

- We want to separate the implementation of the persistence layer, and only expose the interfaces.
- The goal is not to help make switching database implementation easier. But it does makes switching ORMs easier (e.g. in golang, go-pg has been superseeded by bun). 

## About repository
- there are two parts, the _interface_ and also the _implementation_ layer. The interface only defines the contracts (input/output), and the implementation is storage specific. In hexagonal architecture, _interface_ is the _port_, _implementation_ is the _adapter_.
- what is the input and output of repository? The input/output is the domain entity, but can also be basic primitives.
- are repositories entity specific? Nope, design entity for an aggregate, see [1].
- can a repository have business logic? Nope :smile: 
- are filter logic and/or conditional logic in repository business logic? Yes and no. You can treat them as domain business logic or aggregate logic, e,g fetching total issues count in a bug tracking application, or checking uniqueness of email could be aggregate logic
- are repositories implemented as class or functions? Ideally class, but usually it is for namespacing purposes, and also allows an abstract class/interface to be defined alongside. 
- is partial update allowed in repository. Usually the whole entity is fetched first prior to update, and then the fields are updated before saved. However, for performance reason (and also simplicity, since we know what field will be modified), partial updates should be allowed. This is either done by having partial dtos at the repository layer for updates 


on partial updates
https://stackoverflow.com/questions/19407390/domain-driven-design-how-to-handle-updates-for-parts-of-your-aggregrate-roots

https://enterprisecraftsmanship.com/posts/partially-initialized-entities-anti-pattern/

what other alternatives for partial update? perhaps we can diff the changes into a generic data structure, and then only update those fields that changed. However, this might not detect fields tha are removed.

for repository, we may also want to distinguish between update put and patch. cause for null columns, there might be times we want to set them back to null.

## About the layer
- where does the repository _interface_ layer belongs to? In the domain layer, a.k.a (domain > repository) (See [3])
- where does the repository _implementation_ layer belongs to? Ideally between the __application service__ and __domain model__. Some place it in the __persistence layer__. 
- can a repository method updates multiple entities? Repository are aggregate-specific (see [1]), and may update only an aggregate and the child. A method should only update one of them though (?)
- are repository meant for a single entity? Nope, it is meant for an aggregate root. An aggregate root may have several children entities (which can be nested). This adds some complication, cause when fetching, we want to fetch the aggregate root as a whole to protect against invariance. However, when updating, we must also consider if we want to return the whole object, or just partial. When using Postgres, we cannot utilise the `returning *` after an update/insert for example, because it does not return the child entities. When updating, we may also want to consider if we want to only apply the update or fetch the whole aggregate root before calling the setter method to update. There's also complexity with bulk operations, as rebuilding the whole entity may require a lot of queries to the database.
- what is the difference between _repository_ and _factory_? A __factory__ handles the beginning of an object’s life; a __repository__ helps manage the middle and the end. (See [2])
- who call the repository? Application services, domain services, factories (I have not seen any example, but ideally the repository should be called at application service layer, and the entities be passed to factory instead)
- can entity, value object call repository? Nope 
- can domain service call repository? Ideally shouldn't, but there are cases that is mentioned in [4]
- where do we start the transaction? At the application service
- why shouldn't we start the transaction at this layer? 
   - Simple, if we do this, we cannot chain it with other repository methods
   - It makes testing harder, it is harder to override the transaction or will lead to nested transaction
- do i need repository for the reason of enabling swapping or db or mocking of db? Nope, both is not a strong reason to have repository layer because you need to test your queries against your external store/remote data anyway (not necessarily sql). But repository do hide the complexity of the query and make it reusable. 
- repository has its own dto model, which will map back to the domain model. the mapping is done in the repository layer. https://softwareengineering.stackexchange.com/questions/404076/domain-vs-entities-model-domain-driven-design-ddd
- should repository return domain models? Yes

## About persistence
- should one repository map to one database table? Nope, you should never create a repository for each table in the database. (See [1]).

## DB Partial Update

A lot of solution uses DataMapper pattern - mapping only non-empty fields from the DTO to domain Entity.
https://www.baeldung.com/spring-data-partial-update
https://stackoverflow.com/questions/15329436/partial-updates-for-entities-with-repository-dto-patterns-in-mvc-prepping-for-a
https://softwareengineering.stackexchange.com/questions/389542/updating-the-db-in-the-repository-from-a-dto-in-a-layered-architecture
https://auth0.com/blog/automatically-mapping-dto-to-entity-on-spring-boot-apis/
http://codebetter.com/iancooper/2011/04/12/repository-saveupdate-is-a-smell/

Most repository pattern assumes the following pseudo code:

```js
class User {}

class UserRespository {
   async save(user: User) {
      // The save does not take into consideration whether the User is new or partially updated.
   }
}
```

There are a few problems here:
- the save does not take into consideration whether the fields are filled or partial, so there may be overriding of fields with value to empty 
- one way to mitigate this is to load the entire entity, perform updates and save the entity back. This handles the limitation above, but at the price of performance. 
- a better way is to perform partial updates, where only a subset of fields are updated
- it is also possible to diff the fields with values that has changed, and pass them to the ORM
- the idea above is the same as event sourcing, except that with event sourcing, the events are defined manually


## Should I map repository types to domain types?

There's a lot of work mapping from one layer to another, especially when they are almost a 1-to-1 representation. However, when they are not, it is always beneficial to separate repository types from domain types. 

There are some scenarios that would eventually lead to that:
- custom columns in db layer that cannot be mapped directly to domain types, e.g. go's sqlc might generate multiple types depending on the query, like `GetUserWithNameRows`, `GetUserWithEmailRows` that has same fields but different types that needs to be mapped back to User entity
- composition of multiple entities that are fetched from joins instead (ORMs are not efficient at separating them, so fetching say `user` joined to `accounts` requires custom mapping to separate them if the tool you use does not separate the `User` and `Account` type

# How to fetch aggregate

Instead of fetching the data you need individually with many small repository method, handle them in a single method and return the aggregate instead. This simplifies the usecase layer.

# Stores

The repository ahould just be a facade that connects multiple sources.

- repository is not equal to ORM or raw SQL queries, though it can be
- sometimes it is easier to abstract orm and table specific queries into stores
- stores are essentially implementation details, and one repository can have multiple stores
- an example is an AuthRepository, which may call both the UserStore and AccountStore


# References

1. [Microsoft DDD CQRS Pattern: Infrastructure Persistence Layer Design](https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/infrastructure-persistence-layer-design#:~:text=Repositories%20are%20classes%20or%20components,required%20to%20access%20data%20sources.&text=Conceptually%2C%20a%20repository%20encapsulates%20a,closer%20to%20the%20persistence%20layer.)
2. [StackOverflow: DDD Repository and Factory](https://stackoverflow.com/questions/31528368/ddd-repository-and-factory)
3. [StackOverflow: Which layer do DDD repositories belong to?](https://softwareengineering.stackexchange.com/questions/396151/which-layer-do-ddd-repositories-belong-to)
4. [Google Groups: Can Domain Service Access Repositories](https://groups.google.com/g/dddcqrs/c/66zbcL97ilk?pli=1)
