package engine

import (
	"sort"

	"example.com/coldchain-lake/internal/model"
)

func DetectExcursions(samples []model.Sample, config Config) []Excursion {
	config = config.normalized()
	byShipment := make(map[string][]model.Sample)
	for _, sample := range samples {
		byShipment[sample.ShipmentID] = append(byShipment[sample.ShipmentID], sample)
	}

	var result []Excursion
	for shipmentID, items := range byShipment {
		sort.Slice(items, func(i, j int) bool { return items[i].TimestampMs < items[j].TimestampMs })
		var current *Excursion
		for _, sample := range items {
			outside := sample.TemperatureC < config.MinTempC || sample.TemperatureC > config.MaxTempC
			if outside {
				if current == nil {
					current = &Excursion{
						ShipmentID: shipmentID,
						StartMs:    sample.TimestampMs,
						MinTempC:   sample.TemperatureC,
						MaxTempC:   sample.TemperatureC,
					}
				}
				current.EndMs = sample.TimestampMs
				if sample.TemperatureC < current.MinTempC {
					current.MinTempC = sample.TemperatureC
				}
				if sample.TemperatureC > current.MaxTempC {
					current.MaxTempC = sample.TemperatureC
				}
				continue
			}
			if current != nil {
				current.EndMs = sample.TimestampMs
				closeExcursion(current, config.MinExcursionMs, &result)
				current = nil
			}
		}
		if current != nil {
			closeExcursion(current, config.MinExcursionMs, &result)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].StartMs == result[j].StartMs {
			return result[i].ShipmentID < result[j].ShipmentID
		}
		return result[i].StartMs < result[j].StartMs
	})
	return result
}

func closeExcursion(item *Excursion, minimum int64, result *[]Excursion) {
	item.DurationMs = item.EndMs - item.StartMs
	if item.DurationMs < minimum {
		return
	}
	if item.Attribution == "" {
		item.Attribution = "unattributed"
	}
	*result = append(*result, *item)
}
