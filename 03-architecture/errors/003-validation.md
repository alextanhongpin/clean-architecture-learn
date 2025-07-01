# Validation in DDD

Where does validation happens? First of all, we need to discern between the types of errors we have. Field-level errors (or input validation) are usually not a concern of a domain (e.g. empty value, space untrimmed, check length, correct input format, correct data type), see [2]. Domain errors are rules for the domain, so `name can only consist of alphabets` and `name cannot be longer than 55 characters` should have their respective errors. A __value object__ can also be created to represent the `name` value. As such, we can delegate input validation to the application service, and handle domain errors in domain itself. This is similar to the __two-step validation__ is proposed by Microsoft see [3]. 

Validation patterns:
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
2. [StackOverflow: Should value objects contains technical validation for input parameters](https://stackoverflow.com/questions/39224430/should-value-objects-contain-technical-validation-for-input-parameters/39234321)
3. [Design validations in the domain model layer](https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/domain-model-layer-validations#two-step-validation)
