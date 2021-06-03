package nominatim

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var lettersregex = regexp.MustCompile("([A-Za-z ]+)")

type Nominatim struct {
	BaseURL           string
	FormatHouseNumber bool
}

type SearchParameters struct {
	// Search query
	Query string
	// Detailed search, don't combine with Query
	Street     string
	City       string
	County     string
	State      string
	Country    string
	PostalCode string
	// Limit results
	CountryCodes []string
	Limit        int
	Viewbox      []float64 // x1,y1,x2,y2
	// Additional features
	IncludeAddress bool
	IncludeGeoJSON bool
}

type SearchResult struct {
	PlaceID        int                    `json:"place_id"`
	License        string                 `json:"licence"`
	OSMType        string                 `json:"osm_type"`
	OSMID          int                    `json:"osm_id"`
	BoundingBoxStr []string               `json:"boundingbox"`
	LatStr         string                 `json:"lat"`
	LngStr         string                 `json:"lon"`
	DisplayName    string                 `json:"display_name"`
	Class          string                 `json:"class"`
	Type           string                 `json:"type"`
	Importance     float64                `json:"importance"`
	Icon           string                 `json:"icon"`
	GeoJSON        map[string]interface{} `json:"geojson"`
	Address        SearchAddress          `json:"address"`
	BoundingBox    []float64
	Lat            float64
	Lng            float64
}

type SearchAddress struct {
	HouseNumber   string `json:"house_number"`
	Road          string `json:"road"`
	Building      string `json:"building"`
	City          string `json:"city"`
	Suburb        string `json:"suburb"`
	Neighbourhood string `json:"neighbourhood"`
	County        string `json:"county"`
	State         string `json:"state"`
	PostalCode    string `json:"postcode"`
	Country       string `json:"country"`
	CountryCode   string `json:"country_code"`
}

func (n *Nominatim) Search(p SearchParameters) ([]SearchResult, error) {
	// Defaults
	if n.BaseURL == "" {
		n.BaseURL = "https://nominatim.openstreetmap.org"
	}
	// Get initial url
	nurl, err := url.Parse(n.BaseURL)
	if err != nil {
		return []SearchResult{}, err
	}
	// Set path
	nurl.RawPath = "/search"
	// Build query
	q := nurl.Query()
	// Basics
	q.Set("format", "json")
	// Query
	if p.Query != "" {
		q.Set("q", p.Query)
	}
	if p.Street != "" {
		q.Set("street", p.Street)
	}
	if p.City != "" {
		q.Set("city", p.City)
	}
	if p.County != "" {
		q.Set("county", p.County)
	}
	if p.State != "" {
		q.Set("state", p.State)
	}
	if p.Country != "" {
		q.Set("country", p.Country)
	}
	if p.PostalCode != "" {
		q.Set("postalcode", p.PostalCode)
	}
	// Filtering
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if len(p.CountryCodes) != 0 {
		q.Set("countrycodes", strings.Join(p.CountryCodes, ","))
	}
	if len(p.Viewbox) != 0 {
		vb := []string{}
		for _, v := range p.Viewbox {
			vb = append(vb, strconv.FormatFloat(v, 'f', 8, 64))
		}
		q.Set("viewbox", strings.Join(vb, ","))
	}
	// Features
	if p.IncludeAddress {
		q.Set("addressdetails", "1")
	}
	if p.IncludeGeoJSON {
		q.Set("polygon_geojson", "1")
	}
	// Set query
	nurl.RawQuery = q.Encode()
	// Make request
	req, err := http.NewRequest("GET", nurl.String(), nil)
	if err != nil {
		return []SearchResult{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []SearchResult{}, err
	}
	// Decode results
	var results []SearchResult
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return []SearchResult{}, err
	}
	// Convert types
	for i, r := range results {
		lat, err := strconv.ParseFloat(r.LatStr, 64)
		if err != nil {
			return results, nil
		}
		lng, err := strconv.ParseFloat(r.LngStr, 64)
		if err != nil {
			return results, nil
		}
		results[i].Lat = lat
		results[i].Lng = lng
		if len(r.BoundingBoxStr) != 0 {
			bounding := []float64{}
			for _, vstr := range r.BoundingBoxStr {
				vfloat, err := strconv.ParseFloat(vstr, 64)
				if err != nil {
					return results, nil
				}
				bounding = append(bounding, vfloat)
			}
			results[i].BoundingBox = bounding
		}
	}
	if n.FormatHouseNumber {
		for i, result := range results {
			results[i].Address.HouseNumber = lettersregex.ReplaceAllString(result.Address.HouseNumber, "")
		}
	}
	// Return
	return results, nil
}
