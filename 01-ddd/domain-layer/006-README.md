# Domain Layer

This directory contains ADR (Architecture Decision Records) on writing better domain layer in Golang. The implementation of domain layer in other language may differ due to language features, and hence the examples here may not be idiomatic in other languages.


Creational
- factory, builder and constructor
- avoiding mega constructor

Behavioral
- getter, setter and wither
- dealing with enums, value object
- enums as data from database
- constraints in database vs in memory
- separate types (date, set, stringcase, range)
- dealing with null types vs pointers
- error handling
- polymorphism
- computed value
- statuses and state machine
- concurrency

Tactical
- usecase and user story
- entity vs model representation

## What is it

Domain layer are where entities lies.

It is called domain, because it stored domain-specific business logic, e,g ecommerce, payment etc.

Business logic is just logic required by business to function. Without the business component, it is just logic.

Some example includes 
- the user can have max 1 membership discount per year
- or as simple as only alphanumric name can be used to register an account

Business logic should be separated from access control logic (only user that own the x can delete x), because it is kinda expected behaviour.

The same goes to most validation logic, e.g. optional vs required. Occassionally you will find some overlap, e.g price cannot be negative (validation rule) but also cannot exceed certain amount (business rule). Of course, it is not wrong to deduplicate the logic to domain layer (innermost layer), since the validation is guaranteed, but usually validation is done at the top layer at the ease of short-circuiting the operation (fail fast) and returning the reason for failure on the client side.

With that said, business logic appears in most layers, e.g repository has some postgres contraints on unique column or not null etc. However, if the logic has no external dependencies, it should belong in the domain layer, where the entire apo can call it, since most layer has access to domain types (depending on your usecase, you may want to convert those domain types to immutable read only values when returning them to presentation layer).

If the business logic is placed outside of domain layer, ensure that the dependency rule is respected, only the top layer can call the bottom layer and not the other way around.


## Functional vs OOP

That depends on you, but if the domain type has all the necessary info, then just using methods would be better. The exceptionns are
- the data depends on external source
- the logic does not require instantiation of a domain type, or it is expensive to do so
- normally only a property of the domain needs to have the logic executed, e,g password hashing, checking email validation match


