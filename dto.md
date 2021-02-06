# DTO

- data transfer object
- plain data structure without any business logic (repeat, not a single business logic)
- factories may turn DTO into Entity
- usually used by the application service as an input/output contract
- input entity are specific to the application service's method, but the output may be shared, e.g. `CreateUserDto` for `createUser` method, but `PartialUser` as output dto.
