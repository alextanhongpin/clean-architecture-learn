# Entity

## Introduction

An entity represents a domain object with a distinct identity that persists throughout its lifecycle. Entities are the building blocks of domain models and encapsulate both data and behavior relevant to the business domain.

## Key Characteristics

- **Identity**: Each entity has a unique identifier that distinguishes it from other entities
- **Lifecycle**: Entities have a lifecycle from creation to deletion
- **Behavior**: Entities encapsulate business logic and domain rules
- **State**: Entities maintain internal state that can change over time

## Design Principles

### Database Independence
- Typically modelled from the database table, but entity implementation should be database-free
- In `golang`, there shouldn't be a `sql.NullString` in the entity fields, or a mongo bson tag in the field tags
- Entity is just dto between repository and app, so it should not bubble up the whole layer

### Business Logic Encapsulation
- Entity should contain the primary business logic
- Entity encapsulates behaviour required by the (sub)domain, making it natural, expressive, understandable
- Business logic not in entity can be placed at service layer (validation not for entity layer, validation for fields, and also logic like encrypting password where there may be multiple strategy argon vs bcrypt can be done at service layer)

### Layer Interactions
- Repositories/domain services operates on entity
- Application service does not return entities, expose a separate dto instead
- Entity return errors
- Does not perform IO
- Does not call third party libs, define interface and leave the implementation as infrastructure level and call them at domain service/application service/repository, whichever makes sense

## Encapsulation Guidelines

### Field Access
- While not a must, entity's fields are usually private, and may only be modified through setters/factories/constructors to ensure they are set to valid states
- Instead of setter for individual fields, use setters for a group of fields to validate them as a business logic, `changePersonalDetails` (name, age etc) instead of setters for each

### Getters and Setters
- Entity can have getters
- Entity getters should be immutable (for list, make a copy first)
- Entity should not have individual setters field, we usually modify a set of value objects (e.g. updateUserDetail, updateAddress vs setName, setAge)
- Entity should not be modifiable outside the domain
- Unrestrictive setters makes it easy to violate domain rules, set a group of value objects and always check for invariants

## Validation Responsibilities

Entity are responsible for validation. Entities should always be in a valid state. There are however some usecases where the validation does not make sense to put in the entity.

### Entity-Level Validation
- Field validation and business rules that can be determined from the entity's own state
- Invariant enforcement within the entity boundaries

### External Validation (Not Entity Responsibility)
- When the validation rules requires access to repository, e.g. enforcing uniqueness of email
- When the validation should operate on a collection of entity, e.g. a user can only have 5 accounts

In any of the cases above, the entity is limited by the fact that it should not have access to the repository layer. If the domain service cannot access repository too, then the only possible layer to place this logic is in the application service itself.

## Construction Patterns

### Constructor Design
There are several options for constructor, but ideally we want one that can scale with growing fields:

- **Minimal constructors**: Constructors should have minimal required fields
- **Builder pattern**: Use required constructor field + functional optional for handling default args
- **Valid state**: Entity should be constructed in its valid state
- **Validation**: On build for builder pattern, pass in a validator to validate for different constructor args

### Value Objects
- Group related entity fields as value objects, e.g. Address, Confirmable
- Separate construction by using builder pattern
- Constructors can evolve to mega constructor, so plan accordingly

## Entity vs Domain Model

Depending on your understanding, there may be no difference between domain model and domain entity. However, for certain applications, it is better to split the domain model and entity:

- **Entity**: Data that mimics the database table
- **Domain Model**: May be similar to or represent a subset of the entity
- **Mapping Considerations**: 
  - In some languages or when using ORMs, the domain entity may contain database specific mappings (e.g. golang's gorm or bun), however domain model should not have those details
  - Domain entity may also be composed purely of domain primitives, but when mapping to domain model, they should be mapped to value objects when possible
  - The opposite is true, domain entity should not contain domain models or value object

## Best Practices

1. **Single Responsibility**: Each entity should have a single, well-defined responsibility
2. **Invariant Protection**: Always maintain entity invariants through proper encapsulation
3. **Side-effect Isolation**: Side-effects (date time, random number) should be modelled externally for testability
4. **Reusability**: If the domain logic does not rely on the identity of the domain model, separate it into value object or service layer, to promote reusability
5. **Method Design**: Domain methods should not contain accepts external input - the computed values should be based solely on itself. Delegate such operations to domain services instead.

## References

1. [Validation in Domain Model of Domain Service](https://stackoverflow.com/questions/35934713/validation-in-domain-model-of-domain-service)
2. [Using Repository in Entity for validation before update](https://stackoverflow.com/questions/55549616/ddd-using-repository-in-entity-for-validation-before-update)
3. [Quasiclass](http://www.idinews.com/quasiClass.pdf)
