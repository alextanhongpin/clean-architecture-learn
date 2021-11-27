
# Domain

- subdomain is part of business
- __core domains__ are where the money is
- __supporting domains__ support your core business (e.g. identity service, authentication)
- __generic domains__ are the ones you need, but don’t care a lot about, so you could probably buy them of the shelf, e.g. upload image to Amazon S3, even authentication, mailing payment provider can be generic subdomain
- how should we represent domain models in golang? Should it just be an interface (akin to abstract class) or should it be the implementation? Implementation or abstraction?


## Are domain useful?

- they are used to encapsulate business rules to avoid repetition. But if they are only used once, there is probably no need for it
- most complex business rules resides in aggregate layer. However, designing a good aggregate layer is not easy, as it can slowly become god object. Also, it hardly fulfils the singke responsibility principle, as aggregate may be composed of 5 entity, but only 2 is used for the business logic. In that case, it is probably better to create a pure function that accepts only the 2 entities required for that logic.


# References

1. [Decompose by Subdomain](https://microservices.io/patterns/decomposition/decompose-by-subdomain.html)
