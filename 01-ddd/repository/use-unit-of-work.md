# Unit of Work

Unit of Work is a common pattern to abstract the transaction logic. Normally, the _application service_ is responsible for initializing the unit of work.

The repository layer should not have transaction logic, as it is hard to compose multiple repository method when the transaction can only be initialized within.
