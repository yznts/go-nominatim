
# GO Nominatim

Wrapper around Nominatim API

## Import

```go
import(
    nm "github.com/yuriizinets/go-nominatim"
)
```

## Usage

```go
n := nm.Nominatim{}
results, err := n.Search(nm.SearchParameters{  // Check SearchResult struct for details
    Query:          "Miami, Florida",
    CountryCodes:   []string{"us"},
    IncludeAddress: true,
    IncludeGeoJSON: true,
})
if err != nil {
    log.Panicln(err)
}
log.Println(results[0].DisplayName)
log.Println(results[0].Lat, results[0].Lng)
log.Println(results[0].BoundingBox)
```
