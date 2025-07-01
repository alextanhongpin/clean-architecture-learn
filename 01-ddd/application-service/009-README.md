# Application Service Implementation

A list of problems faced when implementing DDD, and some possible solutions.



## Pre/post-hook

The examples shown in [persistence validation](https://github.com/alextanhongpin/clean-architecture-learn/blob/main/application-service/persistence-validation.md) and [enforcing unique constraints](https://github.com/alextanhongpin/clean-architecture-learn/blob/main/application-service/enforcing-unique-constraints.md) deals with entity validation that is tied to the persistence layer. 

For most cases, we would like to perform validation before or after the repository method, but having the errors within the domain layer. This can be solved by using the `visitor` as well as `decorator` pattern.


## Responsibility

- pre validation
- orchestrate business logic
- initialize unit of work
- maps primitives and dto to domain types
- handles authorization and permission
- creates and executes donain services
- calls repository, message bus
- returns read only domain types in most cases, otherwise an immutable thpe
- sandwich of imperative shell

## what does this layer does 

control flow of the application 
basically a series of steps to execute 
- each step is usually getting input data, and transforming to output data
- manages transaction in an application, through unit of work

## what this layer shouldnt do

inline business logic

