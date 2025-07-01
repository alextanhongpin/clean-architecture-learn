# Application Service

## Introduction

Application services act as the orchestration layer between the presentation layer and the domain layer. They coordinate business workflows, manage transactions, and handle cross-cutting concerns while delegating core business logic to domain entities and services.

## Key Responsibilities

### Primary Responsibilities
- **Use Case Orchestration**: Coordinate complex business workflows
- **Transaction Management**: Handle transaction boundaries and consistency
- **Security**: Implement authentication and authorization
- **Validation**: Perform application-level validation (distinct from domain validation)
- **Error Handling**: Transform domain errors into appropriate application responses
- **External Integration**: Coordinate with external services and APIs
- **Event Publishing**: Publish domain events for cross-bounded context communication

### What Application Services Should NOT Do
- **Business Logic**: Core business rules belong in domain layer
- **Persistence Logic**: Direct database access should be through repositories
- **Presentation Logic**: UI-specific concerns belong in presentation layer
- **Infrastructure Details**: Low-level technical concerns belong in infrastructure layer

## Input and Output

### Input Handling
- **Accept DTOs**: Receive data transfer objects, not domain entities
- **Validation**: Perform input validation and sanitation
- **Authentication**: Verify user identity and permissions
- **Command Processing**: Handle commands from presentation layer

### Output Handling
- **Return DTOs**: Never expose domain entities directly to presentation layer
- **Error Translation**: Convert domain errors to application-appropriate formats
- **Response Formatting**: Prepare data in format suitable for presentation layer
- **Status Reporting**: Provide operation status and metadata

## Design Patterns

### Command Handler Pattern
Application services often implement command handlers:

```go
type UserRegistrationCommand struct {
    Email    string
    Password string
    Name     string
}

type UserApplicationService struct {
    userRepo     UserRepository
    emailService EmailService
    eventBus     EventBus
}

func (s *UserApplicationService) RegisterUser(ctx context.Context, cmd UserRegistrationCommand) (*UserDTO, error) {
    // 1. Validate input
    if err := s.validateRegistrationCommand(cmd); err != nil {
        return nil, err
    }
    
    // 2. Check business rules via domain service
    if exists, err := s.userRepo.ExistsByEmail(ctx, cmd.Email); err != nil {
        return nil, err
    } else if exists {
        return nil, ErrEmailAlreadyExists
    }
    
    // 3. Create domain entity
    user, err := NewUser(cmd.Email, cmd.Name)
    if err != nil {
        return nil, err
    }
    
    // 4. Persist via repository
    if err := s.userRepo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // 5. Handle side effects
    if err := s.emailService.SendWelcomeEmail(ctx, user.Email()); err != nil {
        // Log error but don't fail the operation
        log.Error("Failed to send welcome email", err)
    }
    
    // 6. Publish events
    s.eventBus.Publish(UserRegisteredEvent{UserID: user.ID()})
    
    // 7. Return DTO
    return s.toUserDTO(user), nil
}
```

### Query Handler Pattern
For read operations, application services can implement query handlers:

```go
type GetUserQuery struct {
    UserID string
}

func (s *UserApplicationService) GetUser(ctx context.Context, query GetUserQuery) (*UserDTO, error) {
    user, err := s.userRepo.FindByID(ctx, query.UserID)
    if err != nil {
        return nil, err
    }
    
    return s.toUserDTO(user), nil
}
```

## Transaction Management

### Transaction Boundaries
Application services define transaction boundaries:

```go
func (s *OrderApplicationService) ProcessOrder(ctx context.Context, cmd ProcessOrderCommand) error {
    return s.repo.WithTransaction(ctx, func(txCtx context.Context) error {
        // All operations within this block are part of the same transaction
        
        order, err := s.orderRepo.FindByID(txCtx, cmd.OrderID)
        if err != nil {
            return err
        }
        
        if err := order.Process(); err != nil {
            return err
        }
        
        if err := s.orderRepo.Save(txCtx, order); err != nil {
            return err
        }
        
        // Update inventory
        return s.inventoryRepo.UpdateStock(txCtx, order.Items())
    })
}
```

### Unit of Work Pattern
For complex operations involving multiple aggregates:

```go
type UnitOfWork interface {
    RegisterNew(entity interface{})
    RegisterDirty(entity interface{})
    RegisterDeleted(entity interface{})
    Commit(ctx context.Context) error
    Rollback() error
}

func (s *ApplicationService) ComplexOperation(ctx context.Context) error {
    uow := s.NewUnitOfWork()
    defer uow.Rollback() // Rollback if not committed
    
    // Register changes
    uow.RegisterNew(newEntity)
    uow.RegisterDirty(modifiedEntity)
    uow.RegisterDeleted(deletedEntity)
    
    // Commit all changes atomically
    return uow.Commit(ctx)
}
```

## Validation Layers

### Application-Level Validation
- **Input Validation**: Format, required fields, data types
- **Authorization**: User permissions and access control
- **Business Rule Validation**: Rules that span multiple aggregates

### Domain-Level Validation
- **Entity Invariants**: Rules that belong to individual entities
- **Aggregate Consistency**: Rules within aggregate boundaries
- **Domain Rules**: Core business logic validation

