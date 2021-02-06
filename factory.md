# Factory

- what is a factory? A Factory is an object that has the single responsibility of creating other objects.
- _factory_ marks the __beginning__ of an entity, _repository_ the __middle__ and __end__
- use _factory_ to build the entity, and _repository_ to persist
- one deals with creation, another deals with persistence.
- domain factory can also be user the other way around, for reconstitution, meaning to build an entity back from external sources, whether it’s from JSON (deserialising) or from another domain service or even repository
- can factory access repository (? yes, but avoid that if possible. Are there any good usecase to demonstrate it?)
- factory is responsible for constructing complex aggregates in their valid states. So a constructor can be a factory too
- factory can be plain functions, for complex factory that requires other 

# References

1. [Domain Driven Design: Services and Factories](https://archiv.pehapkari.cz/blog/2018/03/28/domain-driven-design-services-factories/)
