# Types

In application, it is common to have global reusable data types that are domain independent. We know that in DDD (domain-driven-architecture), it is common to define `value object` for types that are domain specify, yet has no identity (e.g. date of birth, age, sex).

However, most of the time, there are data types that are not native to a language, but is repeatedly used across the application. E.g. datetime interval, range type, set, money, bigint.

In golang, we can probably put them in `internal/types` to indicate that they are global reusable types that are domain independent. Other types and methods that could be included are

- slice operations (e.g. map, filter)
- group operation (e.g. group by etc)
- contextkey (to define keys for custom context)
- error (custom error types, as opposed to native error types)
