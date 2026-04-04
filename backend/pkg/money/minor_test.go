package money

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMinor_JSONRoundTrip(t *testing.T) {
	type row struct {
		Balance Minor `json:"currentBalance"`
	}

	const payload = `{"currentBalance":-1099}`

	var got row
	require.NoError(t, json.Unmarshal([]byte(payload), &got))
	assert.Equal(t, Minor(-1099), got.Balance)

	out, err := json.Marshal(got)
	require.NoError(t, err)
	assert.JSONEq(t, payload, string(out))
}

func TestMinor_JSONNull(t *testing.T) {
	var m Minor
	require.NoError(t, json.Unmarshal([]byte("null"), &m))
	assert.Equal(t, Minor(0), m)
}

func TestMinor_Scan(t *testing.T) {
	var m Minor
	require.NoError(t, m.Scan(int64(42)))
	assert.Equal(t, Minor(42), m)

	require.NoError(t, m.Scan([]byte("-7")))
	assert.Equal(t, Minor(-7), m)

	require.NoError(t, m.Scan(nil))
	assert.Equal(t, Minor(0), m)
}

func TestMinor_Value(t *testing.T) {
	v, err := Minor(-99).Value()
	require.NoError(t, err)
	assert.EqualValues(t, int64(-99), v)
}
