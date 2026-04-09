package money

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// ExchangeRateScale is the number of fractional decimal digits stored for an
// [ExchangeRate], matching budget.organization_currencies.rate NUMERIC(20,10).
const ExchangeRateScale = 10

// exchangeRateMultiplier is 10^ExchangeRateScale.
const exchangeRateMultiplier int64 = 10_000_000_000

// ExchangeRate is units of the quoted currency per one unit of the organization
// base currency. Stored as a fixed-scale integer (rate * 10^ExchangeRateScale)
// for exact equality checks (e.g. base row must be exactly 1).
type ExchangeRate int64

// ExchangeRateOne is exactly 1.0 in wire/database terms.
func ExchangeRateOne() ExchangeRate {
	return ExchangeRate(exchangeRateMultiplier)
}

// Float64 returns the rate as a float for JSON APIs and display.
func (r ExchangeRate) Float64() float64 {
	return float64(r) / float64(exchangeRateMultiplier)
}

// IsOne reports whether r is exactly 1 (base currency row).
func (r ExchangeRate) IsOne() bool {
	return r == ExchangeRateOne()
}

// IsPositive reports whether r is greater than zero.
func (r ExchangeRate) IsPositive() bool {
	return int64(r) > 0
}

// ParseExchangeRate converts a JSON/API float to a fixed-scale [ExchangeRate].
func ParseExchangeRate(f float64) (ExchangeRate, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, fmt.Errorf("money: exchange rate is not a finite number")
	}
	scaled := math.Round(f * float64(exchangeRateMultiplier))
	if scaled > float64(math.MaxInt64) || scaled < float64(math.MinInt64) {
		return 0, fmt.Errorf("money: exchange rate out of range")
	}
	return ExchangeRate(int64(scaled)), nil
}

// MarshalJSON encodes as a JSON number.
func (r ExchangeRate) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Float64())
}

// UnmarshalJSON decodes a JSON number into a fixed-scale rate.
func (r *ExchangeRate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*r = 0
		return nil
	}
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return fmt.Errorf("money: UnmarshalJSON ExchangeRate: %w", err)
	}
	parsed, err := ParseExchangeRate(f)
	if err != nil {
		return fmt.Errorf("money: UnmarshalJSON ExchangeRate: %w", err)
	}
	*r = parsed
	return nil
}

// Scan implements sql.Scanner for PostgreSQL NUMERIC (as string or []byte from pgx).
func (r *ExchangeRate) Scan(src interface{}) error {
	if src == nil {
		*r = 0
		return nil
	}
	var str string
	switch v := src.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("money: cannot scan %T into ExchangeRate", src)
	}
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("money: scan ExchangeRate: %w", err)
	}
	parsed, err := ParseExchangeRate(f)
	if err != nil {
		return err
	}
	*r = parsed
	return nil
}

// Value implements driver.Valuer for PostgreSQL NUMERIC.
func (r ExchangeRate) Value() (driver.Value, error) {
	x := int64(r)
	neg := x < 0
	if neg {
		x = -x
	}
	intPart := x / exchangeRateMultiplier
	frac := x % exchangeRateMultiplier
	if neg {
		return fmt.Sprintf("-%d.%010d", intPart, frac), nil
	}
	return fmt.Sprintf("%d.%010d", intPart, frac), nil
}
