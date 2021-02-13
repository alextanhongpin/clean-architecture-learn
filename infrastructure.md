# Infrastructure Layer

- a.k.a secondary adapter from Hexagonal architecture
- used to abstract technical concerns
- contains implementation of the interfaces from the domain layer, such as repository implementation (should this be in application service?)
- any domain concepts (that includes domain services and repositories) that depends on external resources are defined by interfaces, and are implemented here
