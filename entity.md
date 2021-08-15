# Entity

- Typically modelled from the database table, but entity implementation should be database-free (e.g. in `golang`, there shouldn't be a `sql.NullString` in the entity fields, or a mongo bson tag in the field tags)
- Entity is just dto between repository and app, so it should not bubble up the whole layer.
- Entity encapsulates behaviour required by the (sub)domain, making it natural, expressive, understandable
- Entity should contain the primary business logic
- Business logic not in entity can be placed at service layer (validation not for entity layer, validation for fields, and also logic like encrypting password where there may be multiple strategy argon vs bcrypt can be done at service layer)
- Repositories/domain services operates on entity
- Application service does not return entities, expose a separate dto instead.
- While not a must, entity's fields are usually private, and may only be modified through setters/factories/constructors to ensure they are set to valid states
- instead of setter for individual fields, use setters for a group of fields to validate them as a business logic, `changePersonalDetails` (name, age etc) instead of setters for each.
- entity return errors
- does not perform IO
- does not call third part libs, define interface and leave the implementation as infrastructure level and call them at domain service/application service/repository, whichever makes sense



# Responsibilities

## Validation

Entity are responsible for validation. Entities should always be in a valid state. There are however some usecases where the validation does not make sense to put in the entity.

For example:
- when the validation rules requires access to repository, e.g. enforcing uniqueness of email
- when the validation should operate on a collection of entity, e.g. a user can only have 5 accounts

In any of the cases above, the entity is limited by the fact that it should not have access to the repository layer. If the domain service cannot access repository too, then the only possible layer to place this logic is in the application service itself.

## Design thoughts

There are several options for constructor, but ideally we want one that can scale with growing fields.

- constructors can evolve to mega constructor
- constructors should have minimal required fields
- use required constructor field + functional optional for handling default args

- entity can have getters
- entity getters should be immutable (for list, make a copy first)
- entity should not have individual setters field, we usually modify a set of value objects (e.g. updateUserDetail, updateAddress vs setName, setAge)
- entity should not be modifiable outside the domain
- unrestrictive setters makes it easy to violate domain rules, set a group of value objects and always check for invariants
- entity should be constructed in it's valid state
- separate construction by using builder pattern
- on build for builder pattern, pass in a validator to validate for different constructor args
- group related entity fields as value objects, e.g. Address, Confirmable

# Golang example


## Make properties private
To protect against direct modification, which could lead to inconsistent state, e.g. setting name to empty string.
```go
type User struct {
	name        string
	email       string
	hobbies     []string
	phoneNumber *string
}
```

## Use getters and return copies

```go
func (u User) Name() string {
	return u.name
}
```

For slices, always return the copy of the slice.
```go
func (u User) Hobbies() []string {
	out := make([]string, len(u.hobbies))
	copy(out, u.hobbies)
	return out
}
```

## Use withers instead of setters to avoid mutating the data

Withers should protect against invariant.
```go
func (u User) WithName(name string) (User, error) {
	if len(name) == "" {
		return u, errors.New("name cannot be empty")
	}
	u.name = name
	return u, nil
}
```

## Use pointers to indicate optional field
```go
type User struct {
	/*REDACTED*/
	phoneNumber *string
}
// Setters should not be pointers. So if you want to set a pointer value, pass the value and a field valid to indicate if the value should be set.
func (u User) WithPhoneNumber(phoneNumber string, exists bool) (User, error) {
	if !exists {
		return u, nil
	}
	if !checkPhoneNumberValid(phoneNumber) {
		return u, errors.New("invalid phone number")	
	}
	u.phoneNumber = phoneNumber
	return u, nil
}
```

## Use builder to set fields from repository layer

Setters on domain model shoyld onlybe used to modify entity fields without violating invariants. You dont expose setters for every field. However, when loading model from repository (remember, repository return domain model), you may need to set all fields for the model, whether valid or invalid.

However, due to the private fields, you will end up with huge constructorvwhen setting the fields, or you might need to expose alternative and every setters that are invariant free (or by checking if value is valid before setting). The solution is to just separate the construction of the domain model to the builder.

## References

1. [StackOverflow: Validation in Domain Model of Domain Service](https://stackoverflow.com/questions/35934713/validation-in-domain-model-of-domain-service)
2. [StackOverflow: Using Repository in Entity for validation before update](https://stackoverflow.com/questions/55549616/ddd-using-repository-in-entity-for-validation-before-update)
3. [Quasiclass](http://www.idinews.com/quasiClass.pdf)
