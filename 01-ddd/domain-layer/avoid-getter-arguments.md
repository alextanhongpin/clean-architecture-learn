# Avoid Getter Arguments


Having a getter that relies on external arguments to compute a value is an anti-pattern:


```go
type Subscription struct {
	isPremiumContent bool
}

func (s *Subscription) CanSubscribe(isPremiumAccount bool) bool {
	if !s.isPremiumContent {
		return true
	}

	return s.isPremiumContent && isPremiumAccount
}
```

The example above also violates the domain principle, which is to avoid putting authorization logic in domain model. Getters should be able to derive all computed values from itself without relying on external input.

Any methods on the domain model that is not a setter and requires external input to compute values should be moved to the domain service instead.


```go
type Subscription struct {
	IsPremiumContent bool
}

type UserRole struct {
	IsActiveMembership bool
}

// NOTE: We are not calling this SubscriptionPermission, but rather a more specific name that describes the service responsibility.
// Domain service should be stateless
type SubscriptionPermissionChecker struct{}

func (s *SubscriptionPermissionChecker) CanSubscribe(subscription *Subscription, userRole UserRole) bool {
	if !subscription.IsPremiumContent {
		return true
	}
	return subscription.IsPremiumContent && userRole.IsActiveMembership
}
```

We now split the authorization logic into another separate layer. Since the method has no other dependencies, we can just use pure function:

```go
func CanUserRoleSubscribeToSubscription(subscription *Subscription, userRole UserRole) bool {
	if !subscription.IsPremiumContent {
		return true
	}
	return subscription.IsPremiumContent && userRole.IsActiveMembership
}
```
