# Diffing changes

If we are using an ORM, instead of manually writing the update queries to handle partial fields, we can just diff the changes, and allow updating of certain fields. 
Alternatively is to do this handling on the database side through triggers, or during insert.

```go
package main

import (
	"log"

	"github.com/r3labs/diff/v2"
)

type Order struct {
	ID    string `diff:"id"`
	Items []int  `diff:"items"`
	Name  string `diff:"name"`
}

func main() {
	a := Order{
		ID:    "1234",
		Items: []int{1, 2, 3, 4},
		Name:  "John",
	}

	b := Order{
		ID:    "1234",
		Items: []int{1, 2, 4},
		Name:  "john",
	}

	changelog, err := diff.Diff(a, b)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", changelog)
}
```
