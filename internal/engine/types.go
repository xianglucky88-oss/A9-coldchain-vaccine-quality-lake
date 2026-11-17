package engine

type Config struct {
	MinTempC            float64
	MaxTempC            float64
	MinExcursionMs      int64
	AttributionWindowMs int64
	LoggerOffsetsMs     map[string]int64
}

func (c Config) normalized() Config {
	if c.MaxTempC <= c.MinTempC {
		c.MaxTempC = c.MinTempC + 1
	}
	if c.MinExcursionMs < 0 {
		c.MinExcursionMs = 0
	}
	if c.AttributionWindowMs < 0 {
		c.AttributionWindowMs = 0
	}
	return c
}

type Excursion struct {
	ShipmentID   string   `json:"shipment_id"`
	StartMs      int64    `json:"start_ms"`
	EndMs        int64    `json:"end_ms"`
	DurationMs   int64    `json:"duration_ms"`
	MinTempC     float64  `json:"min_temp_c"`
	MaxTempC     float64  `json:"max_temp_c"`
	DoorEventIDs []string `json:"door_event_ids,omitempty"`
	Attribution  string   `json:"attribution"`
}

type ShipmentQuality struct {
	ShipmentID      string  `json:"shipment_id"`
	Samples         int     `json:"samples"`
	RejectedSamples int     `json:"rejected_samples"`
	DuplicateRows   int     `json:"duplicate_rows"`
	Excursions      int     `json:"excursions"`
	MKT             float64 `json:"mkt_c"`
}

type Result struct {
	Excursions []Excursion       `json:"excursions"`
	Quality    []ShipmentQuality `json:"quality"`
}
