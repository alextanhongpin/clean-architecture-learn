# Aggregate Root

- Aggregate root is a collection of entities
- each entities can be a standalone _aggregate root_ too
- what is the difference between aggregate root and entity? An entity may have reference to another entit. We call this associations.
- However, aggregate consist of an aggregate root, and one or more entities that are required to perform business logic atomically.
- You would normally have repository that returns both individual entity (usually for read) as well as aggregate (usually for write operations where the aggregate is needed to validate the dto). Repository may compose the aggregate from multiple enitity.
- Note that for read layer, you moght not want to return the composed aggregate, since not all presentation layer requires that and fetching may be expensive. Also, GraphQL and modern Rest usually allows fields to be selected by client lazily, so the associations are mostly loaded from other tools using dataloader etc that is more robust against N+1 issues
- Aggregate root is responsible for the construction of the entities they owned, e.g. an `OrderCart` aggregate is responsible for adding or removing the `OrderCartItem`. However, they could be created without an identity (id), depending on whether you delegate id generation to the database or application level.
- Most literature suggests returning the aggregate, not entity from repository. However, most aggregate are single-use only for specific scenarios, and most of them do not represent the read-only data required. So aggregate could be returned when there is a need for mutation, and not for read-only queries.
- Aggregate can get very large if there is a lot of nested (and nested-nested entities). Hence, the construction of the aggregate should be really specific to the operation required, and the association should only be referenced by identity (foreign key) rather than the full entity.
- Aggregate children may not outlived the parent. If the aggregate root is deleted, then the children aggregate should be deleted too.
- Aggregate children may not be dependent on the parent - e.g. a post may have many comments. However, the operations on comments does not have to be dependent on the parent.

1. https://gedgei.wordpress.com/2016/06/10/does-ddd-promote-large-aggregates/
2. https://softwareengineering.stackexchange.com/questions/399184/how-to-create-large-readonly-entities-in-ddd

- any good example of aggregate root?
- does aggregate root contains business logic
- is service layer an aggregate root

## read only
http://gorodinski.com/blog/2012/04/25/read-models-as-a-tactical-pattern-in-domain-driven-design-ddd/

## Why not to update multiple entities
https://softwareengineering.stackexchange.com/questions/356106/ddd-why-is-it-a-bad-practice-to-update-multiple-aggregate-roots-per-transaction

## Nesting entity
https://stackoverflow.com/questions/50889425/creating-nested-entities-through-an-aggregate-root-ddd


## Should aggregate returns error or boolean?

That depends
- boolean only represents yes or no. If the aggregate method could potentially return multiple errors to indicate which steps fails, error is better
- boolean does not give information on the error, hence the error has to be called by the client. If the errors are the same, then it will be repeated multiple times. It depends on whether the error really reside in the usecase layer or domain layer.


https://www.alibabacloud.com/blog/an-in-depth-understanding-of-aggregation-in-domain-driven-design_598034
