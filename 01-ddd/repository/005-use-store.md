# Stores (or tables)

Repository is not the direct representation of the datasource such as mapping to database.


Instead, repository is just responsible for mapping the results from different data sources (API, Database) to the domain model.

`stores` are essentially the mapping from the database to the application.


## Testing


When testing, we do not need to test the repository layer, as the repository should return a _domain model_, which would have private fields etc. Instead, the `stores` returns it's own data type that represents the tables in the database. Those should be tested.
