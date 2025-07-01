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


## Where to map other external client types?

Depending on your application, you may need to call external party API and map the domain types to and from your domain model. For example, a payment type may need to map external response back to your payment domain models. However, the questions lies, where to map the implementation in your app?

One can for example have the folder for the domain types. Then another adapter layer to map the types back to the domain model.


What happens if the types cannot be mapped back? We can supposedly import the types in the application service, it is not violating the principles.

# References

1. [Layers, ports & adapters - Part 3, Ports & Adapters](https://matthiasnoback.nl/2017/08/layers-ports-and-adapters-part-3-ports-and-adapters/)
