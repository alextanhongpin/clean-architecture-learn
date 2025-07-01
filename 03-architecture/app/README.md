# App


## User Story

User story describes the goal of application:

```
As a user,
I want to track my expenses,
In order to ensure my expenses are well-controlled.
```

## Aggregates

Aggregate is our root model.

The main aggregate is the `Expense`. It records what the expense is used for, the amount, and additional notes. The date the expense is made can be customized - which defaults to now.

The user can choose a predefined `Category`, or create a list of customizable `Category`.

```typescript
interface Expense {
	remark: string // What the expense is for
	category: Category
	notes: string
	amount: number
	createdAt: Date
}
```

## Queries

Queries is what our data can tell us.

The `Expense` is just a collection - we will a list of expenses, in which we can perform the following queries:

- total expenses by day/week/month/year
- month-by-month expenses changes
- top expenses by category
- list all expenses in descending order
- list expenses by filter (category, day/week/month)
- average expenses by day/week/month/year


## Usecases

Usecases describes the interactions the user can (or cannot do) with the app, and the different scenarios.

They are mostly mutations.

- user manage expenses (manage includes the basic CRUD)
- user search expense (by name)

## Code


```go
type Expense struct {
	ID int64
	Remark string
	Category Category
	Notes string
	Amount int64
	CreatedAt time.Time
}

type ExpenseList []Expense

func (list ExpenseList) Total() int64 {
	panic("not implemented")
}

type Category struct {
	ID int64
	Name string
	CreatedAt time.Time
}
```
