# Adapter

- what are adapters?
- which layer do they belong to?
- what components are classified as adapters?
- consist of primary and secondary adapters
- external adapter to application (controller, api layer, middleware)
- internal adapter called by application (Postgres, message bus)

## How do they differentiate themselves from ports?

Ports is the interface, while adapter is the implementation. So our port can be for example `Storage`, while the adapter can be `InMemoryStorage` or `PostgresStorage`.

Not every interface should be suffixed with -er, e.g. reader, writer. It is preferable to prefix the implementation with the actual name, e.g. interface name `Cache` will have implementation `RedisCache` instead of `CacheImpl` or `cache`.


## Where do ports and adapters belong?

In golang, do not do Java-style interface. The ports (interface) should belong close to where client is. In short, there should not be a package to declare all the interfaces for the adapter, rather, we should only have the `adapters` package that contains all the ports/implementations.


## What is the input/output of an adapter?


The goal of the adapter is to really provide an interface for the external integration to play well with the domain layer.

They can for example, accept their own types declared in the package, domain types such as value objects and domain model as well as primitives.

However, what is important is that the domain layer should not depend on the adapter types.

## Where to place third party implementation

If you are implementing the third party REST call, then separate the client from the adapter.

The adapter layer should be responsible for converting the domain types to/from the adapter, with the condition that the domain layer should never depend on the types from the adapter.

## Layer

What does an adapter layer does? 

It only does conversion from one type to another, nothing else.

It is easy to validate the conversion happens. The target type should have all the fields filled.

No layers should be accessing each other.
