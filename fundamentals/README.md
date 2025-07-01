# Clean Architecture Fundamentals

This directory contains the core concepts and building blocks of clean architecture and domain-driven design.

## Table of Contents

1. [**Aggregate**](001-aggregate.md) - Understanding aggregate roots and their role in domain modeling
2. [**Entity**](002-entity.md) - Domain entities, their responsibilities, and design principles  
3. [**Repository**](003-repository.md) - Data access patterns and repository implementation guidelines
4. [**Domain Service**](004-domain-service.md) - Business logic that spans multiple entities
5. [**Application Service**](005-application-service.md) - Use case orchestration and coordination layer

## Learning Path

### Beginner
Start with understanding the basic building blocks:
1. **Entity** - Learn about domain entities and their encapsulation principles
2. **Aggregate** - Understand how entities work together in aggregates
3. **Repository** - Learn how to abstract data access

### Intermediate  
Understand the service layers:
4. **Domain Service** - Learn when and how to implement domain services
5. **Application Service** - Master use case orchestration

### Advanced
Apply these concepts in practice:
- See the `patterns/` directory for advanced implementation patterns
- Check `implementation/` for language-specific examples
- Review `examples/` for real-world scenarios

## Key Principles

- **Separation of Concerns**: Each layer has distinct responsibilities
- **Dependency Inversion**: High-level modules don't depend on low-level modules
- **Domain-Centric**: Business logic is the core of the application
- **Testability**: Design enables easy unit and integration testing

## Related Directories

- `../patterns/` - Advanced architectural patterns and practices
- `../implementation/` - Language and framework-specific implementations  
- `../examples/` - Real-world examples and case studies
