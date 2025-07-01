Each layer has its own set of responsibility, following the single responsibility principle.

Imagine calling parse uuid method in the repo, it is a violation of layers. The repo should expect a uuid instead if primitive string.
