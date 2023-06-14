# Testing Clean Architecture

1. What kind of tests do we need for each layer?
2. How to handle infrastructure tests? (Probably the answer is don't, testing usecase behaviour is what we want, infrastructure behaviour is not the main concern)


## Strategy

- test end-to-end for all layers
- test layer by layer


The former does not really reduce the amount of tests written. Depending on the usecase, you may end up with combinatorial explosion of parameters to test. 

Testing layer by layer requires you to mock the layer underneath, which increases the LOC for tests. Although it is not necessary to mock each layer underneath.

