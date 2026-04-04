package money

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// MajorAmount is a create-request balance in major display units (e.g. 40 = 40.00 in USD).
// JSON always unmarshals as a float and converts to Minor using defaultMajorDecimalPlaces.
// Use this on POST bodies; list/get responses use Minor (integer minor units) instead.
type MajorAmount struct {
	Minor Minor
}

func (m *MajorAmount) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		m.Minor = 0
		return nil
	}
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return fmt.Errorf("money: UnmarshalJSON MajorAmount: %w", err)
	}
	mv, err := majorFloatToMinor(f, defaultMajorDecimalPlaces)
	if err != nil {
		return fmt.Errorf("money: UnmarshalJSON MajorAmount: %w", err)
	}
	m.Minor = mv
	return nil
}

// MarshalJSON encodes the resolved minor units (for logging/tests; create responses omit body).
func (m MajorAmount) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(m.Minor))
}

// Value implements driver.Valuer for inserts.
func (m MajorAmount) Value() (driver.Value, error) {
	return m.Minor.Value()
}
