# Authorization

- it is not part of the core domain logic
- which layer does authorization belongs to? I would put it at the outermost layer (aka controller layer). The logic for generating and validating the jwt token should be external to the application, and not even part of the _application service_. They can be explicitly passed to the _application service_ through _DTO_ though.
