How to deal with null values in DDD?

How can values be null? They have not been created yet, especially for aggregate root entities.

Use null object pattern. https://blog.ndepend.com/null-evil/


Alternatively, we can introduce a valid method (r.g. isZero for golang), but this defeats the purpose of always valid object. Entity should be created in its valid state. Alternatively, just set it to null and return a Boolean.
