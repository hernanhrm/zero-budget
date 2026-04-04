// Package money represents monetary amounts as integer minor units (e.g. cents for USD),
// matching JSON numbers on the wire and the approach described in
// https://github.com/techforge-lat/money-as-integer
package money

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// Minor is a monetary amount in the smallest currency unit for the associated ISO 4217 code.
type Minor int64

// Int64 returns the wire / database representation.
func (m Minor) Int64() int64 {
	return int64(m)
}

// MarshalJSON encodes Minor as a JSON number.
func (m Minor) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(m))
}

// UnmarshalJSON decodes a JSON integer into minor units (API responses and internal use).
// For create-account bodies that send major units (e.g. 40 dollars), use MajorAmount instead.
func (m *Minor) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*m = 0
		return nil
	}
	var n int64
	if err := json.Unmarshal(data, &n); err != nil {
		return fmt.Errorf("money: UnmarshalJSON Minor: %w", err)
	}
	*m = Minor(n)
	return nil
}

// Scan implements sql.Scanner for database reads.
func (m *Minor) Scan(src interface{}) error {
	if src == nil {
		*m = 0
		return nil
	}

	var n int64
	switch v := src.(type) {
	case int64:
		n = v
	case int32:
		n = int64(v)
	case int16:
		n = int64(v)
	case int:
		n = int64(v)
	case uint64:
		if v > uint64(math.MaxInt64) {
			return fmt.Errorf("money: uint64 value %d overflows Minor", v)
		}
		n = int64(v)
	case uint32:
		n = int64(v)
	case []byte:
		parsed, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return fmt.Errorf("money: scan []byte: %w", err)
		}
		n = parsed
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("money: scan string: %w", err)
		}
		n = parsed
	default:
		return fmt.Errorf("money: cannot scan %T into Minor", src)
	}

	*m = Minor(n)
	return nil
}

// Value implements driver.Valuer for database writes.
func (m Minor) Value() (driver.Value, error) {
	return int64(m), nil
}
