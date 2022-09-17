# Errors across layers

Different layer have different errors, but how do we translate them across different layer. The domain errors belongs in domain layers.

1. repository. They should return domain errors, so sql.ErrNoRows should be translated to the entity's error, such as UserNotFound
2. application service. Here, they also deal mostly with domain errors. However, the question is what errors to return to the caller, which is the presentation layer? Application service does not return domain models, so by right they should not return domain errors too
3. presentation layer, e.g. rest layer. In rest, errors are represented as status code. So errors such as 404 not found, 400 bad request should be mapped from the errors that are returned from the application service. However, what errors should hte application service return?



## Should each layer has it's own error handling?

Yes, ideally each layer should have it's own error. If the errors from a subpackage needs to be handled by the caller, then it should be mapped to the caller errors.


This is usually subject to how you structure your packages.


For example, errors in the `repository` packages could be mapped to domain layer in the usecase layer when needed. It might not be ideal to expose the `repository` errors to the `presentation` layer. If we need to do so, then the `usecase` layer that is calling the `repository` layer should be mapping the `repository` layer to the `domain` errors before passing it to the `repository` layer.


The reason why the `presentation` layer should only receive `domain` errors, versus `usecase` errors is that most invariant errors should belong to the domain layer, which could also be returned from the domain service, factory etc, as there should not be any domain logic in the usecase layer.

The errors translations etc should not belong in the usecase layer either. It should belong in the presentation layer, for example the REST API layer before returning the errors to the client.
