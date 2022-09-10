# Return types

## What should an application service return?

Read-only types. 

## Can an application service return domain model?

Yes, if only they consist of read-only, immutable attributes. If the domain model contains setter or anything that produces side effects, it should be mapped to a usecase types before returning to the presentation layer.
