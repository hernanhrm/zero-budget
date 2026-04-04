package money

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMajorAmount_UnmarshalJSON_wholeMajorUnits(t *testing.T) {
	var s struct {
		B MajorAmount `json:"currentBalance"`
	}
	require.NoError(t, json.Unmarshal([]byte(`{"currentBalance":40}`), &s))
	assert.Equal(t, Minor(4000), s.B.Minor)
}

func TestMajorAmount_UnmarshalJSON_fractionalMajorUnits(t *testing.T) {
	var s struct {
		B MajorAmount `json:"currentBalance"`
	}
	require.NoError(t, json.Unmarshal([]byte(`{"currentBalance":12.34}`), &s))
	assert.Equal(t, Minor(1234), s.B.Minor)
}

func TestMajorAmount_UnmarshalJSON_zero(t *testing.T) {
	var s struct {
		B MajorAmount `json:"currentBalance"`
	}
	require.NoError(t, json.Unmarshal([]byte(`{"currentBalance":0}`), &s))
	assert.Equal(t, Minor(0), s.B.Minor)
}
