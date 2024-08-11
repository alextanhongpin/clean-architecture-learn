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



## Repository is not the database layer

A lot of repositories are implemented as a layer that the application use to interact with the database. However, that is a wrong usage of repository.

Database sources (or any other data source) still needs to be mapped to the domain layer. The mapping happens in the repository layer, not in the database layer.

You can define a package named `postgres` that does mapping to the tables.

The repository layer should then call this `postgres` layer to return the types.

## Common questions

1. where to place the repository? in a separate package?
   - place it together with the usecase. the usecase will just define the repository interface
   - the `postgres` (or `storage` layer) can be defined in a separate package instead
2. what does the repository return
   - the domain types. When calling the `postgres` layer, it maps the postgres layer to the domain types
3. why do we use separate layer `postgres` to map to the db?
   - this layer is usually library dependent, such as using sqlc compiler or orm or other code generation. This is also database dependent. we do not want to mix external dependencies with our domain, but want to map them in the repository layer.
4. is the repository supposed to return only one entity type?
   - wrong, it should return an aggregate, or related entities. it can return just one, but not limited to it
5. what other responsibility does repository do?
   - external api calls, caching decision, enqueue to message queue is done at the repository layer
   - user did not need to know about caching decision in usecase kayer. Also, this simplifies tesing usecase because there are less dependencies to mock.
   - repository generally hold all side effects operations that are usually useless to test if mocked
6. how are errors handled here
   - repository also transform errors from postgres layer to domain errors. it converts sql violation such as unique constraints or no rows
