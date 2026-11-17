package model

import "fmt"

type Sample struct {
	ShipmentID   string  `json:"shipment_id"`
	LoggerID     string  `json:"logger_id"`
	Sequence     int     `json:"sequence"`
	TimestampMs  int64   `json:"timestamp_ms"`
	ReceivedAtMs int64   `json:"received_at_ms"`
	TemperatureC float64 `json:"temperature_c"`
}

func (s Sample) Validate() error {
	if s.ShipmentID == "" || s.LoggerID == "" {
		return fmt.Errorf("shipment_id and logger_id are required")
	}
	if s.Sequence < 0 {
		return fmt.Errorf("sequence must be non-negative")
	}
	if s.TimestampMs <= 0 || s.ReceivedAtMs <= 0 {
		return fmt.Errorf("timestamps must be positive")
	}
	return nil
}

type DoorEvent struct {
	ShipmentID   string `json:"shipment_id"`
	EventID      string `json:"event_id"`
	TimestampMs  int64  `json:"timestamp_ms"`
	ReceivedAtMs int64  `json:"received_at_ms"`
	Kind         string `json:"kind"`
	Location     string `json:"location,omitempty"`
}

func (e DoorEvent) Validate() error {
	if e.ShipmentID == "" || e.EventID == "" {
		return fmt.Errorf("shipment_id and event_id are required")
	}
	if e.TimestampMs <= 0 || e.ReceivedAtMs <= 0 {
		return fmt.Errorf("timestamps must be positive")
	}
	if e.Kind != "door_open" && e.Kind != "door_close" && e.Kind != "handover" {
		return fmt.Errorf("unsupported event kind %q", e.Kind)
	}
	return nil
}
