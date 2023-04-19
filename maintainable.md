# Maintainable Code

Every time I join a new time, I face the issue with legacy code
- no readme
- no details on how to even setup the local environment
- missing `.env` of course
- frozen cause I don't know where to start


However, that happens for all projects, and I realised even my code will not be something others can understand.

How to document code well, so that anyone who takes over it can understand and be productive with it immediately?

- document architecture
  - most architecture are complex - there are no fixed standard (conway's law at play, or more like it's the choice of the previous maintainer anyway)
  - the train of thought is hard to understand, like how do I add a new api, how to test it, what integration test do I need
- document domain logic 
  - this is the most crucial part, and can somehow be tangled in mess if the domain layer is not there (code mixed in the infrastructure or usecase layer)
  - due to some shortcuts taken by developer, there are values hardcoded or _adjustments_ made that nobody can remember

Things to do
- [ ] document a basic tutorial on _how to do x_
- [ ] add references to the initial tutorial if any
- [ ] add contacts/maintainers in the README so that people can look for them
- [ ] separate domain layer


