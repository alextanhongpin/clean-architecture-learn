# Domain Model


- Depending on your understanding, there may be no difference between domain model and domain entity 
- However, for certain application, it is better to split the domain model and entity
- entity are data that mimics the database table, while domain may be similar to or represent a subset of the entity
- in some language or when using, the domain entity may contain database specific mappings (e.g. golang's gorm or bun), however domain model should not have those details
- domain entity may also be composed purely of domain primitives, but when mapping to domain model, they should be mapped to value objects when possible
- the opposite is true, domain entity should not contain domain models or value object
- domain methods should not contain accepts external input - the computed values should be based solely on itself. Delegate such operations to domain services instead.
- side-effects (date time, random number) should be modelled externally for testability

## Domain model should not have methods depending on external input

```go
type Coupon struct {
	code      string
	expiredAt time.Time
}

func NewCoupon(code string, expiredAt time.Time) *Coupon {
	return &Coupon{
		code:      code,
		expiredAt: expiredAt,
	}
}

func (c Coupon) IsExpired() bool {
	return time.Now().After(c.expiredAt)
}
```

The above implementation, while it's valid, it violates the principle that _side-effects should be modelled externally for testability_.

However, refactoring it to this doesn't make it better to, since now the _method depends on external input_:
```go
func (c Coupon) IsExpired(at time.Time) bool {
	return at.After(c.expiredAt)
}
```



The alternative it separate the logic into a dedicated domain service:

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	coupon := NewCoupon("VALID", time.Now().Add(-1*time.Second))
	validator := NewCouponValidator()
	fmt.Println(validator.IsExpired(coupon))
}

type Coupon struct {
	code      string
	expiredAt time.Time
}

func NewCoupon(code string, expiredAt time.Time) *Coupon {
	return &Coupon{
		code:      code,
		expiredAt: expiredAt,
	}
}

type CouponValidator struct {
	now func() time.Time
}

func NewCouponValidator() *CouponValidator {
	return &CouponValidator{
		now: time.Now,
	}
}

func (c *CouponValidator) IsExpired(coupon *Coupon) bool {
	return c.now().After(coupon.expiredAt)
}
```
