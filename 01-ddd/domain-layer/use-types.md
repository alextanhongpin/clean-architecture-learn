# Use Types

Most applications needs to depends on custom methods for primitives, or even new primitives.

For example, it is common to implement Set in golang due to lacking of such data structure in the standard library.


The general rule is, sincethey are not value objects or domain models, and are usually global to the whole application, they should be placed in the "internal/types" folder.


This applies to the following
- stringcase for strings
- custom math
- null types
- custom context
- custom date types
