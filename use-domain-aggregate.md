# Use domain gateway


In DDD, domain logic may sometimes leak at the repository layer. This is however, acceptable because only the repository layer can perform certain computation (e.g. performing aggregation of data, counting total records).

However, applying back the data to the domain aggregate/model/service tends to be unnatural.

The example below demonstrates an order checkout logic that checks the previous count of orders made for eligibility:


```go
package main

import (
	"context"
	"errors"
	"fmt"
)

var ErrOrderExceedLimit = errors.New("order: limit exceeded")

const OrderLimit int64 = 10

type Order struct {
	ID string
}

func (o *Order) IsWithinOrderLimit(orderCount int64) bool {
	return orderCount < OrderLimit
}

type orderRepository interface {
	FindOrderByUserID(ctx context.Context, userID int64) (*Order, error)
	FindOrderCountByOrderIDForUser(ctx context.Context, userID string) (map[string]int64, error)
}

type OrderUsecase struct {
	repo orderRepository
}

func (uc *OrderUsecase) Checkout(ctx context.Context, userID int64) error {
	order, err := uc.repo.FindOrderByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find order by user id: %w", err)
	}

	orderCountByOrderID, err := uc.repo.FindOrderCountByOrderIDForUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find order count by user id: %w", err)
	}

	if err := CanCheckoutOrder(order, orderCountByOrderID); err != nil {
		return err
	}

	return nil
}

func CanCheckoutOrder(order *Order, orderCountByOrderID map[string]int64) error {
	if !order.IsWithinOrderLimit(orderCountByOrderID[order.ID]) {
		return ErrOrderExceedLimit
	}

	return nil
}
```


There are a few problem with the code above:
- The repository is returning a `map[string]int`, which is not an aggregate. This makes the code very procedural.
- The domain model has a method that checks the if the order is within an order limit. One good rule of thumb for business logic in domain is the method should not depend on external data. 



A refactor using proper domain model:

```go
package main

import (
	"context"
	"errors"
	"fmt"
)

type OrderLimit struct {
	orderCountByOrderID map[string]int64
}

func (ol *OrderLimit) MustBeWithinLimit(o *Order) error {
	orderCount := ol.orderCountByOrderID[o.ID]
	if orderCount < OrderLimitPerUser {
		return nil
	}

	return ErrOrderExceedLimit
}

var ErrOrderExceedLimit = errors.New("order: limit exceeded")

const OrderLimitPerUser int64 = 10

type Order struct {
	ID string
}

type orderRepository interface {
	FindOrderByUserID(ctx context.Context, userID int64) (*Order, error)
	FindOrderLimitByUserID(ctx context.Context, userID int64) (*OrderLimit, error)
}

type OrderUsecase struct {
	repo orderRepository
}

func (uc *OrderUsecase) Checkout(ctx context.Context, userID int64) error {
	order, err := uc.repo.FindOrderByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find order by user id: %w", err)
	}

	orderLimit, err := uc.repo.FindOrderLimitByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find order limit by user id: %w", err)
	}

	if err := CanCheckoutOrder(order, orderLimit); err != nil {
		return err
	}

	return nil
}

func CanCheckoutOrder(order *Order, orderLimit *OrderLimit) error {
	return orderLimit.MustBeWithinLimit(order)
}
```
