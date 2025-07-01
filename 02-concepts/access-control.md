# Designing access-control for applications

Context:
You need to restrict user's access to specific resources or operations based on certain criteria. Access control is not part of domain-driven design.

Some applications benefits greatly from proper access-control
- admin-facing application
- finance application
- app for managing communities/groups



## Interface

```go
type AccessControl interface {
  Allow(actor Actor, op Operation, resource Resource) (allow bool, reason error)
}
```

# References

- https://softwareengineering.stackexchange.com/questions/367659/implementing-ddd-users-and-permissions
- https://stackoverflow.com/questions/23464697/access-control-in-domain-driven-design/23485141#23485141
- https://softwareengineering.stackexchange.com/questions/71847/with-all-of-these-services-how-can-i-not-be-anemic/71883#71883
