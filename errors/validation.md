# Validation in DDD

Where does validation happens? See [1].

1. `isValid` method. This allows model to be in invalid state, which is a huge drawback. 
2. validating in application services layer, a.k.a. two-step check. Ensuring the dtos are valid means the state will always be valid (unless there are multiple steps in between that creates/modifies parameters). However, the business logic leaks to application service.
3. TryExecute: validate when calling domain method. The entity holds the domain knowledge and can ensure consistency.
4. Execute/CanExecute: allows the reuse of the validation.

- does entity throws error? yes, the place to maintain the errors should be close to the entity as possible. Outer layers can wrap the errors for more details. Ideally use the `Execute/CanExecute` method.
- does repository return domain errors? yes, there are some errors that can only be returned by the database, such as unique constraint violiation on unique email during user registration, check constraints in database etc. We do not want to cast as well as leak database error handling to the application service, so the best way is to let the repository return the errors. `EntityNotFound` errors for specific entity should be returned here too.
- does service return domain errors? absolutely, but prioritize putting it in domain entity if possible.
- does application service return domain errors? Nope. Ideally all domain errors should be translated to errors that will be used by the client of application service, e.g. rest adapters that maps domain errors to http errors.

# References
1. [Validation and DDD](https://enterprisecraftsmanship.com/posts/validation-and-ddd/)
