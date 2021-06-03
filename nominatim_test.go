package nominatim

import (
	"errors"
	"strconv"
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

func TestFormatHouseNumber(t *testing.T) {
	n := Nominatim{
		FormatHouseNumber: true,
	}
	results, err := n.Search(SearchParameters{
		Query:          "florida, longwood, 32779",
		IncludeAddress: true,
		IncludeGeoJSON: true,
	})
	if err != nil {
		t.Error(err)
	}
	if len(results) == 0 {
		t.Error(errors.New("No results"))
	}
	if _, ok := strconv.Atoi(results[0].Address.HouseNumber); ok != nil {
		t.Error(errors.New("House number is not a number"))
	}
}
