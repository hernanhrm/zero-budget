package accounttype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	assert.Equal(t, Checking, Normalize("checking"))
	assert.Equal(t, CreditCard, Normalize("  credit_card "))
}

func TestIsValid(t *testing.T) {
	assert.True(t, IsValid(Checking))
	assert.True(t, IsValid(LineOfCredit))
	assert.False(t, IsValid("INVALID"))
	assert.False(t, IsValid(""))
}
