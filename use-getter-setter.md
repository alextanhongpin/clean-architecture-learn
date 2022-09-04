# Use getter/setter (or wither)

In golang, struct fields that starts with lowercase are private. Most people seem to associate having private fields are necessary for domain models, though this is not related to DDD, it is purely OOP paradigm.

Getters in golang does not require the `Get` prefix [^1].

- You do not need to expose a getter/setter pair for each private fields - it is not idiomatic golang. 
- Also, the more fields you have, the longer your constructor becomes - for every private field, the only way to set the initial value is through constructor (for required fields). Setting it through setters makes it possible to forget setting required fields.
- For getter, avoid passing in arguments to compute value - move such logic to the domain service instead.
- Use getter/setter only for fields that you want to protect against invariant. There are no reason to put getter/setter for every fields.
- Separating read-model and write-model or creating another separate model will make it easier to distinguish between what fields to keep private or public


In this example, we split the user model into two separate model - one is purely read-only model without any behaviours, while the other contains credential logic. This avoid leaking the password to the read layer:
```go
type User struct {
	ID        string
	Name      string
	Email     string
	Birthdate time.Time
}

type UserCredential struct {
	id                string
	encryptedPassword string
}
```

[^1]: https://go.dev/doc/effective_go#Getters
