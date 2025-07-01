## Using interface as a contract

- specifying concrete types means modifying input structs when a new field is added
- On the other hand, changing interface just requires changing the internal implementation (granted that all those fields are provided)

Eg. Type User has 10 fields

There are two ways to update

Domain update
1. fetch the existing entity first
2. Make changes
3. Save changes

while this may sound like its perfect, it will eventually be prone to data race. If multiple users are updating the same entity, you may want to lock the row first. This is also the same issue when returning an entity after update (not insert), because their representation may differ. Hence the changeset patterns makes sense.

Updating whole row is not a performance issue on postgres, so updates can be simplified greatly.

You may need to fetch and update if
- past data is required, e.g for state transition (though this could be part of sql filter query)
- there is not enough information to validate if the changes are valid, e.g acl, bool flags or statuses
- there may be other associations that are required for cross validation

Change set
1. pre validate changes locally
2. Save changes
3. Handle errors

Change set makes sense for
- bulk
- Non business logic fields, like remarks
- there is no need to validate against past data
- changeset are value object can be validated as a whole, e.g. user address

Ddd, vs crud
Data vs behaviour
Ddd however depends on data validation too.
Data for ddd should be valid.
Modifying data vs applying behaviour (are value objects part of ddd)
Validation for aggregates vs single entity.
Validation before and after
Batch validation, precise row with error. 
Using interfaces to represent entity.
What is an entity? Something that has lifecycle?

# Repository
- business logic in db
- mapping errors to domain
- passive errors, errors returned after insert or update
- active errors, errors returned before insert
- bulk errors
- errors from triggers, constraint, integrity
- if the data can be validated in bulk before insert, then aggregate them before insertion
- else, insert them line by line in a single transaction to get the specific error message
- behaviour is global, since the db is the main aggregate. data fetched in application is just partial local aggregate. 

How to represent an entity?
- struct with public fields
- Struct with private fields, and getters
- Interfaces
-   Immutability
-   Good for whitelist
- Use valid method to check if it’s instantiated, or use a custom vet to check unwanted initialisation, or add panic on method call on initialised fn. run the vet before app run instead of after test.

Service/usecase/domain
- returning, error or bool/nil?

Updating in repository
- find and update pattern, returning the new submission will save all the data, returning nil will skip
- Wrapping in tx 
- Tx with event context for publishing outbox messages


Scheduling
The art of scheduling
Scheduling ahead of time
When things go south 
Retries
Accuracy in time
History of changes

Time range for scheduling
Unschedule
Race conditions
Schedule once only
Clearing schedules 
Bitemporal


Workflow step and domain model
Embedding model in flow


# No to clean architecture


http://johannesbrodwall.com/2014/07/10/the-madness-of-layered-architecture/
https://garywoodfine.com/why-i-dont-like-layered-architecture-for-microservices/
