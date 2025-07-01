# Domain Layer 

Domain layer is where the business logic lies.

Since almost every layer has access to dimain, this makes it the best place to put stateless business logic.

The domain layer is also the single source of truth for the types.

This is because the application can retrieve data from different sources, e.g. database, third party api calls. However, they will all be translated to domain layer.


## Domainless layer

