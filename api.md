# Filter

How to handle filter options with query string?

1. https://www.npmjs.com/package/json-query
2. https://github.com/persvr/rql
3. https://www.moesif.com/blog/technical/api-design/REST-API-Design-Filtering-Sorting-and-Pagination/


Try not to expose all fields to the public. Ideally every endpoint should serve a specific usecase.

Not all filter combinations are possible, so having a more specific endpoint makes the API more reliable. 

Having a generic endpoint with all possible filter combinations makes the client more reliant on it, making it harder to change
.
