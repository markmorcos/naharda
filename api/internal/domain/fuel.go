package domain

import "time"

// FuelPrice is an official pump price for a product (project.md §4).
type FuelPrice struct {
	Product       string    `json:"product"` // gasoline-80 | gasoline-92 | gasoline-95 | diesel
	PriceEGP      float64   `json:"price_egp"`
	EffectiveFrom time.Time `json:"effective_from"`
}

// DefaultFuelPrices is the seed/fallback set, overridable per product via the
// manual_override table. These are PLACEHOLDERS pending the operator's manual
// entry on the next announced EGPC price change (fuel is manual-first — §4).
var DefaultFuelPrices = []FuelPrice{
	{Product: "gasoline-80", PriceEGP: 12.25, EffectiveFrom: fuelDate(2025, 10, 17)},
	{Product: "gasoline-92", PriceEGP: 13.75, EffectiveFrom: fuelDate(2025, 10, 17)},
	{Product: "gasoline-95", PriceEGP: 15.00, EffectiveFrom: fuelDate(2025, 10, 17)},
	{Product: "diesel", PriceEGP: 13.50, EffectiveFrom: fuelDate(2025, 10, 17)},
}

func fuelDate(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
