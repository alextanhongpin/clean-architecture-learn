# Repository

In golang, there is no domain repository layer, which is purely the interface to the repository. Instead, the repository is just an implementation, usually the adapter that is passed in the usecase.

- dealing with transaction, unit of work
- dealing with events, outbox
- saga pattern
- row locking and custom locking
- access control
- orm vs raw sql
- mapping values to and from domain model and dto
- pagination, sort and filter
- error handling
- constraints and business logic in db

- storage, clients and pubsub


## Pagination


usually pagination is not really a domain concern. Most data related to query, if void of business logic are usually not part of domain, unless they have some complex conditional (e.g. dont show out of stock products for certain categories).
