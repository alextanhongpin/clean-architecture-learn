# Use Layer Mappers


## Problem

When working with layered architecture, you need to map types from one layer to another. Mapping types is error prone and manual. Also, there will also be a question on which layer should handle the mapping.


## Proposed Solution


For each layer, introduce a mapping layer. Assuming we have two layers, parent and child layer, and the parent layer calls the child layer.


When the parent layer needs to call the child layer, the parent layer must first convert the types to child layer.
The response will also be converted by the parent layer. The child layer should not know anything about the parent layer types.


There is also another option, which is to map to another known types by all the layers, e.g. repository mapping the type to domain types, which the usecase layer knows.


## Conclusion


Map in one direction.
