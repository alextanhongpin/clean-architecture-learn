# Domain service

- See [2], this provides a lot of info about domain service.
- part of domain layer
- operates on aggregates
- contains business logic such as validation
- should have as few services as possible to avoid anaemic entity model
- may call repository, but should be avoided. Instead, query the entity though _repository_ at the _application service_ layer, and pass the entity to the `domain service`. See [1]
- does not deal with persistence
- often misunderstood with _application service_ layer


# Reference

1. [StackOverflow: Can Domain Services access Repositories?](https://stackoverflow.com/questions/26930131/can-domain-services-access-repositories)
2. [Domain Driven Design- Tactical Patterns](http://domaindrivendesigns.blogspot.com/2018/11/domain-services-domain-model-domain.html?m=0)
