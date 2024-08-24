# Repository


Aka persistence layer.
Use repository to
- read/write data to any data source, not limited to database, can be API calls
- returns domain models from those data sources
- calls external API and transforms the response to the domain models
- abstracts cache/message queue implementation
- can return more than one domain models type
- does not imply one repo = one database table, this is wrong
- can have business logic like unique constraints in db, converts those to domain errors

