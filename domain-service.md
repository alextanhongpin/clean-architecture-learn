# Domain service

- See [2], this provides a lot of info about domain service.
- part of domain layer
- operates on aggregates
- contains business logic such as validation
- should have as few services as possible to avoid anaemic entity model
- may call repository, but should be avoided. Instead, query the entity though _repository_ at the _application service_ layer, and pass the entity to the `domain service`. See [1]
- does not deal with persistence
- often misunderstood with _application service_ layer
- since the _application service_ shouldn't house business logic, if there is a logic that requires access to repository, then the best place is to put it in _domain service_. Ideally _domain service_ should not call repository. If possible, fetch all related entities (or aggregates) in the _application service_, then pass them to the _domain service_ to execute the business logic.


# Reference

1. [StackOverflow: Can Domain Services access Repositories?](https://stackoverflow.com/questions/26930131/can-domain-services-access-repositories)
2. [Domain Driven Design- Tactical Patterns](http://domaindrivendesigns.blogspot.com/2018/11/domain-services-domain-model-domain.html?m=0)
