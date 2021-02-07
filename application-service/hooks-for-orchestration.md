# Using hooks for orchestration

```go
package main

import (
	"context"
	"errors"
	"log"
)

var (
	ErrUserNotFound    = errors.New("not found: user")
	ErrUniqueViolation = errors.New("sql: unique violation")
	ErrEmailExists     = errors.New("email exists")
)

type User struct {
	unconfirmedEmail  string
	email             string
	encryptedPassword string
	plaintextPassword string
}

type Hook interface {
	Before() error
	Exec() error
	After() error
}

type Pipeline struct {
	steps []func() error
}

func (p *Pipeline) Add(step func() error) {
	if step == nil {
		return
	}
	p.steps = append(p.steps, step)
}

func (p *Pipeline) Exec() error {
	for _, step := range p.steps {
		if err := step(); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) Create(hook interface{}) error {
	if hook == nil {
		return nil
	}
	p := new(Pipeline)
	if h, ok := hook.(interface{ Before() error }); ok {
		p.Add(h.Before)
	}
	if h, ok := hook.(interface{ Exec() error }); ok {
		p.Add(h.Exec)
	}
	if h, ok := hook.(interface{ After() error }); ok {
		p.Add(h.After)
	}
	return p.Exec()
}

func main() {
	u := new(User)
	u.unconfirmedEmail = "john.doe@mail.com"
	hook := &CreateUserHook{
		user: u,
		repo: new(UserRepository),
		ctx:  context.Background(),
	}
	if err := u.Create(hook); err != nil {
		log.Println(err)
	}
}

type userRepository interface {
	Create(ctx context.Context, u *User) error
}
type CreateUserHook struct {
	user *User
	repo userRepository
	ctx  context.Context
}

func (h *CreateUserHook) Before() error {
	h.user.encryptedPassword = "password"
	return nil
}

func (h *CreateUserHook) Exec() error {
	if err := h.repo.Create(h.ctx, h.user); err != nil {
		// This suffer the same problem, the hook knows about the database errors.
		if errors.Is(ErrUniqueViolation, err) {
			return ErrEmailExists
		}
	}
	return nil
}

type UserRepository struct{}

func (r *UserRepository) Create(ctx context.Context, u *User) error {
	if u.unconfirmedEmail == "john.doe@mail.com" {
		return ErrUniqueViolation
	}
	return nil
}
```
