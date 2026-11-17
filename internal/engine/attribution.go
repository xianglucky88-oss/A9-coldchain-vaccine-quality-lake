package engine

import "example.com/coldchain-lake/internal/model"

func AttributeDoors(excursions []Excursion, events []model.DoorEvent, windowMs int64) []Excursion {
	if windowMs < 0 {
		windowMs = 0
	}
	result := append([]Excursion(nil), excursions...)
	for i := range result {
		ids := make([]string, 0)
		for _, event := range events {
			if event.ShipmentID != result[i].ShipmentID || event.Kind != "door_open" {
				continue
			}
			if event.TimestampMs >= result[i].StartMs-windowMs && event.TimestampMs <= result[i].EndMs {
				ids = append(ids, event.EventID)
			}
		}
		if len(ids) > 0 {
			result[i].DoorEventIDs = ids
			result[i].Attribution = "door_open"
		}
	}
	return result
}
