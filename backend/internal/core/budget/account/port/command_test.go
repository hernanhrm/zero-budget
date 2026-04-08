package port

import (
	"context"
	"testing"

	"backend/core/budget/account/accounttype"
	"backend/infra/money"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount_Validate_invalidType(t *testing.T) {
	ctx := context.Background()
	c := CreateAccount{
		ID:             uuid.New(),
		OrganizationID: "org-1",
		Name:           "Primary",
		Type:           "NOT_A_TYPE",
		CurrencyCode:   "USD",
		CurrentBalance: money.MajorAmount{},
		IsActive:       true,
	}
	err := c.Validate(ctx)
	require.Error(t, err)
	assert.ErrorContains(t, err, "must be a valid value")
}

func TestCreateAccount_Validate_validType(t *testing.T) {
	ctx := context.Background()
	c := CreateAccount{
		ID:             uuid.New(),
		OrganizationID: "org-1",
		Name:           "Primary",
		Type:           accounttype.Checking,
		CurrencyCode:   "USD",
		CurrentBalance: money.MajorAmount{},
		IsActive:       true,
	}
	require.NoError(t, c.Validate(ctx))
}

func TestUpdateAccount_Validate_typeOptional(t *testing.T) {
	ctx := context.Background()
	u := UpdateAccount{Name: null.StringFrom("Updated name")}
	require.NoError(t, u.Validate(ctx))
}

func TestUpdateAccount_Validate_invalidTypeWhenSet(t *testing.T) {
	ctx := context.Background()
	u := UpdateAccount{Type: null.StringFrom("bad")}
	err := u.Validate(ctx)
	require.Error(t, err)
}

func TestUpdateAccount_Validate_validTypeWhenSet(t *testing.T) {
	ctx := context.Background()
	u := UpdateAccount{Type: null.StringFrom(accounttype.Investment)}
	require.NoError(t, u.Validate(ctx))
}
