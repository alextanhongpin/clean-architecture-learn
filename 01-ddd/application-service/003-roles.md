# Roles

Once you have a usecase, we can further breakdown the roles in the usacase.

The common pattern is communicating with a repository.

Each usecase always starts with similar logic

- validation
- calls repository
- transform output

Within a usecase, we can repeat this steps as many times to achieve a result.

However, it makes testing harder because we need to cover more path.


How about cross cutting concerns?

- logging
- tracing
- message queue
- errors

Ideally they dont belong here.




