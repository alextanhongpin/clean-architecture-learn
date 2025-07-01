# Use access delegation

Access control is common in applications. However, many has scorn against putting authorization in the domain layer.

- where do authorization belongs
- how do validate access for particular user (owner of the entity created)
- what is the difference between policies and rules?


# Use layers to separate authorization level

Most backend applications are n-tier applications. Common layers includes controller, usecase and repository layer.

Each layer has separate responsibility. When moving from the outer to infer layer, the responsibility changes from checking permission to executing business rules.

For application that uses JWT for example, the logic to generate the JWT token as well as validating it belongs to the controller layer, not usecase.


Authorization should happen at the outermost layer.

However, there are some domain specific access control that can only be addressed by the domain layer.

For example, multi tenant application as well as ownership of entity.
