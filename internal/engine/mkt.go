package engine

import (
	"math"

	"example.com/coldchain-lake/internal/model"
)

const (
	activationEnergyJ = 83144.0
	gasConstant       = 8.314
)

// MeanKineticTemperature applies the Arrhenius-equivalent MKT definition.
// Invalid absolute temperatures are ignored so one corrupt sensor row does not
// make the entire shipment non-computable.
func MeanKineticTemperature(samples []model.Sample) float64 {
	values := make([]float64, 0, len(samples))
	for _, sample := range samples {
		kelvin := sample.TemperatureC + 273.15
		if kelvin > 0 {
			values = append(values, kelvin)
		}
	}
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, kelvin := range values {
		sum += math.Exp(-activationEnergyJ / (gasConstant * kelvin))
	}
	mean := sum / float64(len(values))
	if mean <= 0 || mean >= 1 {
		return 0
	}
	return -activationEnergyJ/gasConstant/math.Log(mean) - 273.15
}
