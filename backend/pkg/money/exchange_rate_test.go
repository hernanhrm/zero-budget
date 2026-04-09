package money

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseExchangeRate_one(t *testing.T) {
	r, err := ParseExchangeRate(1)
	require.NoError(t, err)
	assert.True(t, r.IsOne())
	assert.InDelta(t, 1.0, r.Float64(), 1e-12)
}

func TestParseExchangeRate_positive(t *testing.T) {
	r, err := ParseExchangeRate(0.92)
	require.NoError(t, err)
	assert.False(t, r.IsOne())
	assert.InDelta(t, 0.92, r.Float64(), 1e-9)
}

func TestExchangeRate_JSON_roundTrip(t *testing.T) {
	orig, err := ParseExchangeRate(31)
	require.NoError(t, err)
	b, err := json.Marshal(orig)
	require.NoError(t, err)

	var back ExchangeRate
	require.NoError(t, json.Unmarshal(b, &back))
	assert.Equal(t, orig, back)
}

func TestExchangeRate_Value_scanRoundTrip(t *testing.T) {
	orig, err := ParseExchangeRate(1.2345678901)
	require.NoError(t, err)

	v, err := orig.Value()
	require.NoError(t, err)
	str, ok := v.(string)
	require.True(t, ok)

	var scanned ExchangeRate
	require.NoError(t, scanned.Scan([]byte(str)))
	assert.Equal(t, orig, scanned)
}
