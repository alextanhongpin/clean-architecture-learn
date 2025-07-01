# Use external calls


There are implementation where third-party API calls are called with the repository being the interface, e.g. third-party payment etc.

Do not do this. Instead, create another interface specifically for the payment, separate from the repository interface.


The exception is when the implementation does return a domain model.
