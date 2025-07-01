# Unclean architecture 

Are there such thing as unclean architecture? Definitely.

Any attempts to create a clean architecture without actually understanding the need of it often leads to unclean architecture.

Architecture is rarely language agnostic. Sure, the fundamentals may be the same, but the actual implementation differs a lot by language. A clean architecture implementation in Java or Dotnet is very different from javascript or golang.

Golang specifically is quite unique in the sense that its way of defining packages can lead to "gradient" layers, where there is a tight coupling between layers even when that is not the intent in clean architecture.

It is hard to create fully isolated layers, since types needs to be imported in order to be used.

This guide serves as a guideline, since I have seen many unnecessarily poor design of clean architecture in golang. Golang is simple, and we should keep it that away.

Avoid over abstraction.

We will cover the different layers 
- presentation
- domain
- usecase
- adapter

and suggestions on how to apply them in golang.

We also suggest not adding anything that doesn't make sense.

The less code, the better.

