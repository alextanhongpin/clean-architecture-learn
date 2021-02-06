# Application Service Implementation

A list of problems faced when implementing DDD, and some possible solutions.



## Pre/post-hook

The examples shown in [persistence validation](https://github.com/alextanhongpin/clean-architecture-learn/blob/main/application-service/persistence-validation.md) and [enforcing unique constraints](https://github.com/alextanhongpin/clean-architecture-learn/blob/main/application-service/enforcing-unique-constraints.md) deals with entity validation that is tied to the persistence layer. 

For most cases, we would like to perform validation before or after the repository method, but having the errors within the domain layer. This can be solved by using the `visitor` as well as `decorator` pattern.
