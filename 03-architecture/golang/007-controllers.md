# Controllers

Controllers bridges the presentation layer and usecase.

Controllers and usecases are application services. However, controllers are specific for http requests.

The reason for splitting the usecase into separate layer is because it can be used by other transport layers, such as grpc, cli, graphql, background jobs, cron etc.

By having separate usecase, we avoid scattering application logic.

Each transport layer's responsibility is just to map the types that can be consumed by usecase.
