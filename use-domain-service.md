# Use Service

## Pure functions
Example of good domain service:

```go
// model.go

type User struct {
	bannedAt *time.Time
}

func (u *User) IsBanned() bool {
	return u.bannedAt != nil
}
```

Domain service calls the domain model methods, and returns specific errors for the invariant.
```go
var ErrUserBanned = errors.New("user: banned")

func UserIsActive(u *User) error {
	if u == nil {
		panic("user: not found")
	}

	if u.IsBanned() {
		return ErrUserBanned
	}

	return nil
}
```

## With dependencies

If an operations contains side-effects, such as external dependencies, IO, random numbers or date, it might be better to separate them from the domain model for testability:

```go
type Card struct {
	expiredAt time.Time
}

// Bad: hard to test expiry logic
func (c *Card) IsExpired() bool {
	return time.Now().After(c.expiredAt)
}
```

Better design:
```go
type CardValidator struct {
	now func() time.Time // Prefer function over interface for simple factories.
}

var ErrCardExpired = errors.New("card: expired")

func (v *CardValidator) CardMustNotBeExpired(c *Card) error {
	if v.now().After(c.expiredAt) {
		return ErrCardExpired
	}
	return nil
}
```
