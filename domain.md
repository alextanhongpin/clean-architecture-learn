
# Domain

- subdomain is part of business
- __core domains__ are where the money is
- __supporting domains__ support your core business (e.g. identity service, authentication)
- __generic domains__ are the ones you need, but don’t care a lot about, so you could probably buy them of the shelf, e.g. upload image to Amazon S3, even authentication, mailing payment provider can be generic subdomain
- how should we represent domain models in golang? Should it just be an interface (akin to abstract class) or should it be the implementation? Implementation or abstraction?

# References

1. [Decompose by Subdomain](https://microservices.io/patterns/decomposition/decompose-by-subdomain.html)
