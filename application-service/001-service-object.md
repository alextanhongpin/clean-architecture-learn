# Service Object Pattern


## Context

One of the common issues when designing the application service layer is it often composed as a large transaction script.


This has several implications:

- the resulting code is hard to test, and needs to be executed from start to end
- the individual steps will be tested as a whole, either by changing the input arguments, or mocking the internal dependencies to return certain payload that would trigger the scenarios to be tested

## Decisions

There are few options to refactor a large application logic

1. use domain objects
2. decompose into service object

The former should be attempted before moving with the latter.

The latter is preferred if the usecase itself has a lot of intermediate steps, which would have its own dependencies too.

Todo: show example below


## Consequences

Each steps can be tested independently.
