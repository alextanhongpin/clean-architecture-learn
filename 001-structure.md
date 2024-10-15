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


To accomplish a feature, we need to implement a service, which requires a store.

## Service

A service is simply a function. This function can be composed of multiple steps in order to complete the task at hand.

For example, calling an external service or making calls to the database.

There is an input and output. 

The most important thing is we have proper metrics for the service we call.

## Store

Aka repository, we just need a single repository for all services. This makes perfect sense, since we only have one database per app anyway.

We can also combine multiple data connector into one struct 

