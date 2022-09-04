# Use Separate Read/Write model

Separating read-model and write-model or creating another separate model will make it easier to distinguish between what fields to keep private or public


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

