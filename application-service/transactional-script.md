# Transactional script


Most of the time, transaction script keeps logic very simple:
```go
package main

import (
	"context"
)

type Librarian struct {
	id int64
}

type Member struct {
	id int64
}

type LibraryUseCase struct {
	repo interface {
		BlockMember(ctx context.Context, librarianID, memberID int64) error
	}
}

func (uc *LibraryUseCase) BlockMember(ctx context.Context, librarianID, memberID int64) error {
	return uc.repo.BlockMember(ctx, librarianID, memberID)
}
```


Had we redesign it as DDD:

```go
package main

import (
	"context"
)

type Librarian struct {
	id int64
}

func (l *Librarian) Block(m *Member) {
	m.isBlocked = true
}

type Member struct {
	id        int64
	isBlocked bool
}

type LibraryUseCase struct {
	repo interface {
		FindMemberByID(ctx context.Context, memberID int64) (*Member, error)
		FindLibrarianByID(ctx context.Context, librarian int64) (*Librarian, error)
		UpdateMember(ctx context.Context, m *Member) error
	}
}

func (uc *LibraryUseCase) BlockMember(ctx context.Context, librarianID, memberID int64) error {
	librarian, err := uc.repo.FindLibrarianByID(ctx, librarianID)
	if err != nil {
		return err
	}
	member, err := uc.repo.FindMemberByID(ctx, memberID)
	if err != nil {
		return err
	}
	librarian.Block(member)

	return uc.repo.UpdateMember(ctx, member)
}
```

The number of code has increased dramatically, and it is questionable whether it is more readable or maintainable than the approach using transaction script.


For the implementation above, the `Librarian` is not required to be loaded to ensure the invariant when blocking `Member`. We might have violated the layer responsibilities by placing domain logic (who can block member) in the repository layer. But in some scenarios, it is worth the trade-off (one single query vs multiple queries to fetch the domain entities).
