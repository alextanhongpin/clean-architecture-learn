# Application Structure Patterns

## Feature-Oriented Architecture

Instead of layer-based architecture, organize code around business features.

### Core Principles

- **Simple and Straightforward**: Structure should be easy to understand at a glance
- **No Indirection**: Avoid unnecessary abstraction layers
- **Metric Driven**: Focus on measurable outcomes
- **Flat Structure**: Minimize deep nesting and complex layer hierarchies
- **Minimal Imports**: Reduce coupling between components

### Why Feature-Oriented?

Traditional layered architecture creates several problems:
- **Cross-cutting Changes**: Simple feature changes require modifications across multiple layers
- **Coupling**: Layers become tightly coupled through shared interfaces
- **Complexity**: Deep inheritance hierarchies make code hard to understand
- **Import Hell**: Complex dependency graphs between packages

Feature-oriented architecture solves these by:
- **Vertical Slicing**: All code for a feature lives together
- **Reduced Coupling**: Features are independent and self-contained  
- **Clear Boundaries**: Easy to understand what belongs where
- **Team Alignment**: Teams can own entire features

## Structure Definition

### Feature Definition
A feature helps users accomplish a goal. If you have a sentence that starts with _"I want to..."_, you have a feature.

Examples:
- **Authentication Feature**: Login, register, password reset
- **Payment Feature**: Process payments, refunds, billing
- **Notification Feature**: Send emails, SMS, push notifications

### Action-Based Components
Each feature contains **actions** - specific operations users can perform.

Actions have properties:
- **Idempotent**: Can be safely retried
- **Observable**: Proper metrics and logging
- **Testable**: Clear input/output contracts
- **Composable**: Can be combined with other actions

### Store Pattern
Instead of multiple repositories, use a single **Store** per feature that aggregates all data sources:

```go
// Store represents data connector to multiple sources
type Store struct {
    *Postgres
    *Redis  
    *Kafka
}
```

Benefits:
- **Single Entry Point**: One place for all data access
- **Simplified Testing**: Mock one interface instead of many
- **Consistent Transactions**: Easier to manage consistency across data sources

## File Organization

### Directory Structure
```
feature-name/
├── action1_controller.go    # HTTP/gRPC handlers
├── action1.go              # Business logic
├── action2_controller.go   
├── action2.go              
├── store.go                # Data access layer
├── entity.go               # Domain models
└── types.go                # Feature-specific types
```

### Example: Authentication Feature
```
auth/
├── register_controller.go  # Registration HTTP handlers
├── register.go             # Registration business logic  
├── login_controller.go     # Login HTTP handlers
├── login.go                # Login business logic
├── store.go                # User data access
├── user.go                 # User entity
└── types.go                # Auth-specific types
```

### Benefits of This Structure
- **Co-location**: Related code lives together
- **Independent Features**: Features don't depend on each other
- **Clear Ownership**: Easy to assign team ownership
- **Simplified Imports**: No complex cross-layer dependencies

## Implementation Guidelines

### 1. Action Design
Each action should:
- Accept input DTO
- Return output DTO or error
- Handle its own validation
- Log important events
- Emit metrics

```go
type LoginAction struct {
    store LoginStore
    logger Logger
    metrics Metrics
}

func (a *LoginAction) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
    // Validate input
    if err := req.Validate(); err != nil {
        a.metrics.InvalidRequest.Inc()
        return nil, err
    }
    
    // Business logic
    user, err := a.store.FindUserByEmail(ctx, req.Email)
    if err != nil {
        a.metrics.LoginError.Inc()
        return nil, err
    }
    
    if !user.ValidatePassword(req.Password) {
        a.metrics.InvalidCredentials.Inc()
        return nil, ErrInvalidCredentials
    }
    
    // Success metrics
    a.metrics.SuccessfulLogin.Inc()
    a.logger.Info("User logged in", "userID", user.ID)
    
    return &LoginResponse{
        Token: generateToken(user),
        User:  toUserDTO(user),
    }, nil
}
```

### 2. Store Interface Segregation
Each action declares only the store methods it needs:

```go
type LoginStore interface {
    FindUserByEmail(ctx context.Context, email string) (*User, error)
}

type RegisterStore interface {
    CreateUser(ctx context.Context, user *User) error
    ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// Actual store implements all interfaces
type AuthStore struct {
    db *sql.DB
}

func (s *AuthStore) FindUserByEmail(ctx context.Context, email string) (*User, error) { ... }
func (s *AuthStore) CreateUser(ctx context.Context, user *User) error { ... }
func (s *AuthStore) ExistsByEmail(ctx context.Context, email string) (bool, error) { ... }
```

### 3. Metrics Strategy
Track three key metrics for every action:
- **Request Rate**: How often is this action called?
- **Error Rate**: What percentage of requests fail?
- **Duration**: How long does the action take?

These metrics enable:
- **Performance Monitoring**: Identify slow operations
- **Reliability Tracking**: Monitor error rates
- **Usage Analytics**: Understand feature adoption
- **Capacity Planning**: Plan for scale based on usage

## Anti-Patterns to Avoid

### 1. Generic Services
```go
// Bad: Generic service
type UserService struct {}
func (s *UserService) DoEverything() {}

// Good: Specific actions
type RegisterUserAction struct {}
type LoginUserAction struct {}
```

### 2. Deep Layer Hierarchies
```go
// Bad: Too many layers
UserController -> UserService -> UserRepository -> UserStore -> Database

// Good: Direct path
UserController -> UserAction -> UserStore -> Database
```

### 3. Cross-Feature Dependencies
```go
// Bad: Features depending on each other
auth.LoginAction depends on profile.UserProfileService

// Good: Independent features communicating via events
auth.LoginAction emits UserLoggedInEvent
profile feature listens to UserLoggedInEvent
```

## Migration Strategy

### From Layer-Based to Feature-Based

1. **Identify Features**: Group related use cases
2. **Extract Vertically**: Move all related code to feature directory
3. **Consolidate Stores**: Combine related repositories into feature stores
4. **Simplify Interfaces**: Remove unnecessary abstraction layers
5. **Add Metrics**: Instrument actions with proper monitoring

### Gradual Transition
- Start with new features using the new structure
- Gradually refactor existing features
- Maintain both structures during transition
- Use metrics to validate improvements

## Conclusion

Feature-oriented architecture provides:
- **Clarity**: Easy to understand and navigate
- **Productivity**: Faster development and debugging
- **Maintainability**: Easier to modify and extend
- **Team Efficiency**: Clear ownership and boundaries
- **Observability**: Built-in metrics and monitoring

This approach prioritizes pragmatism over architectural purity, focusing on delivering business value efficiently.