```go
func (s *UserApplicationService) UpdateUserProfile(ctx context.Context, cmd UpdateProfileCommand) error {
    // Application-level validation
    if !s.authService.CanUpdateProfile(ctx, cmd.UserID) {
        return ErrUnauthorized
    }
    
    // Load domain entity
    user, err := s.userRepo.FindByID(ctx, cmd.UserID)
    if err != nil {
        return err
    }
    
    // Domain-level validation happens in entity method
    if err := user.UpdateProfile(cmd.Name, cmd.Email); err != nil {
        return err
    }
    
    return s.userRepo.Save(ctx, user)
}
```

## Error Handling Strategy

### Error Translation
Transform domain errors into application-appropriate formats:

```go
func (s *ApplicationService) translateError(err error) error {
    switch {
    case errors.Is(err, domain.ErrUserNotFound):
        return &ApplicationError{
            Code:    "USER_NOT_FOUND",
            Message: "User not found",
            Status:  404,
        }
    case errors.Is(err, domain.ErrInvalidEmail):
        return &ApplicationError{
            Code:    "INVALID_INPUT",
            Message: "Invalid email format",
            Status:  400,
        }
    default:
        return &ApplicationError{
            Code:    "INTERNAL_ERROR",
            Message: "Internal server error",
            Status:  500,
        }
    }
}
```

## Integration Patterns

### External Service Integration
Handle external service calls with proper error handling and fallbacks:

```go
func (s *ApplicationService) ProcessPayment(ctx context.Context, cmd PaymentCommand) error {
    // Domain validation
    order, err := s.orderRepo.FindByID(ctx, cmd.OrderID)
    if err != nil {
        return err
    }
    
    if !order.CanProcessPayment() {
        return ErrOrderNotReady
    }
    
    // External service call with retry and timeout
    payment, err := s.paymentService.ProcessPayment(ctx, PaymentRequest{
        Amount:   order.Total(),
        OrderID:  order.ID(),
        Method:   cmd.PaymentMethod,
    })
    if err != nil {
        return s.handlePaymentError(err)
    }
    
    // Update domain state
    order.MarkPaid(payment.ID)
    return s.orderRepo.Save(ctx, order)
}
```

### Event Publishing
Publish domain events for eventual consistency:

```go
func (s *ApplicationService) PublishEvents(ctx context.Context, events []DomainEvent) error {
    for _, event := range events {
        if err := s.eventBus.Publish(ctx, event); err != nil {
            // Log error and potentially retry
            log.Error("Failed to publish event", "event", event, "error", err)
            return err
        }
    }
    return nil
}
```

## Testing Strategy

### Unit Testing
Focus on testing the orchestration logic:

```go
func TestUserApplicationService_RegisterUser(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    mockEmailService := &MockEmailService{}
    service := NewUserApplicationService(mockRepo, mockEmailService)
    
    cmd := UserRegistrationCommand{
        Email:    "test@example.com",
        Password: "password123",
        Name:     "Test User",
    }
    
    // Act
    result, err := service.RegisterUser(context.Background(), cmd)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockRepo.AssertSaveCalled(t)
    mockEmailService.AssertWelcomeEmailSent(t)
}
```

### Integration Testing
Test the integration between layers:

```go
func TestUserApplicationService_Integration(t *testing.T) {
    // Use real repository with test database
    db := testdb.Setup(t)
    repo := postgres.NewUserRepository(db)
    service := NewUserApplicationService(repo, &MockEmailService{})
    
    // Test actual database interactions
    result, err := service.RegisterUser(context.Background(), validCommand)
    
    assert.NoError(t, err)
    // Verify data was actually persisted
    user, err := repo.FindByEmail(context.Background(), validCommand.Email)
    assert.NoError(t, err)
    assert.Equal(t, validCommand.Name, user.Name())
}
```

## Best Practices

1. **Thin Layer**: Keep application services thin - they should orchestrate, not implement business logic
2. **Single Responsibility**: Each application service method should handle one use case
3. **Dependency Injection**: Use dependency injection for repositories and external services
4. **Error Handling**: Implement comprehensive error handling and logging
5. **Transaction Management**: Clearly define transaction boundaries
6. **Security**: Implement proper authentication and authorization
7. **Monitoring**: Add metrics and monitoring for application service operations

## Common Pitfalls

1. **Fat Services**: Putting too much logic in application services
2. **Domain Logic Leakage**: Implementing business rules in application layer
3. **Direct Entity Exposure**: Returning domain entities instead of DTOs
4. **Transaction Abuse**: Creating overly long transactions
5. **Infrastructure Coupling**: Coupling application services to specific infrastructure

## References

1. [Domain-Driven Design: Tackling Complexity in the Heart of Software](https://www.amazon.com/Domain-Driven-Design-Tackling-Complexity-Software/dp/0321125215)
2. [Application Services in DDD](https://enterprisecraftsmanship.com/posts/application-services-or-controllers/)
3. [DDD and Hexagonal Architecture](https://vaadin.com/blog/ddd-part-3-domain-driven-design-and-the-hexagonal-architecture)
