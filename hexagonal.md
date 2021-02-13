# Hexagonal Architecture

Hexagonal architecture share some similarity with Clean Architecture (so does Onion architecture). While they are not entirely the same, some patterns such as port/adapters can be applied successfully.

There are primary and secondary ports/adapters, which is essentially just another name to describe the input and output of the system.

Ports are just `interfaces`, while `adapters` are the implementations. Following _dependency injection_ principles, we only depend on the interfaces, and not the implementation. Hence `adapters` are rarely called in any parts of the application, except its constructor.

We will typically separate the primary ports and secondary ports folders as such:
- primary ports can be `http/`, `graph/` (for graphql), `grpc/`, `cli/` etc
- secondary ports are named `infra/` and contains the interface as well as implementation for database, repository, cache, logging, message bus etc.

# References

1. [Medium: Ports & Adapters architecture](https://wkrzywiec.medium.com/ports-adapters-architecture-on-example-19cab9e93be7)
2. [DevTo: Hexagonal architecture: Ports and Adapters](https://dev.to/jofisaes/hexagonal-architecture-ports-and-adapters-1h4m)
3. [Layers, ports & adapters - Part 3, Ports & Adapters](https://matthiasnoback.nl/2017/08/layers-ports-and-adapters-part-3-ports-and-adapters/)
