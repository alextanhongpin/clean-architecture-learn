# Working with Features to Usecases



When working on a new feature, we start with a feature definition:


```
Feature: Manage shipment
```

Then, we write a user story that describes the value:

```
As a user,
I want to get the cheapest shipment,
In order to reduce cost
```

From there, we can get into more detail usecase:


```
Scenario: User view shipment option

1. User checkout several items
2. System presents shipment option
Business rule:
- when cost is above $100, then the shipment should be free.
```

This usecases can then be translated to code:

```go
type CheckoutUsecase struct {}

func (uc *CheckoutUsecase) CheckOut(order Order) error {
	// Shipping fee must be 0 if cost is above 100.
	fee, err := CalculateShippingFee(order)
	if err != nil {
		return err
	}

	return nil
}
```
