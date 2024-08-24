# Presentation

Presentation layer is probably the easiest to work with.

This layer is where the client (aka end users) connects to our API.

It can be 
- http (rest)
- grpc
- cli
- graphql
- etc

We can have multiple. In golang, it is common to have a `cmd` directory that contains scripts to be executed.

This layer has the following responsibility
- defines request/response types
- parse inputs, such as command line args, query string, forms or json body
- sanitize inputs
- validates inputs aka validation errors
- ~converts inputs into DTO for usecase layer~ this is done at app layer, the presentation layer defines the interface
- handles error
- converts domain errors to presentation specific error (e.g. exit code or json or html errors)
- performs authorization, which is not business logic
- ~transforms usecase output to presentation specific output~, this is done at app layer, the presentation layer defines the expected types 
- middlewares for presentation layer specific logic (e.g. passing user id context, gzip compression)

This layer does not contain any business logic

- presentation | parse | sanitize | validate | transform | usecase | render

Although most people are familiar with controllers, the more appropriate naming for golang is handlers

we can use the following folder to represent one or more presentation layer

```
rest/
- api/v1/users.go
- api/v1/v1.go
- api/health.go
- handlers/user.go
- middleware/auth.go
graphql/
- resolvers/
- queries/
- mutations/
grpc/
cmd/send_user_email.go
```

```go
type UserHandler struct {
  uc userUsecase
}

func (h *UserHandler) Get/List/Create/Update/Delete/Search
```

We defined our handlers separately from routes.

We also define a generic transformer that maps the input/output to the usecase layer.

```go
GetUserEncoder(ctx, GetUserRequest) (T, error)
GetUserDecoder(ctx, T) (GetUserResponse, error)
```

Or we should probably put it at usecase layer, wrapping the original usecase.

All these mappers should be handled at the application constructor.

Since each layer is not supposed to know each other (ala no imports), the mapping should only be done in one layer, which is the application layer (not to be confused with application service and controller layer). 
By having all the mappings done in one layer, you avoided props drilling.


```go
func toUsecaseDto(api.Request) usecase.Dto {}
func toRepoParams(usecase.Dto) db.Params {}
func userResult(u *domain.User, err error) (*api.User, error) {}
```
