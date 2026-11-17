package engine

import (
	"sort"

	"example.com/coldchain-lake/internal/model"
)

type Normalized struct {
	Samples             []model.Sample
	Rejected            int
	Duplicate           int
	RejectedByShipment  map[string]int
	DuplicateByShipment map[string]int
}

func Normalize(samples []model.Sample, offsets map[string]int64) Normalized {
	index := make(map[string]model.Sample)
	rejected := 0
	duplicate := 0
	rejectedByShipment := make(map[string]int)
	duplicateByShipment := make(map[string]int)
	for _, sample := range samples {
		shipmentKey := sample.ShipmentID
		if shipmentKey == "" {
			shipmentKey = "(unknown)"
		}
		if err := sample.Validate(); err != nil {
			rejected++
			rejectedByShipment[shipmentKey]++
			continue
		}
		sample.TimestampMs += offsets[sample.LoggerID]
		key := sampleIdentity(sample)
		if old, ok := index[key]; ok {
			duplicate++
			duplicateByShipment[shipmentKey]++
			if sample.ReceivedAtMs <= old.ReceivedAtMs {
				continue
			}
		}
		index[key] = sample
	}
	out := make([]model.Sample, 0, len(index))
	for _, sample := range index {
		out = append(out, sample)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].ShipmentID == out[j].ShipmentID {
			if out[i].TimestampMs == out[j].TimestampMs {
				return out[i].LoggerID < out[j].LoggerID
			}
			return out[i].TimestampMs < out[j].TimestampMs
		}
		return out[i].ShipmentID < out[j].ShipmentID
	})
	return Normalized{
		Samples:             out,
		Rejected:            rejected,
		Duplicate:           duplicate,
		RejectedByShipment:  rejectedByShipment,
		DuplicateByShipment: duplicateByShipment,
	}
}

func sampleIdentity(sample model.Sample) string {
	return sample.ShipmentID + "|" + sample.LoggerID + "|" + itoa(sample.Sequence)
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for value > 0 {
		pos--
		buf[pos] = byte('0' + value%10)
		value /= 10
	}
	return string(buf[pos:])
}
