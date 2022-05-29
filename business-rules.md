# Business rules

- https://wiki.c2.com/?ExplicitBusinessRules
- https://www2.deloitte.com/content/dam/Deloitte/dk/Documents/Grabngo/Business%20rules%20management_030221.pdf

- what is the difference between business rules and usecase
- where should business rules be stored? (application layer requires redeployment, database requires changes too)
- putting business rules in database
  - as data: json schema, scripting language like lua (only data change, but needs to be executed at application layer)
  - as part of db: unique constraints, stored procedures, views, materialized view, triggers (requires schema migration)

