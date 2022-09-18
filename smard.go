package smard

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RequestForm struct {
	Format        string `json:"format"`
	Language      string `json:"language"`
	ModuleIds     []int  `json:"moduleIds"`
	Region        string `json:"region"`
	TimestampFrom int    `json:"timestamp_from"`
	TimestampTo   int    `json:"timestamp_to"`
	Type          string `json:"type"`
}

type RequestBody struct {
	RequestForm []RequestForm `json:"request_form"`
}

type ProductionDataRow struct {
	Timestamp         time.Time
	Biomass           int
	Hydropower        int
	WindOffshore      int
	WindOnshore       int
	Photovoltaic      int
	OtherRenewables   int
	Nuclear           int
	Lignite           int
	HardCoal          int
	NaturalGas        int
	PumpedStorage     int
	OtherConventional int
}

type ProductionForecastDataRow struct {
	Timestamp           time.Time
	Total               int
	PhotovoltaicAndWind int
	WindOffshore        int
	WindOnshore         int
	Photovoltaic        int
	Other               int
}

type ConsumptionDataRow struct {
	Timestamp     time.Time
	GridLoad      int
	ResidualLoad  int
	PumpedStorage int
}

func parseTimestampFromColumns(columns []string) (time.Time, error) {
	timestamp, err := time.Parse("02.01.2006", columns[0])
	if err != nil {
		return time.Time{}, err
	}

	ot, err := time.Parse("15:04", columns[1])
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), ot.Hour(), ot.Minute(), 0, 0, time.Local), nil
}

// takes a Megwatt string in and converts it to a watt integer
func mwhStrToWhInt(mwh string) int {
	parsed, err := strconv.ParseInt(strings.TrimSpace(strings.Replace(mwh, ".", "", -1)), 10, 32)
	if err != nil {
		return -1
	}
	return int(parsed * 1_000_000)
}

// gets data from smard.de in the given range for given modules. Returns the body as plain csv with semicolons as seperators.
// don't confuse the dot with a decimal comma, it divides every three digits. I have no idea why you would put that in a CSV but that is the way it is.
func getRawData(from time.Time, to time.Time, moduleIds []int) (string, error) {
	filters, err := json.Marshal(RequestBody{
		RequestForm: []RequestForm{{Format: "CSV", Language: "de", ModuleIds: moduleIds, Region: "de", TimestampFrom: int(from.UnixMilli()), TimestampTo: int(to.UnixMilli()), Type: "discrete"}},
	})
	if err != nil {
		return "", err
	}
	resp, err := http.Post("https://www.smard.de/nip-download-manager/nip/download/market-data", "application/json", bytes.NewBuffer(filters))

	if err != nil {
		return "", err
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	return string(body), nil
}

// gets all production data in the given time range from smard.de and returns them as a list of ProductionDataRows
func GetProductionData(from time.Time, to time.Time) ([]ProductionDataRow, error) {
	productionModuleIds := []int{1001224, 1004066, 1004067, 1004068, 1001223, 1004069, 1004071, 1004070, 1001226, 1001228, 1001227, 1001225}

	data, err := getRawData(from, to, productionModuleIds)
	if err != nil {
		return []ProductionDataRow{}, err
	}

	var parsedData []ProductionDataRow

	for rowIndex, rawRow := range strings.Split(data, "\n") {
		// skip first row, it contains only headings
		if rowIndex == 0 {
			continue
		}

		parsedRow := ProductionDataRow{}
		columns := strings.Split(rawRow, ";")

		ts, err := parseTimestampFromColumns(columns)
		if err != nil {
			// we have no date, continue
			continue
		}
		parsedRow.Timestamp = ts

		parsedRow.Biomass = mwhStrToWhInt(columns[2])
		parsedRow.Hydropower = mwhStrToWhInt(columns[3])
		parsedRow.WindOffshore = mwhStrToWhInt(columns[4])
		parsedRow.WindOnshore = mwhStrToWhInt(columns[5])
		parsedRow.Photovoltaic = mwhStrToWhInt(columns[6])
		parsedRow.OtherRenewables = mwhStrToWhInt(columns[7])
		parsedRow.Nuclear = mwhStrToWhInt(columns[8])
		parsedRow.Lignite = mwhStrToWhInt(columns[9])
		parsedRow.HardCoal = mwhStrToWhInt(columns[10])
		parsedRow.NaturalGas = mwhStrToWhInt(columns[11])
		parsedRow.PumpedStorage = mwhStrToWhInt(columns[12])
		parsedRow.OtherConventional = mwhStrToWhInt(columns[13])

		parsedData = append(parsedData, parsedRow)
	}
	return parsedData, nil
}

// gets all production forecast data in the given time range from smard.de and returns them as a list of ProductionForecastDataRows
func GetProductionForecastData(from time.Time, to time.Time) ([]ProductionForecastDataRow, error) {
	data, err := getRawData(from, to, []int{2000122, 2005097, 2000715, 2000125, 2003791, 2000123})
	if err != nil {
		return []ProductionForecastDataRow{}, err
	}

	var parsedData []ProductionForecastDataRow
	for rowIndex, rawRow := range strings.Split(data, "\n") {
		// skip first row, it contains only headings
		if rowIndex == 0 {
			continue
		}

		parsedRow := ProductionForecastDataRow{}
		columns := strings.Split(rawRow, ";")

		ts, err := parseTimestampFromColumns(columns)
		if err != nil {
			// we have no date, continue
			continue
		}
		parsedRow.Timestamp = ts

		parsedRow.Total = mwhStrToWhInt(columns[2])
		parsedRow.PhotovoltaicAndWind = mwhStrToWhInt(columns[3])
		parsedRow.WindOffshore = mwhStrToWhInt(columns[4])
		parsedRow.WindOnshore = mwhStrToWhInt(columns[5])
		parsedRow.Photovoltaic = mwhStrToWhInt(columns[6])
		parsedRow.Other = mwhStrToWhInt(columns[7])

		parsedData = append(parsedData, parsedRow)
	}
	return parsedData, nil
}

// gets all consumption data in the given time range from smard.de and returns them as a list of ConsumptionDataRows
func GetConsumptionData(from time.Time, to time.Time) ([]ConsumptionDataRow, error) {
	consumptionModuleIds := []int{5000410, 5004387, 5004359}

	data, err := getRawData(from, to, consumptionModuleIds)
	if err != nil {
		return []ConsumptionDataRow{}, err
	}

	var parsedData []ConsumptionDataRow

	for rowIndex, rawRow := range strings.Split(data, "\n") {
		// skip first row, it contains only headings
		if rowIndex == 0 {
			continue
		}

		parsedRow := ConsumptionDataRow{}
		columns := strings.Split(rawRow, ";")

		ts, err := parseTimestampFromColumns(columns)
		if err != nil {
			// we have no date, continue
			continue
		}
		parsedRow.Timestamp = ts

		parsedRow.GridLoad = mwhStrToWhInt(columns[2])
		parsedRow.ResidualLoad = mwhStrToWhInt(columns[3])
		parsedRow.PumpedStorage = mwhStrToWhInt(columns[4])

		parsedData = append(parsedData, parsedRow)
	}
	return parsedData, nil
}

// func main() {
// 	data := GetProductionForecastData(time.Date(2022, 9, 10, 0, 0, 0, 0, time.Local), time.Now())
// 	for _, row := range data {
// 		fmt.Printf("time: %v: pump: %v\n", row.Timestamp, row.Other)
// 	}
// }
