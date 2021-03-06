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

# Good practices

- DDD implementation should not be influenced by persistence and database
- avoid anaemic domain model, domain classes full of getters and setters, void of behaviours.
- not all logic are domain logic, e.g. field validations
- security (e.g. authentication) is not part of the core domain (it could be a generic/supporting domain), however, most of the time, the implementation can lie in the `ui`/`application service` layer, out of the domain model.


Questions
- do i need to define a type for each layer, e.g. usecase.User, repository.User, entity.User? Nope, https://softwareengineering.stackexchange.com/questions/303478/uncle-bobs-clean-architecture-an-entity-model-class-for-each-layer
- why not user service, repository and entity in user package, vs separate repository for each? Circular dependency

# References

1. [StackOverflow: UseCase-Drive vs Domain-Driven](https://stackoverflow.com/questions/3173070/design-methodology-use-case-driven-vs-domain-driven)
