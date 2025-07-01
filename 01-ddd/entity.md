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

Another option is to just live with immutable domain models or use init only setter https://docs.microsoft.com/en-us/dotnet/csharp/language-reference/keywords/init#:~:text=init%20(C%23%20Reference)&text=An%20init%2Donly%20setter%20assigns,%2DImplemented%20Properties%2C%20and%20Indexers.

## Public/Private fields

Should fields be set to private with getter/setter?


There were some thoughts where all fields in a class should be private, and they can only be set through constructor or setters with validation applied to ensure the class is always valid.


However, this approach usually ends up with a lot of code. Imagine a hypothetical scenario where a class have 80 fields - that would mean that the constructor needs to accept 80 args, as well as having 80 getters and potentially 80 setters.

Instead of doing that, we can just keep the fields public, and keep the setters still. The fields can only be set using setters (this is up to PR reviews to ensure no direct setters are done), and the code is pretty much more maintainable in a way that one write less code.

https://stackoverflow.com/questions/35832379/oop-private-field-or-private-property-setter-in-regards-to-ddd


## Guard against invalid entity

To guard against invalid entity, just validate them at the beginning of each layer:

```go
package main

import (
	"errors"
	"fmt"
)

func main() {
	fmt.Println("Hello, 世界")
	u := &User{}
	saveUser(NewValidate(u))
}

func saveUser(val *Validate[*User]) {
	fmt.Println(val.Value())
}

type validatable interface {
	Validate() error
}

// Validate is a container type that ensures that the wrapped object fulfills the interface.
type Validate[T validatable] struct {
	val T
}

func NewValidate[T validatable](t T) *Validate[T] {
	return &Validate[T]{val: t}
}

func (v *Validate[T]) Value() T {
	if err := v.val.Validate(); err != nil {
		panic(err)
	}

	return v.val
}

type User struct {
	Name string
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("user: name is required")
	}

	return nil
}
```

We can write less code this way:

```go
package main

import (
	"errors"
	"fmt"
)

func main() {
	fmt.Println("Hello, 世界")
	u := &User{}
	saveUser(u)
}

func saveUser(u *User) {
	MustValidate(u)

	fmt.Println(u)
}

type validatable interface {
	Validate() error
}

func MustValidate(v validatable) {
	if err := v.Validate(); err != nil {
		panic(err)
	}
}

type User struct {
	Name string
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("user: name is required")
	}

	return nil
}
```

## Separating data from methods

If we dont use private fields, we end up exposing the fields to mutations. This is fine, granted that we go with immutable domain models. This domain models are mostlyread only fields with possible some getters for computed value. Those getters could call another function or service to decorate the model with more meaningful data. For example, a person object may have a getter age, that calls a function calculate age.
This allows calculate age to function to be used elsewhere, while avoiding imports when calling the getter.

Thedomain model should not have any methods that mutates itself. Instead, it should be passed to a service or function that computes the next state, akin to redux in React. The reason is simple - it is more natural to call updateDateOfBirth(user, dob) than user.withDateOfBirth/setDateOfBirth/updateDateOfBirth(dob). Generally, the setter methods assumes the user exists or have been loaded before the setter could be called. However, most of the time, it is more performant to just call update directly to the repository layer, after validating that the params are valid.

Another reason is that the behaviour of the mutation might vary based on business logic. Having several more specific methods is better than having conditionals. Setters for the purpose of setting individual fields are rarely useful (validation could have been replaced by construction of valueobject which is always valid).

https://wiki.c2.com/?FunctionsAndDataAreSeparate


## Read-only entity

The traditional OOP way is to create a class with private properties, and the properties can only be accessed through getters, and mutated through setters. However, that does not seem to work well for every scenarios, especially modern applications due to the following reasons:


- it is expensive to load entity. Say for example you want to update a user's email. It is expensive to query the user from the storage first, then calling the method `SetEmail` and then only update the user in the storage. Instead, it will be simpler to just issue `storage.UpdateUserEmail(email)`.
- state lives in database, not app. This is a follow up to the point above. Another scenario to consider is when performing bulk update. It is far easier to update the database columns then to follow the fetch/update/save lifecycle.


How does this affect the way we write code?

- entities are now dumb getters, so we don't need private properties.
- we can separate business logic for mutation/setters. Instead of `SetEmail` which checks if the user can set the email and the email is valid, we encapsulate the behaviour for validating email to another method (service layer). That way, the email can always be validated without first loading the user entity.
- code becomes less OOP, and more data-oriented. We just write methods to transform/process the data.
- models become behaviourless. The less behaviour we tie to our entity, the better. Note that there are certain scenarios where we still want to load the entity before checking if an action can be carried out. For example, if we do soft delete, then the user still exists in the db. So we have to check if the user is not soft deleted before updating the email.

```js
// Before.
class UserUseCase {
	async updateUserEmail(id, email) {
		// Load the entity.
		const user = await this.repo.find(id)

		// Set (and validate) the email.
		user.setEmail(email)

		// Save the entity.
		await this.repo.save(user)
	}
}

// After.
class UserUseCase {
	async updateUserEmail(id, email) {
		// Validate the email.
		assertIsEmail(email)

		// Issue command to db to update the user's email.
		await this.repo.updateUserEmail(id, email)
	}
}
```

What other implications that this have to other layers?
- since our domain models are purely getters, we can pass them outside the usecase layer without worrying that clients will call setters that may modify the state of the entity
- since behaviour is external, we can reuse it across other entities.


## References

1. [StackOverflow: Validation in Domain Model of Domain Service](https://stackoverflow.com/questions/35934713/validation-in-domain-model-of-domain-service)
2. [StackOverflow: Using Repository in Entity for validation before update](https://stackoverflow.com/questions/55549616/ddd-using-repository-in-entity-for-validation-before-update)
3. [Quasiclass](http://www.idinews.com/quasiClass.pdf)
