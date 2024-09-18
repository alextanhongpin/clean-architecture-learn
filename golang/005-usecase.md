# Usecase

- orchestrates business logic (calls domain business logic, not inline business logic)
- handles transaction through unit of work


## Error Handling


How does usecase handles error from another layer?

Take for example, a usecase that calls a repository to fetch a user. If the repository returns an error, how should the usecase handle it?

```go
// user_usecase.go
type UserRepository struct {
    FindUser(ctx context.Context, id string) (*domain.User, error)
}

type UserUsecase struct {
    userRepo UserRepository
}

func(uc *UserUsecase) FindUser(ctx context.Context, id string) (*domain.User, error) {
    user, err := uc.userRepo.FindUser(ctx, id)
    if err != nil {
        return nil, err
    }

    return user, nil
}
```

Above, the usecase is simply returning the error from the repository. This is a common pattern in golang, but it is not always the best way to handle errors.

We simply do not know what error to expect, and we are not handling the error in any way. This is not a good practice.

Instead, we should wrap the error with a domain error, and return it to the presentation layer.

```go
func(uc *UserUsecase) FindUser(ctx context.Context, id string) (*domain.User, error) {
    user, err := uc.userRepo.FindUser(ctx, id)
    if errors.Is(err, sql.ErrNoRow) {
        return nil, domain.ErrUserNotFound
    }
    if err != nil {
        return nil, err
    }

    return user, nil
}
```

Above, we are wrapping the error with a domain error, `domain.ErrUserNotFound`. This is a domain error that is specific to the usecase, and the presentation layer can handle it accordingly.

Unfortunately, this is not ideal. The usecase layer should not be aware of the repository error. Instead, the repository should return a domain error.

```go
// user_repository.go
type UserRepository struct {
}

func (r *UserRepository) FindUser(ctx context.Context, id string) (*domain.User, error) {
    user, err := r.db.FindUser(ctx, id)
    if errors.Is(err, sql.ErrNoRow) {
        return nil, domain.ErrUserNotFound
    }
    if err != nil {
        return nil, domain.ErrInternal
    }

    return user, nil
}
```

And our usecase is back to the original implementation:
```go
// user_usecase.go
func(uc *UserUsecase) FindUser(ctx context.Context, id string) (*domain.User, error) {
    user, err := uc.userRepo.FindUser(ctx, id)
    if err != nil {
        return nil, err
    }

    return user, nil
}
```

Which leaves us back to the original question, how should the usecase handle the error from another layer?
The usecase should know which errors to expect, and the repository should implement it. We choose to make the errors explicit in the usecase layer, and wrap the error as unknown error if it is not expected.


```go
func(uc *UserUsecase) FindUser(ctx context.Context, id string) (*domain.User, error) {
    user, err := uc.userRepo.FindUser(ctx, id)
    if err != nil {
        if !errors.Is(err, domain.ErrUserNotFound)  {
            return nil, UnknownError(err)
        }

        return nil, err
    }

    return user, nil
}
```

With this, we now know what errors to expect in the usecase layer, and we can handle it accordingly. The presentation layer can then handle the error as needed.

This makes it easier for the repository implementer to know what errors to return, and the usecase implementer to know what errors to expect.


## Conclusion ...

Stick with the original solution...
```go
// user_usecase.go
func(uc *UserUsecase) FindUser(ctx context.Context, id string) (*domain.User, error) {
    user, err := uc.userRepo.FindUser(ctx, id)
    if err != nil {
        return nil, err
    }

    return user, nil
}
```

And define the list of errors in the domain layer instead ...

```go
// domain/user.go
var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user exists")
```

## Patterns

- break usecases: when the usecases gets too large, there may be sign it could be composed of unrelated usecases. We can apply the break pattern and separate unrelated usecases
- group usecase: the opposite of break, we group related usecases, e.g. login and register under AuthenticationUsecase. CRUD under manage entity.
- 
