# Seeding nested data

Not part of clean architecture, but we want a "clean" option to seed nested data, as well as allow overriding of entity at each nested level.

```go
package main

import (
	"context"
)

type User struct {
	ID int64
}
type Product struct {
	ID     int64
	UserID int64
}

type Item struct {
	ID        int64
	ProductID int64
}

// Builder creates a default builder with all data faked but overridable.
func buildProduct() *ProductBuilder {
	return &ProductBuilder{}
}

// Creator creates the entry in the db by calling builder.Build()
func createProduct(ctx context.Context, pb *ProductBuilder) *Product {
	// Creates a Product in DB
	return &Product{}
}

func buildUser() *UserBuilder {
	return &UserBuilder{}
}

func createUser(ctx context.Context, b *UserBuilder) *User {
	// Creates a User in DB
	return &User{}
}

func buildItem() *ItemBuilder {
	return &ItemBuilder{}
}

func createItem(ctx context.Context, b *ItemBuilder) *Item {
	// Creates a Item in DB
	return &Item{}
}

// For simple use cases.
func seedUser(ctx context.Context) *User {
	b := buildUser()
	b.WithName("john")
	return createUser(ctx, b)
}

func seedProduct(ctx context.Context) *Product {
	user := seedUser(ctx)
	b := buildProduct()
	b.WithUserID(user.ID())
	return createProduct(ctx, b)
}

func seedItem(ctx context.Context) *Item {
	product := seedProduct(ctx)
	b := buildItem()
	b.WithProductID(product.ID)
	return createItem(b)
}

func seedItems(ctx context.Context, n int) []Item {
	product := seedProduct(ctx)
	items := make([]Item, n)
	for i := 0; i < n; i++ {
		b := buildItem()
		b.WithProductID(product.ID)
		item := createItem(b)
		items[i] = item
	}
	return items
}

// For greater granularity, allow overwriting of fields for each associations.
type ItemSeeder struct {
	UserBuilder *UserBuilder
	User        *User // We want to access User, Product and Items after creation.

	ProductBuilder *ProductBuilder
	Product        *Product

	ItemBuilders []*ItemBuilder
	Items        []Item
}

func NewItemSeeder() *ItemSeeder {
	return &ItemSeeder{
		UserBuilder:    buildUser(), // Defaults.
		ProductBuilder: buildProduct(),
		ItemBuilders:   buildItems(),
	}
}

func (s *ItemSeeder) Seed(ctx context.Context) {
	ub := s.UserBuilder
	s.User = createUser(ctx, ub)

	pb := s.ProductBuilder
	pb.WithUserID(s.User.ID())
	s.Product = createProduct(ctx, pb)

	s.Items = make([]Item, n)
	for i := 0; i < n; i++ {
		b := s.ItemBuilders[i]
		b.WithProductID(s.Product.ID)
		item := createItem(b)
		s.Items[i] = item
	}
}

func main() {
	is := NewItemSeeder()
	is.UserBuilder.
		WithName("john").
		WithAge(23) // Overrides user's fields.
	is.ProductBuilder().
		WithName() // Overrides product's fields.
	is.Seed(context.Background()) // Creates User, Product and Items
}
```
