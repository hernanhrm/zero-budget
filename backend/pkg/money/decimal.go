package money

import (
	"fmt"
	"math"
)

// defaultMajorDecimalPlaces is used when a JSON number is interpreted as major units (e.g. dollars).
// TODO: derive per currency from ISO 4217 when CreateAccount carries currency context end-to-end.
const defaultMajorDecimalPlaces = 2

// majorFloatToMinor converts a major-unit floating amount (e.g. 12.34) to minor units
// using half-up rounding at decimalPlaces.
func majorFloatToMinor(f float64, decimalPlaces int) (Minor, error) {
	if decimalPlaces < 0 || decimalPlaces > 18 {
		return 0, fmt.Errorf("money: invalid decimalPlaces %d", decimalPlaces)
	}

	factor := math.Pow10(decimalPlaces)
	rounded := math.Round(f * factor)
	if rounded > float64(math.MaxInt64) || rounded < float64(math.MinInt64) {
		return 0, fmt.Errorf("money: amount out of int64 range")
	}

	return Minor(int64(rounded)), nil
}
