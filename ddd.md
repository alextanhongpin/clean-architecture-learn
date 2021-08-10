# Separate DDD from CRUD

Admin API designs are usually CRUD. the problems faced there are not DDD as you will most likely be mapping your data store closer to your db design, and they dont have much business logic nor will they follow the common ddd design, especially when dealing with bulk operations.


# Separate your read from mutation

Your read data could just be loading data directly from Repository, without the need for unnecessary layers. There are very few business logic pertaining to read, except mapping some columns to flags. After all, the UI should not have business logic. 

For example, if you are running an ecommerce Backend, you might want to return a flag available to indicate there are stock available. Or if there are discount, you might want to show the discounted price derived from the discount value.

For other usecases, it is also common to return a flag indicating you owned the entity, such as me (when the entity is related to a being) or mine if its other nouns.

You might also have more complex rules such as permission, or rule engine that will compute some value. 

But they all dont belong to the domain. Getters and computed value does not count as business value unless they are there to ensure invaraint.
