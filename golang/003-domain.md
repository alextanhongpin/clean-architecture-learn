# Domain

Domain is most important layer in clean architecture 


The following is the responsibility of domain layer

- define the entities
- performs global business logic, aka business logic that can be used at usecase layer (presentation layer should not have access to domain, since it may mutate domain or have side-effects)
- defines the global errors that is tied to the business logic, note that there are also errors specific to repo/usecase or maybe aggregate that might not be here
- can have aggregate and aggregate business logic
- does domain logic validation
- can have value objects, see more on `types` proposal
- can have constructors for default valuea
- holds constants, enums
- when an entity has reference to another entity, it is called aggregate
- aggregate implies lifefcycle, if the root is deleted, all the children is deleted
- service performs operations on aggregate's entities
- can have collections of entities
  
What it isn't 
- a representation of db table
- anaemic without business logic
- nested entities, we just have reference to the id, not nested
- should not contain any external dependencies or config

