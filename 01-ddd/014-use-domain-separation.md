# Use domain separation

How to differentiate between domain model, domain service, and domain aggregate?

They are very similar in many ways, and sometimes it is unclear where to place the domain logic.


## Domain Service

- It is easier to separate domain service from domain model and domain aggregate - if you have any external dependencies, use domain service.
- External dependencies can be side-effects such as random number generator and datetime. 
- Any other IO may be passed (though should be limited) such as HTTP service calls, but they should only be passed as interface abstraction rather than concrete implementation. 
- Avoid passing db calls, that belongs to the repository layer.

## Domain Model

- Domain model is represents a single entity loaded from the database. 
- The domain model should not have any associations, except references to those association (e.g. foreign key).
- Domain model should contain methods that does not accept any arguments. Computed values could only be derived from the existing fields.


## Domain Aggregate

- Domain aggregate represents a concept that involves multiple domain models
- The rules above applies, the domain aggregate should not have access to external dependencies
