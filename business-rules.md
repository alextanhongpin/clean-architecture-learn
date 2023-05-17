# Business rules

## What constitutes of a good business rule?

E.g for login usecase
- password must be valid
- user must notbe banned

both rules above can just be a function Password


- https://wiki.c2.com/?ExplicitBusinessRules
-https://www2.deloitte.com/content/dam/Deloitte/dk/Documents/Grabngo/Business%20rules%20management_030221.pdf

- what is the difference between business rules and usecase
- where should business rules be stored? (application layer requires redeployment, database requires changes too)
- putting business rules in database
  - as data: json schema, scripting language like lua (only data change, but needs to be executed at application layer)
  - as part of db: unique constraints, stored procedures, views, materialized view, triggers (requires schema migration)


## Business rule vs domain rules


Business rules and domain rules are not the same. The only similarity is that they don't have dependencies to external packages. Business rules are applied to domain entities.


Also, if you were to switch your app to another language and have to bring over two main things, it is the domain and business rules.



How to "document" and apply business rules in your application?

https://ericnormand.me/podcast/what-is-the-difference-between-a-domain-model-and-business-rules
