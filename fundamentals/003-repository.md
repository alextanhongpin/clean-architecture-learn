# Repository Pattern

## Introduction

Repository is the layer that requests external data and returns a domain entity. External data is not necessarily related to persistence - it can be external API calls, cache, or storage objects. Repository connects multiple storage and client layers and returns domain entities.

## Core Concepts

### What is a Repository?
- **Definition**: A layer that abstracts data access and returns domain entities
- **Purpose**: Encapsulates the logic needed to access data sources
- **Scope**: Not limited to databases - includes APIs, caches, message queues, etc.

### Why Use Repository Pattern?
- **Separation of Concerns**: Separates the implementation of the persistence layer from business logic
- **Testing**: Makes switching ORMs easier (e.g. in golang, go-pg has been superseded by bun)
- **Abstraction**: Hides complexity of data access from the domain layer
- **Note**: The goal is not to help make switching database implementation easier, but it does make switching ORMs easier

## Input and Output

### Input
- Domain entities or primitives
- Query parameters and filters
- Transaction contexts

### Output
- Domain entities and aggregates
- Should return domain types, not database-specific types
- Maps external data sources to domain models

## Repository Design Principles

### Single Repository Approach
**Recommended**: Instead of creating many specific repositories, consolidate them under a single repository.

```go
// Good: Single repository
type Repository struct {
    db *sql.DB
}

func (r *Repository) CreateUser(ctx context.Context, user User) error {}
func (r *Repository) CreateProduct(ctx context.Context, product Product) error {}

// Bad: Multiple specific repositories
type UserRepository struct {}
type ProductRepository struct {}
```

**Advantages of Single Repository:**
- Simple construction
- The caller may need to access multiple resources, and using a single repository allows that
- We can define interface on what method is allowed
- Makes testing easier - you don't need to inherit all the other methods that are not used

### Repository per Usecase
For large applications, consider one repository per usecase:

```go
type AuthUsecase struct {
    repo authRepository // Contains only methods needed for auth
}
```

**Benefits:**
- **Interface Segregation**: Repository interface only exposes methods needed by that usecase
- **Easier Testing**: Smaller, more focused interfaces are easier to mock
- **Reduced Coupling**: Usecase only depends on methods it actually uses

### Aggregate-Focused Design
- Repositories should be designed around aggregates, not individual entities
- A repository can return multiple related entities as an aggregate
- Repository methods should return aggregates when write operations require business rule validation
- For read-only operations, returning individual entities may be more efficient

## Implementation Guidelines

### Business Logic Boundaries
- **Repository Responsibilities**:
  - Setting up transactions
  - Handling transaction context
  - Converting database errors to domain errors (e.g. not found, unique constraints)
  - Converting database types to domain entities
  - Mapping params to database params

- **Not Repository Responsibilities**:
  - Primary business logic (belongs in domain layer)
  - Complex business rule validation (belongs in domain services or entities)

### Error Handling
Repository should transform errors from the persistence layer to domain errors:
- Convert `sql.ErrNoRows` to domain-specific errors like `UserNotFound`
- Handle constraint violations and convert to appropriate domain errors
- Maintain abstraction between database and domain layers

### Transaction Management
- Transactions should be started at the application service layer, not repository layer
- Repository methods should accept transaction context when needed
- Avoid nested transactions by keeping transaction control at the service layer

## Layer Relationships

### Where Repository Belongs
- **Interface**: Domain layer (defines the contract)
- **Implementation**: Between application service and domain model, or in persistence layer

### Who Calls Repository?
- **Primary**: Application services
- **Secondary**: Domain services (in rare cases, read-only operations)
- **Never**: Entities or value objects

### Repository and Other Patterns
- **Factory vs Repository**: Factory handles the beginning of an object's life; Repository helps manage the middle and the end
- **Domain Service**: Repository can be injected into domain services in rare scenarios, but should avoid mutations

## Data Access Patterns

### Stores Pattern
Repository should be a facade that connects multiple sources:

```go
// Repository is not equal to ORM or raw SQL queries
type Repository struct {
    userStore    UserStore
    accountStore AccountStore
    cache        Cache
    messageQueue MessageQueue
}
```

**Benefits:**
- Repository is not the direct representation of the datasource
- Easier to abstract ORM and table-specific queries into stores
- Stores are implementation details, one repository can have multiple stores

### Partial Updates
Repository should support both full and partial updates:

```go
// Full entity update
func (r *Repository) UpdateUser(ctx context.Context, user User) error

// Partial update with specific fields
func (r *Repository) UpdateUserEmail(ctx context.Context, userID string, email string) error
```

**Considerations:**
- Full updates: Load entire entity, modify, then save
- Partial updates: More efficient, update only changed fields
- Consider using event sourcing patterns for complex update scenarios

## Best Practices

### 1. Aggregate Fetching
Instead of many small repository methods, handle related data in single methods:

```go
// Good: Return aggregate
func (r *Repository) GetUserWithAccounts(ctx context.Context, userID string) (*UserAggregate, error)

// Avoid: Multiple separate calls
func (r *Repository) GetUser(ctx context.Context, userID string) (*User, error)
func (r *Repository) GetAccountsByUserID(ctx context.Context, userID string) ([]Account, error)
```

### 2. Interface Segregation
Define repository interfaces specific to usecase needs:

```go
type AuthRepository interface {
    CreateUser(ctx context.Context, user User) error
    FindUserByEmail(ctx context.Context, email string) (*User, error)
    // Only methods needed for auth usecase
}
```

### 3. Separation of Concerns
- **Repository Layer**: Maps external data to domain entities
- **Store/Database Layer**: Handles database-specific operations
- **Domain Layer**: Contains business logic and rules

### 4. Testing Strategy
- Repository interfaces should be easily mockable
- Focus testing on the mapping logic between external data and domain entities
- Use integration tests for database-specific logic in stores

## Common Pitfalls

1. **Table-per-Repository**: Don't create one repository per database table
2. **Business Logic in Repository**: Keep business logic in domain layer
3. **Database Leakage**: Don't expose database-specific types in repository interfaces
4. **Over-abstraction**: Don't abstract things that don't need abstraction
5. **Transaction Coupling**: Don't start transactions in repository methods

## References

1. [Microsoft DDD CQRS Pattern: Infrastructure Persistence Layer Design](https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/infrastructure-persistence-layer-design)
2. [DDD Repository and Factory](https://stackoverflow.com/questions/31528368/ddd-repository-and-factory)
3. [Which layer do DDD repositories belong to?](https://softwareengineering.stackexchange.com/questions/396151/which-layer-do-ddd-repositories-belong-to)
4. [Can Domain Service Access Repositories](https://groups.google.com/g/dddcqrs/c/66zbcL97ilk?pli=1)
