# Persistence Validation

Validations should ideally be part of _Entity_, and there are different variations to it:

During mutation:
1. update entity state locally first, then validate, before persisting
2. validate locally first before updating entity state, then persist
3. update at persistence layer directly and let the errors propagate

During creation:
1. Create in-memory first, then validate before persisting
2. Validate before create, ensuring the entity will always be in it's valid state. Typically done using factory/constructor.
3. Create and let the persistence layer warn on errors. Typically for uniqueness constraints (or check constraints in db)

During deletion:
1. Validate if can be deleted
2. Delete and let persistence layer throw error if fails

Above we notice some patterns
- the entity knows the state required
- the persistence layer knows the list of aggregate states. However, the errors thrown by persistence layer does not translate to domain errors.


How do we solve this?
