# Presentation Layer


- The presentation layer is where the end users interacts with the application service
- Presentation layer can be Graphql, REST API, CLI, GRPC or any other transport mechanism
- Presentation layer should be void of domain business logic, even basic if/else
- Versioning is usually handled in presentation layer
- Top-level authorization should be handled here (e.g. JWT Token verification). Domain layer _may_ contain domain-specific authorization.
