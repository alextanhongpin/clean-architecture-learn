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
