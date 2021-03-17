package nominatim

import (
	"errors"
	"testing"
)

func TestSearch(t *testing.T) {
	// Search for Key Biscayne
	n := Nominatim{}
	results, err := n.Search(SearchParameters{
		Query:          "Sunny Isles Beach, Miami, Florida",
		IncludeAddress: true,
		IncludeGeoJSON: true,
	})
	if err != nil {
		t.Error(err)
	}
	if len(results) == 0 {
		t.Error(errors.New("No results"))
	}
}
