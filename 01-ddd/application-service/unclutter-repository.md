# Uncluttering repositories in Usecase Layer


## Problem


When dealing with large usecases, it is not uncommon to have a lot of dependencies. One of the common dependencies is repository.

```go
type UserUsecase struct {
	accountRepo accountRepository
	userRepo userRepository
	membershipRepo membershipRepository
}
```

However, this becomes a problem, because we will have to deal with interface explosion where there are simply too many interfaces to be defined. Some developers took the shortcut by just reusing a common interface (which IMHO is a bad practice):

```go
type UserUsecase struct {
	accountRepo repository.Account
	userRepo repository.User
	membershipRepo repository.Membership
}
```

It is a bad practice, because not all methods of the interface will be used and they become less specific. Also, during testing, it requires more effort to create mocks that implements all the methods (although code generation maeks it possible, it is still not idiomatic).


## Proposed Solution


The solution is to just create one repository per usecase. In the example above, we should only have:


```go
type UserUsecase struct {
	userRepo userRepository
}
```

Note that this changes the meaning of repository substancially. The `userRepository` does not mean the layer that interacts with the `users` table in the database.

Instead, it just means all the operations that the `UserUsercase` requires to interact with the database. Most repositories are designed wrong - as they are just table mappers in disguise.


## Conclusion


Use one repository per usecase. This makes testing much easier too, and limits the scope of access of your usecase layer.
