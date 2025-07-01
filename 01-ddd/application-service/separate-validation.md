# Separate requeset validation


Say for example, you have a usecase that handles user registration. However, notice that we first need to validate the request.


```go
package usecase

import "context"

type AuthUsecase struct {
	repo authRepository
}

type RegisterDto struct {
	Email    string
	Password string
}

func (uc *AuthUsecase) Register(ctx context.Context, dto RegisterDto) error {
	if err := Email(dto.Email).Validate(); err != nil {
		return err
	}

	if err := PlaintextPassword(dto.Password).Validate(); err != nil {
		return err
	}

	// Do login ...
	return nil
}
```

What this means in the integration tests is that we need to provide different input to test the validation logic, which will only add more _branching_.


A better approach is to test the validation separately.
We can shift the validation logic to the `dto` (Data Transfer Object).


```go
package usecase

import "context"

type AuthUsecase struct {
	repo authRepository
}

type RegisterDto struct {
	Email    string
	Password string
}

func (dto RegisterDto) Validate() error {
	return Validate(
		Email(dto.Email),
		PlaintextPassword(dto.Password),
	)
}

func (uc *AuthUsecase) Register(ctx context.Context, dto RegisterDto) error {
	if err := Validate(dto); err != nil {
		return err
	}

	// Do login ...
	return nil
}

type validatable interface {
	Validate() error
}

func Validate(vs ...validatable) error {
	for _, v := range vs {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}
```

For testing now, we can assume that the input for the `Register` method is always valid. We can also test the `RegisterDto` separately.
