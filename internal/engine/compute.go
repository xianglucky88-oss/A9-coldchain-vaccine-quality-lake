package engine

import (
	"example.com/coldchain-lake/internal/model"
)

func Compute(samples []model.Sample, events []model.DoorEvent, config Config) Result {
	config = config.normalized()
	normalized := Normalize(samples, config.LoggerOffsetsMs)
	excursions := DetectExcursions(normalized.Samples, config)
	excursions = AttributeDoors(excursions, events, config.AttributionWindowMs)

	byShipment := make(map[string][]model.Sample)
	for _, sample := range normalized.Samples {
		byShipment[sample.ShipmentID] = append(byShipment[sample.ShipmentID], sample)
	}
	quality := make([]ShipmentQuality, 0, len(byShipment))
	for shipmentID, items := range byShipment {
		count := 0
		for _, excursion := range excursions {
			if excursion.ShipmentID == shipmentID {
				count++
			}
		}
		quality = append(quality, ShipmentQuality{
			ShipmentID:      shipmentID,
			Samples:         len(items),
			RejectedSamples: normalized.RejectedByShipment[shipmentID],
			DuplicateRows:   normalized.DuplicateByShipment[shipmentID],
			Excursions:      count,
			MKT:             MeanKineticTemperature(items),
		})
	}
	sortQuality(quality)
	return Result{Excursions: excursions, Quality: quality}
}
