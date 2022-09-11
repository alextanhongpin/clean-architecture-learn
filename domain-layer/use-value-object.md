# Use value object

Value object does not have identity.

Value object is used as an alternative to primitive obsession.

Another good candidate is value object in service layer. Often, we want to express a business logic without depending on model identity.

For example, we may have a discount service that returns the possible discount that returns discount value object


When client select the discount from the API, they can then be converted back to the value object and be compared for optimistic concurrency.
