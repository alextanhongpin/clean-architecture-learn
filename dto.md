# DTO

- data transfer object
- plain data structure without any business logic (repeat, not a single business logic)
- factories may turn DTO into Entity
- usually used by the application service as an input/output contract
- input entity are specific to the application service's method, but the output may be shared, e.g. `CreateUserDto` for `createUser` method, but `PartialUser` as output dto.


https://professionalbeginner.com/the-dto-dilemma
https://khalilstemmler.com/articles/typescript-domain-driven-design/repository-dto-mapper/
https://softwareengineering.stackexchange.com/questions/267612/when-is-it-appropriate-to-map-a-dto-back-to-its-entity-counterpart/277440
https://softwareengineering.stackexchange.com/questions/408385/domain-driven-design-updating-and-persisting-aggregates
https://buildplease.com/pages/repositories-dto/
