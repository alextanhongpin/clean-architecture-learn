# clean-architecture-learn
Learnings from clean architecture


Domain layer (innermost)
- Domain events
- Entities
- Services
- Repositories
- Value objects
- Aggregates

Application layer (middle)
- Application service, transaction happens here
- Repositories
- Message brokers, adapters
- Workers
- Authorization

Presentation (outermost)
- controllers (REST) or resolver (Graphql)
- Authentication 
- Serialiser, logging, tracing metrics

# Good practices

- DDD implementation should not be influenced by persistence and database
- avoid anaemic domain model, domain classes full of getters and setters, void of behaviours.
- not all logic are domain logic, e.g. field validations
- security (e.g. authentication) is not part of the core domain (it could be a generic/supporting domain), however, most of the time, the implementation can lie in the `ui`/`application service` layer, out of the domain model.


Questions
- do i need to define a type for each layer, e.g. usecase.User, repository.User, entity.User? Nope, https://softwareengineering.stackexchange.com/questions/303478/uncle-bobs-clean-architecture-an-entity-model-class-for-each-layer
- why not user service, repository and entity in user package, vs separate repository for each? Circular dependency


What is a Model?
- a representation of your domain object
- contains rich behaviour
- Contains reference to other model
- Attributes may be primitive or value object

What is an Entity?
- a subset of your model (may or may not mimic your model)
- The persisted view of your model
- The model is always translatable to an entity, and vice versa
- Attributes may be db specific (but not necessarily)


What is a repository?
- fetches model and aggregate from a remote store (external request, remote store, client)

What is an Aggregate?
- composed of one or more model
- For business logic that spans across multiple entities
- Can be returned from a repository
- Why not service? Service may not require the aggregate to be called together (user may forget)
(Show an example of an aggregate and non aggregate)
(Show a good aggregate vs bad aggregate)



## Complex Application

- https://www.alibabacloud.com/blog/how-to-code-complex-applications-core-java-technology-and-architecture_595506?spm=a2c65.11461447.0.0.2ae134ffVGA3P9


# References

[^1]: [StackOverflow: UseCase-Drive vs Domain-Driven](https://stackoverflow.com/questions/3173070/design-methodology-use-case-driven-vs-domain-driven)

where to put domain services
- https://stackoverflow.com/questions/48200345/ddd-where-should-logic-go-that-tests-the-existence-of-an-entity
- https://medium.com/geekculture/an-introduction-to-domain-driven-design-f29fc1877e2a
- https://stackoverflow.com/questions/48200345/ddd-where-should-logic-go-that-tests-the-existence-of-an-entity
