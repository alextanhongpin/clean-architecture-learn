# Repository business rule

- repository can contain business rule
- filtering is part of business rule too (e.g. find only blue balls, find blue balls count)

Bad: 
- database-specific errors leaking to the application service
- swapping database is not possible
- entity cannot call repository, but the opposite holds true, use entity in repository to return domain errors. Better, provide an interface for the entity method to be call upon succes/failure

```go
package main

import (
	"context"
	"database/sql"
	"errors"
)

var ErrUserNotFound = errors.New("not found: user")

type User struct {
	ID string
	Name string
}

type userRepository interface {
	Find(ctx context.Context, id string) (User, error)
}

type ApplicationService struct {
	repo userRepository
}

func (s *ApplicationService) FindUser(ctx context.Context, id string) (User, error) {
	u, err := s.repo.Find(ctx, id)
	if errors.Is(sql.ErrNoRows, err) {
		return u, ErrUserNotFound
	}
	return u, nil
}
```

Better:
- let the repository return the domain error
```go
type UserRepository struct {}

func (r *UserRepository)FindUser() (User, error) {
	// REDACTED
	return User{}, ErrUserNotFound
}
```
