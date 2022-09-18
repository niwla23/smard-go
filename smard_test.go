package smard

import (
	"math"
	"testing"
	"time"
)

var startTime = time.Date(2022, 9, 18, 0, 0, 0, 0, time.Local)
var endTime = time.Date(2022, 9, 18, 8, 0, 0, 0, time.Local)
var factorMWToW = int(math.Pow10(6))

func TestGetProductionData(t *testing.T) {
	got, err := GetProductionData(startTime, endTime)

	if err != nil {
		t.Errorf("error: %v; expected no error", err)
	}

	if len(got) != 32 {
		t.Errorf("len(got) = %v; want 44", len(got))
	}

	r := got[len(got)-1]

	wantedBiomass := 1157000000
	wantedHydropower := 405000000
	wantedWindOffshore := 1005000000
	wantedWindOnshore := 5512000000
	wantedPhotovoltaic := 377000000
	wantedOtherRenewables := 29000000
	wantedNuclear := 1015000000
	wantedLignite := 1365000000
	wantedHardCoal := 968000000
	wantedNaturalGas := 498000000
	wantedPumpedStorage := 162000000
	wantedOtherConventional := 245000000

	if r.Biomass != wantedBiomass {
		t.Errorf("r.Biomass = %v; want %v", r.Biomass, wantedBiomass)
	}

	if r.Hydropower != wantedHydropower {
		t.Errorf("r.Hydropower = %v; want %v", r.Hydropower, wantedHydropower)
	}

	if r.WindOffshore != wantedWindOffshore {
		t.Errorf("r.WindOffshore = %v; want %v", r.WindOffshore, wantedWindOffshore)
	}

	if r.WindOnshore != wantedWindOnshore {
		t.Errorf("r.WindOnshore = %v; want %v", r.WindOnshore, wantedWindOnshore)
	}

	if r.Photovoltaic != wantedPhotovoltaic {
		t.Errorf("r.Photovoltaic = %v; want %v", r.Photovoltaic, wantedPhotovoltaic)
	}

	if r.OtherRenewables != wantedOtherRenewables {
		t.Errorf("r.OtherRenewables = %v; want %v", r.OtherRenewables, wantedOtherRenewables)
	}

	if r.Nuclear != wantedNuclear {
		t.Errorf("r.Nuclear = %v; want %v", r.Nuclear, wantedNuclear)
	}

	if r.Lignite != wantedLignite {
		t.Errorf("r.Lignite = %v; want %v", r.Lignite, wantedLignite)
	}

	if r.HardCoal != wantedHardCoal {
		t.Errorf("r.HardCoal = %v; want %v", r.HardCoal, wantedHardCoal)
	}

	if r.NaturalGas != wantedNaturalGas {
		t.Errorf("r.NaturalGas = %v; want %v", r.NaturalGas, wantedNaturalGas)
	}

	if r.PumpedStorage != wantedPumpedStorage {
		t.Errorf("r.PumpedStorage = %v; want %v", r.PumpedStorage, wantedPumpedStorage)
	}

	if r.OtherConventional != wantedOtherConventional {
		t.Errorf("r.OtherConventional = %v; want %v", r.OtherConventional, wantedOtherConventional)
	}
}

func TestGetProductionForecastData(t *testing.T) {
	got, err := GetProductionForecastData(startTime, endTime)

	if err != nil {
		t.Errorf("error: %v; expected no error", err)
	}

	if len(got) != 32 {
		t.Errorf("len(got) = %v; want 44", len(got))
	}

	r := got[len(got)-1]

	wantedTotal := -1
	wantedPhotovoltaicAndWind := 7188 * factorMWToW
	wantedWindOffshore := 1326 * factorMWToW
	wantedWindOnshore := 5527 * factorMWToW
	wantedPhotovoltaic := 335 * factorMWToW
	wantedOther := -1

	if r.Total != wantedTotal {
		t.Errorf("r.Biomass = %v; want %v", r.Total, wantedTotal)
	}

	if r.PhotovoltaicAndWind != wantedPhotovoltaicAndWind {
		t.Errorf("r.PhotovoltaicAndWind = %v; want %v", r.PhotovoltaicAndWind, wantedPhotovoltaicAndWind)
	}

	if r.WindOffshore != wantedWindOffshore {
		t.Errorf("r.WindOffshore = %v; want %v", r.WindOffshore, wantedWindOffshore)
	}

	if r.WindOnshore != wantedWindOnshore {
		t.Errorf("r.WindOnshore = %v; want %v", r.WindOnshore, wantedWindOnshore)
	}

	if r.Photovoltaic != wantedPhotovoltaic {
		t.Errorf("r.Photovoltaic = %v; want %v", r.Photovoltaic, wantedPhotovoltaic)
	}

	if r.Other != wantedOther {
		t.Errorf("r.Other = %v; want %v", r.Other, wantedOther)
	}

}

func TestGetConsumptionData(t *testing.T) {
	got, err := GetConsumptionData(startTime, endTime)

	if err != nil {
		t.Errorf("error: %v; expected no error", err)
	}

	if len(got) != 32 {
		t.Errorf("len(got) = %v; want 44", len(got))
	}

	r := got[len(got)-1]

	wantedGridLoad := 11087000000
	wantedResidualLoad := 4192000000
	wantedPumpedStorage := 518000000

	if r.GridLoad != wantedGridLoad {
		t.Errorf("r.GridLoad = %v; want %v", r.GridLoad, wantedGridLoad)
	}

	if r.ResidualLoad != wantedResidualLoad {
		t.Errorf("r.ResidualLoad = %v; want %v", r.ResidualLoad, wantedResidualLoad)
	}

	if r.PumpedStorage != wantedPumpedStorage {
		t.Errorf("r.PumpedStorage = %v; want %v", r.PumpedStorage, wantedPumpedStorage)
	}
}
