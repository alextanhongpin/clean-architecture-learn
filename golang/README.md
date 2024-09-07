# Clean architecture Golang

Demonstrates how to implement a whole project based on clean architecture using mostly standard packages.

## Structure

### Config

Application config should be loaded through environment variables only. You don't need a package manager like `viper` to load the environment variables.

A `makefile` and a `.env` goes a long way. See the guide here on changing environment.

#### Loading

```go
// Tips: Load the full url instead if partial parts like DB_USER, DB_PASS... to reduce nunber of environment variables.
var DatabaseURL = env.Load[string]("DATABASE_URL")
```

Bad:
```
REDIS_HOST=
REDIS_PORT=
```

Good:
```
REDIS_URL=
```
Sample implementation here:

https://github.com/alextanhongpin/go-clean-arch
