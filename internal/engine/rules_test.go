package engine

import (
	"testing"

	"example.com/coldchain-lake/internal/model"
)

func TestNormalizeDeduplicatesAndAlignsClock(t *testing.T) {
	rows := []model.Sample{
		{ShipmentID: "s1", LoggerID: "l1", Sequence: 1, TimestampMs: 1000, ReceivedAtMs: 1001, TemperatureC: 4},
		{ShipmentID: "s1", LoggerID: "l1", Sequence: 1, TimestampMs: 1000, ReceivedAtMs: 1002, TemperatureC: 5},
	}
	got := Normalize(rows, map[string]int64{"l1": 500})
	if len(got.Samples) != 1 || got.Samples[0].TimestampMs != 1500 || got.Samples[0].TemperatureC != 5 {
		t.Fatalf("unexpected normalized samples: %+v", got)
	}
}

func TestDetectAndAttributeExcursion(t *testing.T) {
	samples := []model.Sample{
		{ShipmentID: "s1", LoggerID: "l1", Sequence: 1, TimestampMs: 1000, ReceivedAtMs: 1000, TemperatureC: 4},
		{ShipmentID: "s1", LoggerID: "l1", Sequence: 2, TimestampMs: 2000, ReceivedAtMs: 2000, TemperatureC: 9},
		{ShipmentID: "s1", LoggerID: "l1", Sequence: 3, TimestampMs: 3000, ReceivedAtMs: 3000, TemperatureC: 10},
		{ShipmentID: "s1", LoggerID: "l1", Sequence: 4, TimestampMs: 4000, ReceivedAtMs: 4000, TemperatureC: 5},
	}
	events := []model.DoorEvent{{ShipmentID: "s1", EventID: "d1", TimestampMs: 1500, ReceivedAtMs: 1500, Kind: "door_open"}}
	result := Compute(samples, events, Config{MinTempC: 2, MaxTempC: 8, MinExcursionMs: 500, AttributionWindowMs: 1000})
	if len(result.Excursions) != 1 || result.Excursions[0].Attribution != "door_open" {
		t.Fatalf("expected attributed excursion, got %+v", result.Excursions)
	}
}

func TestMeanKineticTemperatureUsesElapsedSamples(t *testing.T) {
	samples := []model.Sample{{TemperatureC: 5}, {TemperatureC: 5}}
	got := MeanKineticTemperature(samples)
	if got < 4.9 || got > 5.1 {
		t.Fatalf("unexpected MKT: %f", got)
	}
}
