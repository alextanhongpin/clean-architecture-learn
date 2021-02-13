# Structure

How do we structure our application to use clean architecture? It really depends on the language used, because each has their own definition of modules (or packages/dependencies etc). Here's how for `golang`.

Four main folders:
```
your-domain
|- ui # Primary adapter lies here. 
|  |- api
|  |  |- root/
|  |  |- v1/
|  |  |- v2/
|  |  `- api.go
|  |- middleware
|  |- security # contains authorization logic (authorization header extraction, jwt signing and validating)
|  `- session # contains propagation context, but only within ui layer.
|
|- infra # Secondary adapter lies here.
|  |- postgres
|  |  |- migration/
|  |  |- seed/
|  |  `- repository/ # Repository implementation may lie here, or in the application service.
|  |- redis.go
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
   |- adapter/
   |  |- repository/ # Another option to place repository implementation.
   |  |- error/ # Maps domain errors to client errors, e.g, http errors
   |  |- dto/ # Maps domain entity to external in/out DTOs
   |  `- translations/ 
   |- usecase/
   `- another usercase
```
