# Third-party Integration

There are some questions when dealing with third-party integration:

- Where do we place the implementation? Definitely not in the domain layer, it will most likely be in the infrastructure layer, alongside repository implementation. External dependencies should not influence your domain model. On the other hand, the domain owns the interface, and decides the number and signature of the method. In short, domain defines interface, the third-party layer at infrastructure the implementation.
- What if the domain and third-party share the same entity (e.g. Payment)? Look at _Context Mapping_ such as Anti-corruption layer. We can probably separate the types, third-party types stays at the infrastructure layer, and domain type stays in domain. But ideally (?) we want our types to adhere to the domain model, as the domain should not be modelled on external dependencies (counter-argument: external dependencies may change, which will break the integration at the domain layer too). The solution for this is either to implement the _Anti-corruption_ layer, of the _Conformist_ pattern, but not both.
