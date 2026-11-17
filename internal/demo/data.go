package demo

import (
	"fmt"

	"example.com/coldchain-lake/internal/model"
)

func Samples() []model.Sample {
	samples := make([]model.Sample, 0, 12)
	for i := 0; i < 12; i++ {
		temp := 4.0
		if i >= 7 && i <= 9 {
			temp = 9.5
		}
		samples = append(samples, model.Sample{
			ShipmentID:   "SHP-001",
			LoggerID:     "LOGGER-A",
			Sequence:     i,
			TimestampMs:  int64(1_767_225_600_000 + i*60_000),
			ReceivedAtMs: int64(1_767_225_601_000 + i*60_000),
			TemperatureC: temp,
		})
	}
	return samples
}

func Events() []model.DoorEvent {
	base := int64(1_767_225_600_000)
	return []model.DoorEvent{
		{ShipmentID: "SHP-001", EventID: "DOOR-1", TimestampMs: base + 5*60_000, ReceivedAtMs: base + 5*60_000 + 500, Kind: "door_open", Location: "transfer-hub"},
		{ShipmentID: "SHP-001", EventID: "DOOR-2", TimestampMs: base + 6*60_000, ReceivedAtMs: base + 6*60_000 + 500, Kind: "door_close", Location: "transfer-hub"},
	}
}

func PrintSummary(resultCount int, excursions int) {
	fmt.Printf("demo complete: shipments=%d excursions=%d\n", resultCount, excursions)
}
