# Use ranking

Sorting and filtering is not a domain concern - the responsibility lies solely in the repository layer. We cannot afford to query the entire table and perform sorting in-memory.

When ranking entities that requires computed data, we can split the query into two part. The first part is responsible for computing the ranking, and returning the ids of the entity,
and the second query actually queries the entities based on the id.

The result will then be the collection of entity with the ranking.

```go
type PostgresArtistRepository struct{}

type Artist struct {
	ID   string
	Name string
}

type ArtistRank struct {
	Rank   int
	Artist Artist
}

func (r *PostgresArtistRepository) ArtistsByRanking(after int64) ([]ArtistRank, error) {
	type artistRank struct {
		artistID string
		rank     int
	}

	// Query the artists rank using another query.
	// You can use CTE too, but for this example, we demonstrate splitting it into two separate queries.
	var artistsRanks []artistRank
	rankByArtistID = make(map[string]int64)

	artistIDs := make([]int64, len(artistRanks))
	for i, ar := range artistRank {
		artistIDs[i] = ar.artistID
		rankByArtistID[ar.artistID] = ar.rank
	}

	// Query artists by id
	_ = artistIDs

	var artists []Artist

	result := make([]ArtistRank, len(artists))
	// Sort the artist by rank again
	for i, artist := range artists {
		result[i] = ArtistRank{
			Artist: artist,
			Rank:   rankByArtistID[artist.ID],
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].rank < result[j].rank
	})

	return result
}
```
