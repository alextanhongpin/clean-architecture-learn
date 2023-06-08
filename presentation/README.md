# Presentation Layer

- The presentation layer is where the end users interacts with the application service
- also commonly known as the controller layer
- Presentation layer can be Graphql, REST API, CLI, GRPC or any other transport mechanism
- Presentation layer should be void of domain business logic, even basic if/else
- Versioning is usually handled in presentation layer
- Top-level authorization should be handled here (e.g. JWT Token verification). Domain layer _may_ contain domain-specific authorization.
- the presentation layer is also responsible for transforming the inputs (query string, params, headers) into usecase inputs
- will usually be documented using tools like Swagger, because it will be the main interface for clients to communicate to the backend

## Roles and Responsibilities 

- request parsing: parse raw data into DTOs, or data transfer object
- response parsing: map response from application service layer into 
 presentation layer type. For JSON API, it can be serialising the output object/class into JSON payload

