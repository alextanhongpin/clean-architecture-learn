# Repository


Aka persistence layer.
Use repository to
- read/write data to any data source, not limited to database, can be API calls
- returns domain models from those data sources
- calls external API and transforms the response to the domain models
- abstracts cache/message queue implementation
- can return more than one domain models type
- does not imply one repo = one database table, this is wrong
- can have business logic like unique constraints in db, converts those to domain errors

Acts as a facade to
- storage db
- API
- cache (can be implemented as wrapper)
- message queue

While each of them can be an interface, it is not strictly required since they are all side-effects, aka mocking them is useless.

On the other hand, we want to mock the entire repository, so placing all the side effects reduces a lot of code in the usecase layer.

```go
func GetUser() {
}
```


## Defining repository 

In golang, we only need to define a **single repository** in `repository/repository.go`>

So if you have multiple tables, you don't create a `UserRepository` or `ProductRepository` or so on.

```go
type Repository struct {
}
func (r *Repository) CreateUser()
func (r *Repository) CreateProduct()
```

Instead, we have a single repository, with the methods containing the entity suffix.

The advantages are
- simple construction
- the caller may need to access multiple resources, and using a single repository allows that
- we can define interface on what method is allowed

## Business Logic

The repository handles the following

- setting up transaction
- handling transaction context (creating a method that supports transaction)
- convert sql errors to domain errors (e.g. not found, unique constraints)
- convert postgres types to domain entities
- mapping params to postgres params

We will cover all of them shortly.

## Repository Params

Where does repository params belong? The answer is just in the same directory.

The repository should not depend on external types.

Take for example `CreateUser`. If we have a multiple usecases that calls the create user, they both depends on the type, rather than two separate usecase types.

Should we define the repository types in entity folder? While logical, it does not belong there.


Conclusion, define repository params in repository directory.

## External calls

Aside from database, we should also place remote api calls or message queue or redis.

For redis, we do embedding.

For the rest, we use special prefix like `RemoteCreateUser` to indicate it is an external API call or `Enqueue` to indicate it is a message queue.

