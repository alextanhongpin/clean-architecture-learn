# Presentation Layer

- deserialization of request
- serialization of response
- localization of content/error messages. While this can be done on the Frontend, it is more flexible if done on the Backend
- authorization through middleware
- can be CLI/REST/GraphQL etc
- calls the application service layer (usecase layers)
- strictly no business logic (even if-else for rendering content)
