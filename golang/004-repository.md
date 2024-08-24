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
