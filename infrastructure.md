# Infrastructure Layer

- a.k.a secondary adapter from Hexagonal architecture
- used to abstract technical concerns
- contains implementation of the interfaces from the domain layer, such as repository implementation. Should repository implementation be in application service? See [1] `Since adapter code is code related to connecting an application to the world outside, adapter code is infrastructure code and should therefore reside in the infrastructure layer`
- any domain concepts (that includes domain services and repositories) that depends on external resources are defined by interfaces, and are implemented here
- how is the infrastructure layer structured? See [1]

```
src/
    <BoundedContext>/
        Domain/
            Model/
        Application/
        Infrastructure/
            <Port>/
                <Adapter>/
                <Adapter>/
                ...
            <Port>/
                <Adapter>/
                <Adapter>/
                ...
            ...
    <BoundedContext>/
        ...
```

# References

1. [Layers, ports & adapters - Part 3, Ports & Adapters](https://matthiasnoback.nl/2017/08/layers-ports-and-adapters-part-3-ports-and-adapters/)
