# Structure

How do we structure our application to use clean architecture? It really depends on the language used, because each has their own definition of modules (or packages/dependencies etc). Here's how for `golang`.

Four main folders:
```
your-domain
|- http # Primary adapter (aka inputs) lies here. 
|  |- api
|  |  |- root/
|  |  |- v1/
|  |  |- v2/
|  |  `- api.go # Maps domain entity to external in/out DTOs
|  |- error # Maps domain errors to client errors, e.g, http errors
|  |- middleware
|  |- security # contains authorization logic (authorization header extraction, jwt signing and validating, NOTE: jwt implementation is not part of infra)
|  `- session # contains propagation context, but only within ui layer.
|
|- graph/ # You can also have multiple primary adapters (websocket, graphql, grpc, rpc, cli)
|- cli/
|
|- infra # Secondary adapter (aka outputs) lies here.
|  |- postgres # Postgres port (port is the interface, adapter is the implementation)
|  |  |- repository/ # Repository implementation may lie here, or in the application service.
|  |  |- migration/ # Adapters
|  |  |- seed/ 
|  |  `- postgres.go
|  |- oauth/
|  |- mailer/
|  |- payment/ # The domain services should depend on the interface.
|  |- validator.go
|  |- logger.go
|  |- redis.go # Redis port/adapter
|  `- kafka.go
|
|- domain # This layer does not depend on any external layer.
|  |- entity # Entities fields are primitives or value object, no sql.NullString etc.
|  |  |- user.go
|  |  `- book.go
|  |- service # Domain services lies here.
|  |  `- user.go 
|  |- repository # Only interface, not implementation.
|  |  |- user.go
|  |  `- book.go
|  |- value # Value objects
|  `- event
|
`- app # Application services lies here. Group by subdomain. They can depend on infra + domain layer.
   |- authentication
   |  |- register.go
   |  |- reset_password.go
   |  |- ...
   |  `- login.go
   |- authentication.go # Facade to individual authentication usecases. This is meant as a "stable" layer, a contract that doesn't change. The implementation may change.
   `- another usercase/
```
