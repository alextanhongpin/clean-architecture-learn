# Packages

How do we design modules? When do we place files in the same module, and when do we separate them?

For most applications, we normally group them by responsibility, e.g.
- controller
  - user.go
  - product.go
- repository
  - user.go
  - product.go

We could have grouped it by features too, e.g.

- user
  - controller.go
  - repository.go
 
Which approach is better? That depends on the question in mind.

- are the modules standalone? then the latter is better, assuming no imports from other modules except the main
