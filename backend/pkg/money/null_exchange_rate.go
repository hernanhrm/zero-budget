package money

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// NullExchangeRate is an [ExchangeRate] that may be null (omitted or JSON null).
// Used for partial PATCH/PUT bodies and sqlcraft COALESCE updates.
type NullExchangeRate struct {
	Rate  ExchangeRate
	Valid bool
}

// NullExchangeRateFrom returns a non-null wrapper around r.
func NullExchangeRateFrom(r ExchangeRate) NullExchangeRate {
	return NullExchangeRate{Rate: r, Valid: true}
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *NullExchangeRate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return fmt.Errorf("money: UnmarshalJSON NullExchangeRate: %w", err)
	}
	r, err := ParseExchangeRate(f)
	if err != nil {
		return fmt.Errorf("money: UnmarshalJSON NullExchangeRate: %w", err)
	}
	n.Rate = r
	n.Valid = true
	return nil
}

// MarshalJSON encodes null when invalid.
func (n NullExchangeRate) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Rate)
}

// Value implements driver.Valuer for partial SQL updates.
func (n NullExchangeRate) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Rate.Value()
}
