# Clean Architecture Learning Repository

A comprehensive guide to clean architecture, domain-driven design, and software engineering best practices.

## 📚 Repository Structure

### 🔧 [fundamentals/](fundamentals/)
Core concepts and building blocks of clean architecture:
- **Aggregate**: Domain aggregate patterns and design
- **Entity**: Domain entity principles and implementation
- **Repository**: Data access patterns and best practices  
- **Domain Service**: Business logic coordination
- **Application Service**: Use case orchestration

### 🏗️ [patterns/](patterns/)
Advanced architectural patterns and practices:
- **Architecture Patterns**: Feature-oriented design, hexagonal architecture
- **Validation Patterns**: Domain vs application validation strategies
- **Error Handling**: Error propagation and transformation patterns
- **Persistence Patterns**: Repository design and data access strategies

### 💻 [implementation/](implementation/)
Language and framework-specific implementations:
- **Golang**: Go-specific patterns and examples
- **Examples**: Real-world code samples and case studies

## 🎯 Learning Path

### Beginner Track
1. Start with [fundamentals/README.md](fundamentals/README.md)
2. Read Entity and Aggregate concepts
3. Understand Repository pattern basics
4. Practice with simple examples

### Intermediate Track  
1. Study Domain and Application Services
2. Learn validation and error handling patterns
3. Explore architectural patterns
4. Implement feature-oriented architecture

### Advanced Track
1. Master complex domain modeling
2. Study event-driven patterns
3. Learn microservices boundaries
4. Practice with complete examples

## 🎨 Design Principles

### Core Principles
- **Domain-Centric**: Business logic is the core of the application
- **Dependency Inversion**: Depend on abstractions, not concretions  
- **Separation of Concerns**: Each layer has distinct responsibilities
- **Testability**: Design enables comprehensive testing

### Architectural Goals
- **Maintainability**: Easy to modify and extend
- **Clarity**: Simple to understand and navigate
- **Performance**: Efficient and scalable
- **Reliability**: Robust error handling and recovery

## 🚀 Quick Start

1. **Understand the Basics**: Read [fundamentals/001-aggregate.md](fundamentals/001-aggregate.md)
2. **See It in Action**: Check [implementation/golang/](implementation/golang/)
3. **Go Deeper**: Explore advanced patterns in [patterns/](patterns/)

## 📖 Key Concepts Summary

### Domain Layer
- **Entities**: Objects with identity and lifecycle
- **Value Objects**: Immutable objects defined by their attributes
- **Aggregates**: Consistency boundaries around related entities
- **Domain Services**: Business logic that spans multiple entities

### Application Layer
- **Application Services**: Use case orchestration
- **Command Handlers**: Process commands from presentation layer
- **Query Handlers**: Handle read operations
- **Event Handlers**: Process domain events

### Infrastructure Layer
- **Repositories**: Data access implementations
- **External Services**: Third-party integrations
- **Persistence**: Database and storage concerns
- **Messaging**: Event publishing and handling

## 🛠️ Best Practices

### Code Organization
- Use feature-oriented structure for better cohesion
- Keep layers properly separated with clear boundaries
- Minimize cross-cutting dependencies
- Design for testability from the start

### Domain Modeling
- Start with understanding the business domain
- Use ubiquitous language consistently
- Keep aggregates small and focused
- Validate business rules in appropriate layers

### Testing Strategy
- Unit test domain logic thoroughly
- Integration test repository implementations
- Use test doubles for external dependencies
- Test error scenarios and edge cases

## 🔍 Common Pitfalls to Avoid

1. **Anemic Domain Models**: Putting all logic in services
2. **God Objects**: Creating overly large aggregates
3. **Layer Violations**: Bypassing architectural boundaries
4. **Over-Engineering**: Adding unnecessary complexity
5. **Infrastructure Coupling**: Binding domain to technical details

## 📚 Philosophy

### Pragmatic Approach
- **Follow patterns when they solve problems you're having**
- **Don't apply patterns blindly**
- **Start simple and grow iteratively**
- **Focus on business value over architectural purity**

### When to Apply Clean Architecture
- ✅ Domain logic is mixed with infrastructure concerns → Separate domain from infrastructure
- ✅ Complex business rules → Use domain services and aggregates
- ✅ Multiple data sources → Repository pattern helps
- ❌ Simple CRUD operations → You probably don't need a Domain Layer

## 🤝 Contributing

This repository is a living document. Contributions are welcome:
- Fix typos and improve clarity
- Add new examples and patterns
- Share real-world experiences
- Suggest better organizations

---

**Note**: This repository focuses on practical implementation over theoretical purity. The goal is to provide actionable guidance for building maintainable software systems.
- Aggregates

Application layer (middle)
- Application service, transaction happens here
- Repositories
- Message brokers, adapters
- Workers
- Authorization

Presentation (outermost)
- controllers (REST) or resolver (Graphql)
- Authentication 
- Serialiser, logging, tracing metrics

# Good practices

- DDD implementation should not be influenced by persistence and database
- avoid anaemic domain model, domain classes full of getters and setters, void of behaviours.
- not all logic are domain logic, e.g. field validations
- security (e.g. authentication) is not part of the core domain (it could be a generic/supporting domain), however, most of the time, the implementation can lie in the `ui`/`application service` layer, out of the domain model.


Questions
- do i need to define a type for each layer, e.g. usecase.User, repository.User, entity.User? Nope, https://softwareengineering.stackexchange.com/questions/303478/uncle-bobs-clean-architecture-an-entity-model-class-for-each-layer
- why not user service, repository and entity in user package, vs separate repository for each? Circular dependency


What is a Model?
- a representation of your domain object
- contains rich behaviour
- Contains reference to other model
- Attributes may be primitive or value object

What is an Entity?
- a subset of your model (may or may not mimic your model)
- The persisted view of your model
- The model is always translatable to an entity, and vice versa
- Attributes may be db specific (but not necessarily)


What is a repository?
- fetches model and aggregate from a remote store (external request, remote store, client)

What is an Aggregate?
- composed of one or more model
- For business logic that spans across multiple entities
- Can be returned from a repository
- Why not service? Service may not require the aggregate to be called together (user may forget)
(Show an example of an aggregate and non aggregate)
(Show a good aggregate vs bad aggregate)



## Complex Application

- https://www.alibabacloud.com/blog/how-to-code-complex-applications-core-java-technology-and-architecture_595506?spm=a2c65.11461447.0.0.2ae134ffVGA3P9


# References

[^1]: [StackOverflow: UseCase-Drive vs Domain-Driven](https://stackoverflow.com/questions/3173070/design-methodology-use-case-driven-vs-domain-driven)

where to put domain services
- https://stackoverflow.com/questions/48200345/ddd-where-should-logic-go-that-tests-the-existence-of-an-entity
- https://medium.com/geekculture/an-introduction-to-domain-driven-design-f29fc1877e2a
- https://stackoverflow.com/questions/48200345/ddd-where-should-logic-go-that-tests-the-existence-of-an-entity
