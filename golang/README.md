# Clean architecture Golang

Demonstrates how to implement a whole project based on clean architecture using mostly standard packages.

## Structure

### Config

Application config should be loaded through environment variables only. You don't need a package manager like `viper` to load the environment variables.

A `makefile` and a `.env` goes a long way. See the guide here on changing environment.

For golang:
1. prefer flat over nested config
2. load all config through variables, avoid struct or map because values may be forgotten or you have to add additional logic
3. all envvar declared cannot be empty, so panic on runtime
4. avoid storing complicated types (json, yaml, slices)
5. avoid third party packages ... stdlib has everything you need
6. don't need dotenv etc
7. keep environment variables small, merge parts when needed e.g. redis url instead of redis host+redis port
8. config values should not be exported
9. config contains the configuration for a dependency, not the values to be used. Config can build with all the necessary values from an environment variablea, but they can still accpet other deps through dependency injection, in case they need to be shared. For shared dependencies, create a root constructor, e.g. repository creator that creates multiple repositories.


For others:
- store environment in .env
- separate different environment with a suffix
- keep secrets away from environment
- use makefile to load environment instead of dotenv
- switch environment with makefile
- guard against unsafe environment like production by prompting or guarding


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

## Commands

1. makefile will suffice for most commands
  - go commands
  - dev setup
  - db migration setup
  - docker setup
  - infra access cli
  - run commands
  - test commands (healthcheck, pings)
  - open dashboard
  - install tools
  - build commands
  - semver bump
2. for stuff that needs to be done programmatically (e.g. scripts to update data, run server etc) place it under `cmd`

## Adapters

Most of your external packages belongs here. 

Adapter contains the implementation, not abstraction, aka dont redefine another constructor etc. You are not building any reusable package here. You initialized deps here.

They should be minimally defined with all configs configured for running the application. Tests environment should use hardcoded values to avoid accidentally hitting production environment.

Bad:

```go
postgres.New(user, pass, host, port, db)
```
Good:
```
pg := postgres.New()
```

If you need to create separate instances, define a private constructor to declare multiple instance:


```go
func NewAsset() *Bucket {
  return newBucket("assets")
}
func newBucket(name string) *Bucket {}
```

All dependencies should be clearly defined with zero customization.

Sample implementation here:

https://github.com/alextanhongpin/go-clean-arch
