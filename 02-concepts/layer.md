# Layers in Clean Architecture

You have a clean layer if they are separated by a namespace, and requires importing with that namespace. The types that belongs to each layer should not leak into another layer.

This may vary across different programming languages, which implements their own module (or package) separately.

For example, in golang, this is not considered a layer:
```
yourapp/
- repository.go
- usecase.go
```

Both `repository.go` and `usecase.go` belongs to the same package, and hence there is no clear separation between them. A better representation would be:

```
yourapp/
- repository/main.go
- usecase/main.go
```

Now the `repository` and `usecase` is a separate layer.

Compare this with `nodejs`, where each file could be a possible layer by itself:

```
yourapp/
- repository.js
- usecase.js
```

Summary
- depending on the language you choose, you may end up with more folders for layer separation
- the term layer refers to a namespaced module, and each layer has a separate namespace

## Layer information

Should each layer has their own input/output types? 
- https://softwareengineering.stackexchange.com/questions/303478/uncle-bobs-clean-architecture-an-entity-model-class-for-each-layer
- https://www.ssw.com.au/rules/rules-to-better-clean-architecture
- https://softwareengineering.stackexchange.com/questions/303478/uncle-bobs-clean-architecture-an-entity-model-class-for-each-layer
- https://www.oncehub.com/blog/explaining-clean-architecture



good article on layering
https://blog.ploeh.dk/2012/02/09/IsLayeringWorththeMapping/

https://khalilstemmler.com/articles/software-design-architecture/domain-driven-design-vs-clean-architecture/
https://buildplease.com/pages/repositories-dto/
https://softwareengineering.stackexchange.com/questions/303478/uncle-bobs-clean-architecture-an-entity-model-class-for-each-layer/303480#303480

great 
https://discourse.world/h/2017/08/11/Misconceptions-Clean-Architecture

https://martinfowler.com/bliki/LocalDTO.html

https://stackoverflow.com/questions/21554977/should-services-always-return-dtos-or-can-they-also-return-domain-models

Having different types at each layer increases the conversion or mapping between dofferent types. Also, changes on the inner fields affects all outer layers.


Note
- one common misunderstanding of layers is that they tend to be separate package or module etc. Doing so could increase the number types conversion.
- we could instead place all of those related code in a single package. This is the concept known as vertical slicing


What happens as we move on the inner layer?

The outer layer implements policies of the inner layer implements rules.

As we go further in, the concerns switches mainly to persistence.

## Why layers fails

Layer are just human concerns. 

Layers are not independent. each layer usually calls another layer forming a gradient. That is why there can be overlap in responsibilities.

Splitting code into layers is just shifting code. if we think about it, the overall logic does not change. 


## How many layers do we need

That depends on the layer 

### Usecase

For usecase layer, we probably only need two layers:

- repository layer
- domain service layer

The repository layer here does not refer to database specific implementation. Instead, it is just a layer that is responsible for mapping any external types to a domain type. (todo: link to better articles).

The repository layer will be bridge (or is it adapter?) between different data source such as APIs, database query calls or even configuration etc. There are many benefits of hiding all external calls behind the repository facade, because in reality, many operations depends on the database layer. Separating them as individual dependencies means creating more utilities to map from one layer to another.

The domain service layer deserves extra attention. The goal is to eliminate inline logic, which usually results in hard to test code (link reference here). This also allows us to swap implementation, inspect and intercept the request and response with hooks. 
