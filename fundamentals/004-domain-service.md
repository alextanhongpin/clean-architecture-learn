# Domain Service

## Introduction

Domain services are part of the domain layer that encapsulate business logic which doesn't naturally fit within a single entity or aggregate. They operate on domain objects and express domain concepts that span multiple entities.

## Key Characteristics

- **Part of Domain Layer**: Belongs to the domain layer, not application layer
- **Stateless**: Unlike entities, domain services are stateless
- **Business Logic Focus**: Contains business logic that doesn't fit in a single entity
- **Pure Functions**: Should ideally be modeled as pure functions
- **Collaboration**: Deals with operations that require multiple domain objects

## When to Use Domain Services

Domain services should be used when:
- Business logic spans multiple entities or aggregates
- The logic doesn't naturally belong to any single entity
- Complex business rules require coordination between domain objects
- Operations that make sense only when done with other collaborators

### Examples
- **Checkout Process**: Involves Product (price), Purchase (items, payment method), and Payment validation
- **Order Validation**: Checking order limits across multiple orders and user accounts
- **Price Calculation**: Computing discounts based on customer type, product category, and current promotions

## Design Principles

### Avoid Anemic Domain Models
- Domain services should complement rich domain models, not replace them
- Don't extract all business logic into services, leaving entities as mere data containers
- Keep entity-specific logic within entities

### Naming and Structure
- Avoid generic names like "UserService" or "OrderService"
- Use intention-revealing names that describe the specific business operation
- Remove "Service" suffix when possible - use verb-based names
- Can be implemented as classes for namespacing or as pure functions

### Dependencies and Coupling
- Should have as few dependencies as possible
- If dependencies are required, pass them as interfaces, not concrete implementations
- Prefer passing required data as parameters rather than injecting dependencies
- Domain services should be unaware of infrastructure concerns

## Relationship with Other Layers

### Domain Service vs Application Service
Domain services are often confused with application services, but they serve different purposes:

**Domain Service:**
- Contains core business logic
- Operates on domain entities and value objects
- Expressed in terms of ubiquitous language
- Input/output are domain types
- Stateless business operations

**Application Service:**
- Orchestrates use cases
- Coordinates between domain services, repositories, and infrastructure
- Handles cross-cutting concerns (transactions, security, logging)
- Input/output are DTOs, not domain objects

### Repository Access
**General Rule**: Domain services should avoid accessing repositories directly.

**Preferred Approach:**
1. Application service queries repository for required entities
2. Application service passes entities to domain service
3. Domain service performs business logic on provided entities

**Exception**: In rare cases, domain services may access repositories for read-only operations, but should never perform mutations.

```go
// Preferred: Application service handles data access
func (as *OrderApplicationService) ProcessOrder(ctx context.Context, orderID string) error {
    order, err := as.repo.GetOrder(ctx, orderID)
    if err != nil {
        return err
    }
    
    customer, err := as.repo.GetCustomer(ctx, order.CustomerID)
    if err != nil {
        return err
    }
    
    // Domain service works with provided entities
    return CanProcessOrder(order, customer)
}

// Domain service as pure function
func CanProcessOrder(order *Order, customer *Customer) error {
    // Business logic here
}
```

## Implementation Patterns

### Pure Functions
For simple business logic without external dependencies:

```go
func CalculateShippingCost(order *Order, customer *Customer, shippingRules ShippingRules) Money {
    // Business logic implementation
}
```

### Service Objects
For more complex scenarios requiring configuration or multiple related operations:

```go
type PricingService struct {
    // Configuration, not external dependencies
}

func (ps *PricingService) CalculatePrice(product *Product, customer *Customer) Money {
    // Implementation
}

func (ps *PricingService) ApplyDiscount(price Money, discount *Discount) Money {
    // Implementation  
}
```

### Actor-Based Naming
Use domain actors/roles as service names when appropriate:

```go
type Librarian struct{}

func (l *Librarian) OrganizeBooks(books []Book) []Book {
    // Implementation
}

func (l *Librarian) AddToCatalog(book Book) error {
    // Implementation
}
```

## Return Types and Error Handling

### Boolean vs Error Returns
- **Boolean**: Use for simple yes/no operations where the reason for failure is obvious
- **Error**: Use when multiple failure reasons are possible or when detailed error information is needed

```go
// Simple boolean return
func IsEligibleForDiscount(customer *Customer) bool {
    return customer.IsVIP() && customer.HasActiveAccount()
}

// Error return for detailed feedback
func ValidateOrderLimits(customer *Customer, order *Order) error {
    if !customer.IsActive() {
        return ErrCustomerInactive
    }
    if order.ExceedsLimit(customer.GetOrderLimit()) {
        return ErrOrderLimitExceeded
    }
    return nil
}
```

### Chaining Operations
Domain services can chain multiple domain model methods and return specific errors:
- Domain entities should return booleans and have single responsibility
- Domain services can orchestrate multiple entity operations and provide detailed error handling

## Best Practices

1. **Keep It Simple**: Domain services should focus on business logic, not technical concerns
2. **Stateless Design**: Domain services should not maintain state between operations
3. **Clear Boundaries**: Don't put infrastructure concerns in domain services
4. **Testability**: Design for easy unit testing with minimal mocking
5. **Ubiquitous Language**: Use domain terminology in service names and methods
6. **Single Responsibility**: Each domain service should have a focused, well-defined purpose

## Common Pitfalls

1. **Infrastructure Coupling**: Avoid coupling domain services to databases, external APIs, or frameworks
2. **Application Logic**: Don't put application orchestration logic in domain services
3. **Entity Replacement**: Don't use domain services as a substitute for rich domain entities
4. **Over-Service**: Not every operation needs a domain service - prefer entity methods when possible

## References

1. [Can Domain Services access Repositories?](https://stackoverflow.com/questions/26930131/can-domain-services-access-repositories)
2. [Domain Driven Design- Tactical Patterns](http://domaindrivendesigns.blogspot.com/2018/11/domain-services-domain-model-domain.html?m=0)
3. [DDD Samples](https://github.com/VaughnVernon/IDDD_Samples)
