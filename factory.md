# Factory

- what is a factory? A Factory is an object that has the single responsibility of creating other objects.
- _factory_ marks the __beginning__ of an entity, _repository_ the __middle__ and __end__
- use _factory_ to build the entity, and _repository_ to persist
- one deals with creation, another deals with persistence.
- domain factory can also be user the other way around, for reconstitution, meaning to build an entity back from external sources, whether it’s from JSON (deserialising) or from another domain service or even repository
- can factory access repository (? yes, but avoid that if possible. Are there any good usecase to demonstrate it?)
- factory is responsible for constructing complex aggregates in their valid states. So a constructor can be a factory too
- factory can be plain functions, for complex factory that requires other 


# Builder

How about builder pattern for creating test data for database, especially with nested aggregates?

```go
package main

import (
	"fmt"
	"log"
)

func main() {
	builder := NewQuestionBuilder().
		SetUserParams(userParams{Name: "john"}).
		SetAnswerParams(answerParams{}).
		SetQuestionParams(questionParams{})
	if err := builder.Build(); err != nil {
		log.Println(err)
	}
	fmt.Println(builder.User())
}

type User struct {
	ID   string
	Name string
}
type Answer struct {
	ID     string
	UserID string
}

type Question struct {
	AnswerID string
	UserID   string
}

type userParams struct {
	Name string
}
type answerParams struct{}
type questionParams struct{}
type QuestionBuilder struct {
	*userParams
	*answerParams
	*questionParams

	user     *User
	answer   *Answer
	question *Question
}

func NewQuestionBuilder() *QuestionBuilder {
	return &QuestionBuilder{}
}

func (b *QuestionBuilder) SetUserParams(u userParams) *QuestionBuilder {
	b.userParams = &u
	return b
}

func (b *QuestionBuilder) SetAnswerParams(a answerParams) *QuestionBuilder {
	b.answerParams = &a
	return b
}

func (b *QuestionBuilder) SetQuestionParams(q questionParams) *QuestionBuilder {
	b.questionParams = &q
	return b
}

func (b *QuestionBuilder) User() User {
	return *b.user
}

func (b *QuestionBuilder) Answer() Answer {
	return *b.answer
}

func (b *QuestionBuilder) Question() Question {
	return *b.question
}

func (b *QuestionBuilder) Build() error {
	if b.userParams == nil {
		b.userParams = &userParams{}
	}
	// CreateUser through repository.
	b.user = &User{Name: b.userParams.Name}

	if b.answerParams == nil {
		b.answerParams = &answerParams{}
	}
	b.answer = &Answer{}

	if b.questionParams == nil {
		b.questionParams = &questionParams{}
	}

	b.question = &Question{}
	return nil
}
```

# References

1. [Domain Driven Design: Services and Factories](https://archiv.pehapkari.cz/blog/2018/03/28/domain-driven-design-services-factories/)
