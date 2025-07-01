# Adapter Layer

Based on ports and adapter, this layer is responsible for converting types that can be used by the app.

They also contains the implementation of the ports.


## Layer

We can create a special `layer` package that is responsible for adapting types from one layer to another.

This package only does type conversion between layers and can be easily tested by verifying the output type fields are never empty.
