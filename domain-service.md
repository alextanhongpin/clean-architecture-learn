# Domain service

- See [2], this provides a lot of info about domain service.
- part of domain layer
- operates on aggregates
- contains business logic such as validation
- should have as few services as possible to avoid anaemic entity model
- often misunderstood with _application service_ layer
- domain services should just be pure function, with the intent in the function name, and without the `Service` suffix
- if the domain service requires dependencies, do not pass them as arguments of the function, as it becomes coupled to the method, unless the dependency is declared as an interface
- domain model should not be able to call domain service, however domain service can call domain model
- it is preferable to return booleans in domain model and domain model methods should only have one responsiblity. Domain service however can chain multiple domain model methods and return specific errors
- unaware of infrastructure or overall application flow, they exclusively encapsulates business logic rules
- a service that expresses a business logic that is not part of any aggregate root
    - e.g. you have two aggregate, Product which contains name and price. Purchase which contains purchase date, list of products ordered with quantity and product price at that time, and the payment method.
    - Checkout is not part of either of these models and is a concept in your business
    - Checkout can be created as  a domain service which fetches all product and compute the total price, pay the total by calling another domain service PaymentService with an implementation part of infrastructure and convert it into purchase 

Can domain services call repository?
- yes
- may call repository, but should be avoided. Instead, query the entity though _repository_ at the _application service_ layer, and pass the entity to the `domain service`. See [1]
- if the service does call repository, it should be for read-only purposes. Leave the persistence of an entity to _repository_.
- since the _application service_ shouldn't house business logic, if there is a logic that requires access to repository, then the best place is to put it in _domain service_. Ideally _domain service_ should not call repository. If possible, fetch all related entities (or aggregates) in the _application service_, then pass them to the _domain service_ to execute the business logic.
- repository can be injected into domain services in rare scenario, it is the application layer that does it most of the time. The domain service should not perform any mutation to the repository (like saving etc), only read
- does not persist entity, that is the role of Repository


When to use domain services?
- encapsulates business logic that doesn’t naturally fit within a domain repository, and are NOT typical CRUD operations
- concepts from the domain that don’t seem to fit as model object end up forming one or more domain services
- encapsulate such behaviour that do not fit a single domain object
- does something which makes sense only when being done with other collaborators (domain objects or other services)
- deals with everything related to domain objects, but go beyond the scope of a single entity with focus on business rules.
- stateless, unlike entities
- can be modelled as pure functions, but typically as a class just for name-spacing purposes
- methods that don’t fit on a single entity or require access to the repository are contained within domain services

What is the input/output:
- the domain services are expressed in terms of ubiquitous language and the domain type, the method arguments and the return values are proper domain classes
- domain services accepts domain entities or value objects, carry out conditional operations on those primitives or objects, or performs business rule calculations, and then return primitives or domain entities or value objects

## Service

We know that domain service layer
- is part of the domain layer
- they should be stateless 
- encapsulate business logic that doesn't naturally belongs to an entity
- operates on a collection of entity

This is how we usually do it
1. Namespace the service using a struct
```
type UserKtpService struct {}
```
2. Have a validate<entity/field> method
```
func (u *UserKtpService) ValidateKTP(nik int64) error {}
func (u *UserKtpService) validateAge(dob time.Time) error {}
func (u *UserKtpService) validateAreaCode(dob time.Time) error {}
```

Disadvantages of the approach above:
- methods can be hard to test, if they are private
- need to instantiate the whole struct before calling the methods
- too procedural and lacks facts, e.g. age < KtpMinAge is not as clear as UserKtpMustHaveLegalAge(age)

Suggestions:
- service methods should read like a business rules, expressing facts
- pattern to use: <entity> must equal <fact> instead of validate<entity>, which is more readable
- the methods should not sound ambiguous (e.g. isPaymentMethodValid vs paymentMethodMustBeValid)
- methods that expresses action is still okay, e.g. calculateTotalCost

So instead of:
```
validate user ktp
- validate age
- validate area code
- validate date of birth
- validate gender
```
We do this:
```
user ktp must be valid
- user ktp age must be legal age
- user ktp area code is valid area
- user ktp date of birth is valid
- user ktp gender is male or female
```
In code, it will be individual functions:
UserKtpMustBeValid(nik)
UserKtpAgeMustBeLegalAge(dob)
UserKtpAreaCodeIsValidArea(areaCode)
UserKtpGenderIsMaleOrFemale(gender)

What do you all think? Is this a better way? :thonking:

# Reference

1. [StackOverflow: Can Domain Services access Repositories?](https://stackoverflow.com/questions/26930131/can-domain-services-access-repositories)
2. [Domain Driven Design- Tactical Patterns](http://domaindrivendesigns.blogspot.com/2018/11/domain-services-domain-model-domain.html?m=0)
3. [GitHub: DDD_Samples](https://github.com/VaughnVernon/IDDD_Samples)
