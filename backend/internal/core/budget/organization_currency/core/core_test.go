package core

import (
	"context"
	"testing"
	"time"

	currencypkg "backend/core/budget/currency/port"
	"backend/core/budget/organization_currency/port"
	basedomain "backend/port"
	"backend/infra/dafi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockOrgRepo struct {
	mock.Mock
}

func (m *mockOrgRepo) FindOne(ctx context.Context, criteria dafi.Criteria) (port.OrganizationCurrency, error) {
	args := m.Called(ctx, criteria)
	return args.Get(0).(port.OrganizationCurrency), args.Error(1)
}

func (m *mockOrgRepo) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[port.OrganizationCurrency], error) {
	args := m.Called(ctx, criteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(basedomain.List[port.OrganizationCurrency]), args.Error(1)
}

func (m *mockOrgRepo) Create(ctx context.Context, entity port.CreateOrganizationCurrency) error {
	return m.Called(ctx, entity).Error(0)
}

func (m *mockOrgRepo) CreateBulk(ctx context.Context, entities basedomain.List[port.CreateOrganizationCurrency]) error {
	return m.Called(ctx, entities).Error(0)
}

func (m *mockOrgRepo) Update(ctx context.Context, entity port.UpdateOrganizationCurrency, filters ...dafi.Filter) error {
	return m.Called(ctx, entity, filters).Error(0)
}

func (m *mockOrgRepo) Delete(ctx context.Context, filters ...dafi.Filter) error {
	return m.Called(ctx, filters).Error(0)
}

func (m *mockOrgRepo) WithTx(tx basedomain.Transaction) port.Repository {
	args := m.Called(tx)
	return args.Get(0).(port.Repository)
}

type mockCurrencyRepo struct {
	mock.Mock
}

func (m *mockCurrencyRepo) FindOne(ctx context.Context, criteria dafi.Criteria) (currencypkg.Currency, error) {
	args := m.Called(ctx, criteria)
	return args.Get(0).(currencypkg.Currency), args.Error(1)
}

func (m *mockCurrencyRepo) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[currencypkg.Currency], error) {
	args := m.Called(ctx, criteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(basedomain.List[currencypkg.Currency]), args.Error(1)
}

func TestService_FindAll_invalidRelation(t *testing.T) {
	org := new(mockOrgRepo)
	cur := new(mockCurrencyRepo)
	svc := New(org, cur, noopLogger{})

	_, err := svc.FindAll(context.Background(), dafi.Criteria{
		Relations: []string{"unknown"},
	})
	require.Error(t, err)
	org.AssertNotCalled(t, "FindAll")
	cur.AssertNotCalled(t, "FindAll")
}

func TestService_FindAll_noRelations_currencyNil(t *testing.T) {
	org := new(mockOrgRepo)
	cur := new(mockCurrencyRepo)
	svc := New(org, cur, noopLogger{})

	now := time.Now()
	row := port.OrganizationCurrency{
		ID:             uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		OrganizationID: "org1",
		CurrencyCode:   "USD",
		IsBase:         true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	org.On("FindAll", mock.Anything, mock.MatchedBy(func(c dafi.Criteria) bool {
		return len(c.Relations) == 0
	})).Return(basedomain.List[port.OrganizationCurrency]{row}, nil)

	out, err := svc.FindAll(context.Background(), dafi.Criteria{})
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Nil(t, out[0].Currency)
	cur.AssertNotCalled(t, "FindAll")
	org.AssertExpectations(t)
}

func TestService_FindAll_withCurrencies_mapsNested(t *testing.T) {
	org := new(mockOrgRepo)
	cur := new(mockCurrencyRepo)
	svc := New(org, cur, noopLogger{})

	now := time.Now()
	usd := port.OrganizationCurrency{
		ID:             uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		OrganizationID: "org1",
		CurrencyCode:   "USD",
		IsBase:         true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	eur := port.OrganizationCurrency{
		ID:             uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		OrganizationID: "org1",
		CurrencyCode:   "EUR",
		IsBase:         false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	org.On("FindAll", mock.Anything, mock.MatchedBy(func(c dafi.Criteria) bool {
		return len(c.Relations) == 0
	})).Return(basedomain.List[port.OrganizationCurrency]{usd, eur}, nil)

	cur.On("FindAll", mock.Anything, mock.MatchedBy(func(c dafi.Criteria) bool {
		if len(c.Filters) != 1 || c.Filters[0].Operator != dafi.In {
			return false
		}
		v, ok := c.Filters[0].Value.([]string)
		if !ok || len(v) != 2 {
			return false
		}
		return true
	})).Return(basedomain.List[currencypkg.Currency]{
		{Code: "USD", Name: "US Dollar", Symbol: "$", DecimalPlaces: 2},
		{Code: "EUR", Name: "Euro", Symbol: "€", DecimalPlaces: 2},
	}, nil)

	out, err := svc.FindAll(context.Background(), dafi.Criteria{
		Relations: []string{port.RelationCurrencies},
	})
	require.NoError(t, err)
	require.Len(t, out, 2)
	require.NotNil(t, out[0].Currency)
	assert.Equal(t, "USD", out[0].Currency.Code)
	assert.Equal(t, "US Dollar", out[0].Currency.Name)
	require.NotNil(t, out[1].Currency)
	assert.Equal(t, "EUR", out[1].Currency.Code)
	org.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestService_FindOne_withCurrencies(t *testing.T) {
	org := new(mockOrgRepo)
	cur := new(mockCurrencyRepo)
	svc := New(org, cur, noopLogger{})

	now := time.Now()
	row := port.OrganizationCurrency{
		ID:             uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		OrganizationID: "org1",
		CurrencyCode:   "USD",
		IsBase:         true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	org.On("FindOne", mock.Anything, mock.MatchedBy(func(c dafi.Criteria) bool {
		return len(c.Relations) == 0
	})).Return(row, nil)

	cur.On("FindOne", mock.Anything, mock.Anything).Return(
		currencypkg.Currency{Code: "USD", Name: "US Dollar", Symbol: "$", DecimalPlaces: 2},
		nil,
	)

	out, err := svc.FindOne(context.Background(), dafi.Criteria{
		Relations: []string{port.RelationCurrencies},
	})
	require.NoError(t, err)
	require.NotNil(t, out.Currency)
	assert.Equal(t, "$", out.Currency.Symbol)
	org.AssertExpectations(t)
	cur.AssertExpectations(t)
}

type noopLogger struct{}

func (noopLogger) With(...interface{}) basedomain.Logger { return noopLogger{} }
func (noopLogger) WithContext(context.Context) basedomain.Logger {
	return noopLogger{}
}
func (noopLogger) Debug(string, ...interface{}) {}
func (noopLogger) Info(string, ...interface{})  {}
func (noopLogger) Warn(string, ...interface{})  {}
func (noopLogger) Error(string, ...interface{}) {}
