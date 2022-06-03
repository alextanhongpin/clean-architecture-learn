# How to design aggregate root/child with polymorphism?

- does OOP principles still matters in DDD?

https://softwareengineering.stackexchange.com/questions/337616/in-ddd-how-do-i-persist-an-aggregate-containing-polymorphism
https://blog.ploeh.dk/2011/05/31/AttheBoundaries,ApplicationsareNotObject-Oriented/
https://stackoverflow.com/questions/30400317/domain-driven-design-and-working-with-polymorphic-child-entities


## API Pattern for polymorphism

```ts
interface Vehicle {
  type: 'car' | 'motor'
  car?: Car
  motor?: Motor
  wheels: number
}

interface Car {
  make: string
  year: number
}

interface Motor {
  brand: string
}
```

```ts

type UserError = {
  __typename: 'error'
  message: string
}

type UserNotFound = UserError & {
  userId: string
}

type User = {
  __typename: 'user'
  name: string
  age: number
}

type UserResult = User | UserNotFound

let user: UserResult = {
  __typename: 'user',
  name: 'john',
  age: 20
}

user = {
  __typename: 'error',
  userId: '1',
  message: 'bad',
}

function printUser(user: UserResult) {
  switch (user.__typename) {
    case 'user':
      console.log(user.name)
      break
    case 'error':
      console.log(user.userId, user.message)
      break
    default:
      throw new Error('invalid user')
  }
}

printUser(user)
```
