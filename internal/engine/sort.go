package engine

import "sort"

func sortQuality(items []ShipmentQuality) {
	sort.Slice(items, func(i, j int) bool { return items[i].ShipmentID < items[j].ShipmentID })
}
