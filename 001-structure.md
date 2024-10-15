# How to structure your application 

The short answer is, it depends. We cannot have a certain answer without a proper goal. Clean architecture by itself does not have any goals, but is ofetn perceived as clean.

We want to avoid fancy architecture such as clean, onion, or hexagonal. Instead, we want to focus on the basics

- the structure is simple and straightforward
- it is easy to reproduce in a glance
- no questions asked
- no indirection, no confusion
- metric driven
- almost flat, no deeply nested structure, less layers and imports

## Rules

We start with a feature-oriented architecture.

We avoid jargons like usecase, domain driven design etc, and instead pick feature. A feature is simple to understand. We for example want to implement authentication feature, that includes login and register. Or amybe a payment system.


A feature is there to help users accomplish a goal. If you have a sentence that starts with _I want to do sth_, you have a feature.


To accomplish a feature, we need to implement an _action_, which may communicate with a _store_.

## Action

Action is just a verb, e.g. _login_, _register_. An action can be broken down into multiple smaller steps, like validating the request, then reading from a store and writing to a store.

An action has many other properties, like idempotent, retryable etc and may have side effects.

There is an input and output. 

The most important thing is we have proper metrics for the action we call. 

- request - how frequent do we call this
- error - what is the error rate
- duration - how long does it take

The metrics measured becomes powerful as we compare different version of the same action, or across a set of actions.

When comparing different version of the same action
- did it become slower
- did the error rate increase
- are the still calls to the older action

Whe comparing across a set of actions
- which has the highest frequency (aka which is more adopted)
- which is having error

## Store

Aka repository, we just need a single repository for all actions. This makes perfect sense, since we only have one database per app anyway.

We can also combine multiple data connector into one struct:

```go
// Store represents data connector to multiple sources.
type Store struct {
  *Postgres
  *Kafka
  *Redis
}
```

Each action then declares the method interface that it requires from the store.

## File structure 

- for each feature, we define a folder
- each folder will have the CASE, controller, actions, store and entity
- this form of vertical slicing is really important, so that we don't have unnecessary imports and type mappings that exists when placed in different layer


```
register_controller.go
register.go
login_controller.go
login.go
store.go // Shared store
user.go // Since entity can be shared across different actions, we create a file by entity name. This entity only exists for this feature.
```

