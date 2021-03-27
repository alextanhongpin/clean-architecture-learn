# Props Drilling

When there are too many layers, adding a new field will involve cascading updates.

Usually each layer should have their own contracts defined for the input, and hence there will also be a lot of conversion when passing props from layer to layer. While it is possible for each layers to use the same props (usually coined DTO or data transfer objects), it is better to separate them so that each layer is stable. The properties per layer may change too (e.g. converting from application service primitive types to domain primitives).

Transformers for each type should be placed close to the DTO. They can transform to and from the entity from a given layer to the next layer.
