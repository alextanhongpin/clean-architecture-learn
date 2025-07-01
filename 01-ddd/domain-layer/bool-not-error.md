The domain model should return a bool, not a specific error since the error may be uo to the client to interpret.

For example, when checking if a user is authorized, we can just return bool. Let the user figure out exactly what errors to return which could be business logic dependent.
