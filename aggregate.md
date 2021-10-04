# Aggregate Root


- Aggregate root is a collection of entities
- each entities can be a standalone _root_ too
- what is the difference between aggregate root and entity? An entity may have reference to another entit. We call this associations.
- However, aggregate consist of an aggregate root, and one or more entities that are required to perform business logic atomically. 
- You would normally have repository that returns both individual entity (usually for read) as well as aggregate (usually for write operations where the aggregate is needed to validate the dto). Repository may compose the aggregate from multiple enitity. 
- Note that for read layer, you moght not want to return the composed aggregate, since not all presentation layer requires that and fetching may be expensive. Also, GraphQL and modern Rest usually allows fields to be selected by client lazily, so the associations are mostly loaded from other tools using dataloader etc that is more robust against N+1 issues

1. https://gedgei.wordpress.com/2016/06/10/does-ddd-promote-large-aggregates/
2. https://softwareengineering.stackexchange.com/questions/399184/how-to-create-large-readonly-entities-in-ddd


## read only
http://gorodinski.com/blog/2012/04/25/read-models-as-a-tactical-pattern-in-domain-driven-design-ddd/

## Why not to update multiple entities
https://softwareengineering.stackexchange.com/questions/356106/ddd-why-is-it-a-bad-practice-to-update-multiple-aggregate-roots-per-transaction

## Nesting entity
https://stackoverflow.com/questions/50889425/creating-nested-entities-through-an-aggregate-root-ddd
