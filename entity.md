# Entity

- Typically modelled from the database table
- Entity is just dto between repository and app, so it should not bubble up the whole layer.
- Entity does not have business logic.
- Entity should contain the primary business logic
- Business logic not in entity can be placed at service layer (validation not for entity layer, validation for fields, and also logic like encrypting password where there may be multiple strategy argon vs bcrypt can be done at service layer)
- Repositories/domain services operates on entity
- Application service does not return entities, expose a separate dto instead.
- While not a must, entity's fields are usually private, and may only be modified through setters/factories/constructors to ensure they are set to valid states
- instead of setter for individual fields, use setters for a group of fields to validate them as a business logic, `changePersonalDetails` (name, age etc) instead of setters for each.


# Responsibilities

## Validation

Entity are responsible for validation. Entities should always be in a valid state. There are however some usecases where the validation does not make sense to put in the entity.

For example:
- when the validation rules requires access to repository, e.g. enforcing uniqueness of email
- when the validation should operate on a collection of entity, e.g. a user can only have 5 accounts

In any of the cases above, the entity is limited by the fact that it should not have access to the repository layer. If the domain service cannot access repository too, then the only possible layer to place this logic is in the application service itself.

## References

1. [StackOverflow: Validation in Domain Model of Domain Service](https://stackoverflow.com/questions/35934713/validation-in-domain-model-of-domain-service)
2. [StackOverflow: Using Repository in Entity for validation before update](https://stackoverflow.com/questions/55549616/ddd-using-repository-in-entity-for-validation-before-update)
