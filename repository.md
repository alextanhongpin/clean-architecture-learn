# Repository

## About repository
- there are two parts, the _interface_ and also the _implementation_ layer. The interface only defines the contracts (input/output), and the implementation is storage specific. In hexagonal architecture, _interface_ is the _port_, _implementation_ is the _adapter_.
- what is the input and output of repository? The input/output is the domain entity, but can also be basic primitives.
- are repositories entity specific? Nope, design entity for an aggregate, see [1].
- can a repository have business logic? Nope :smile:
- are repositories implemented as class or functions? Ideally class, but usually it is for namespacing purposes, and also allows an abstract class/interface to be defined alongside. 


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

## About persistence
- should one repository map to one database table? Nope, you should never create a repository for each table in the database. (See [1]).

# References

1. [Microsoft DDD CQRS Pattern: Infrastructure Persistence Layer Design](https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/infrastructure-persistence-layer-design#:~:text=Repositories%20are%20classes%20or%20components,required%20to%20access%20data%20sources.&text=Conceptually%2C%20a%20repository%20encapsulates%20a,closer%20to%20the%20persistence%20layer.)
2. [StackOverflow: DDD Repository and Factory](https://stackoverflow.com/questions/31528368/ddd-repository-and-factory)
3. [StackOverflow: Which layer do DDD repositories belong to?](https://softwareengineering.stackexchange.com/questions/396151/which-layer-do-ddd-repositories-belong-to)
4. [Google Groups: Can Domain Service Access Repositories](https://groups.google.com/g/dddcqrs/c/66zbcL97ilk?pli=1)
