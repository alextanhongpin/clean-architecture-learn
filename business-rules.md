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

### Difference between business rules and domain rules

Domain includes entities and value objects. The difference between _domain_ and _business rules_ is that _domain_ is usually general, while _business rules_ is specific. Some examples below:

| Domain                            | Business Rule                                               |
| --                                | --                                                          |
| Phone Number: format must be E164 | Supported country code is only for Malaysia and Singapore   |
| OTP: must be at least 4 digits    | Use 6 digits for Payout OTP, use 4 digits for Auth OTP      |

### Errors for domain and business rules

Domain errors are usually more general. Business rules errors are more specific.


### Creation step

Creation is actually part of domain. So the params to create a new entity should belong in the domain.


How to "document" and apply business rules in your application?

https://ericnormand.me/podcast/what-is-the-difference-between-a-domain-model-and-business-rules


## Business flow

Aside from domain and business rules, we also have business flow.

If the business flow involves interactions with the end-users, then we usually call it usecase layer.

However, if you are developing on the backend side, most of the flows that you will design involves mostly interaction between different systems (or internally).

The usecase layer usually treats the system flow as a black box - they only care about the input and the output. However, what usually happens internally can be quite complex.
