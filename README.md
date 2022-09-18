# smard-go
go module to fetch energy data for germany from smard.de

## example

```go
got, err := GetProductionForecastData(time.Date(2022, 9, 10, 0, 0, 0, 0, time.Local), time.Now())
```
this gets the production data from 9-10-2022 to now. The data is returned as a slice of ProductionDataRow. Look at smard.go for available fields.

## publishing
create a tag with the new version in git

publish package with
```sh
GOPROXY=proxy.golang.org go list -m github.com/niwla23/smard-go@TAG
```