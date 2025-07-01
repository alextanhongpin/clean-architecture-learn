# Use One repository (if possible)


Instead of creating many specific repository, consolidate them under a single repository.

This makes testing easier. Also, you do not need to inherit all the other methods that are not used.

If the `userRepo` has 20 methods, you will probably only need 2-3 methods for your specific usecase.

<table>
<thead>
	<tr>
		<th>Bad</th>
		<th>Good</th>
	</tr>
</thead>
<tbody>
<tr> <!--NOTE: It must be indented to the most-left, otherwise it won't work.-->
<td>

Many individual repositories.

```go
package usecase

                               .
type AuthUsecase struct {
	userRepo    userRepository
	accountRepo accountRepository
	cache       cache
	background  background
}
```


</td>
<td>

A single repo per usecase.

```go
package usecase

                               .
type AuthUsecase struct {
	repo authRepository
}




```

</td>
</tr>
</tbody>
</table>

## Use single repository for all usecases

It will be easier to use a single repository for all usecases.

So the usecase only verifies that the single repository implements all the method that is required.

This reduces boilerplate a lot, as well as the number of created objects.

