# Errors across layers

Different layer have different errors, but how do we translate them across different layer. The domain errors belongs in domain layers.

1. repository. They should return domain errors, so sql.ErrNoRows should be translated to the entity's error, such as UserNotFound
2. application service. Here, they also deal mostly with domain errors. However, the question is what errors to return to the caller, which is the presentation layer? Application service does not return domain models, so by right they should not return domain errors too
3. presentation layer, e.g. rest layer. In rest, errors are represented as status code. So errors such as 404 not found, 400 bad request should be mapped from the errors that are returned from the application service. However, what errors should hte application service return?
