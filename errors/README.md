# Domain Errors


There are multiple ways to define errors, and sometimes we want them to be decorated with more capabilities:

- errors with suberrors
- errors with context or rules
- translations
- error wrapping
- multiple errors


## Multiple errors example

This is especially useful when you want to return multiple errors rather than single errors. The implementation below is similar to Formik (see [2]). There is also the _notification_ pattern, see [1]. 
```js
function validateCreateUser(user) {
  const errors = {}
  
  if (!user.email) {
    errors['email'] = 'required'
  } else if (validateEmail(user.email)) {
    errors['email'] = 'invalid'
  }
  
  if (!user.password) {
    errors['password'] = 'required'
  } else if (user.password.length < 8) {
    errors['password'] = 'too short, min 8 characters'
  }
  
  return errors
}
```

Depending on languages, there may already be a way to handle multiple errors, such as JavaScript's __AggregateError__, see [3]:
```js
try {
  throw new AggregateError([
    new Error("some error"),
  ], 'Hello');
} catch (e) {
  console.log(e instanceof AggregateError); // true
  console.log(e.message);                   // "Hello"
  console.log(e.name);                      // "AggregateError"
  console.log(e.errors);                    // [ Error: "some error" ]
}
```

Example with golang:

```go
package main

import (
	"errors"
	"fmt"
	"strings"
)

type AggregateError struct {
	errors map[string]string
}

func (a AggregateError) Error() string {
	var messages []string
	for field, err := range a.errors {
		messages = append(messages, fmt.Sprintf("%s: %s", field, err))
	}
	return strings.Join(messages, "\n")
}

func (a AggregateError) Is(err error) bool {
	_, ok := err.(*AggregateError)
	return ok
}

func (a *AggregateError) Add(field, msg string) {
	a.errors[field] = msg
}

func NewAggregateError() *AggregateError {
	return &AggregateError{
		errors: make(map[string]string),
	}
}

type User struct {
	Email             string
	EncryptedPassword string
}

func encrypt(plaintext string) string {
	return plaintext
}

func NewUser(email, password string) (User, error) {
	err := NewAggregateError()

	if email == "" {
		err.Add("email", "required")
	} else if !strings.Contains(email, "@") {
		err.Add("email", "invalid")
	}

	if password == "" {
		err.Add("password", "required")
	} else if len(password) < 8 {
		err.Add("password", "too short")
	}

	return User{
		Email:             email,
		EncryptedPassword: encrypt(password),
	}, err
}

func main() {
	_, err := NewUser("john.doe", "5")
	fmt.Println(err, errors.Is(&AggregateError{}, err))
}
```

# References
1. [Martin Fowler: Replace throw with Notification](https://martinfowler.com/articles/replaceThrowWithNotification.html)
2. [Formik validation](https://formik.org/docs/guides/validation)
3. [JavaScript AggregateError](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/AggregateError)
