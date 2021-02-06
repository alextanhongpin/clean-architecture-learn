# Domain service

- part of domain layer
- operates on aggregates
- contains business logic such as validation
- should have as few services as possible to avoid anaemic entity model
- may call repository, but should be avoided. Instead, query the entity though _repository_ at the _application service_ layer, and pass the entity to the `domain service`.
- does not deal with persistence
- often misunderstood with _application service_ layer
