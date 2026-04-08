package core

import (
	"context"
	"testing"
	"time"

	"backend/core/budget/account/port"
	transactionport "backend/core/budget/transaction/port"
	basedomain "backend/port"
	apperrors "backend/port/errors"
	"backend/infra/dafi"
	"backend/infra/money"

	"github.com/google/uuid"
	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type noopLogger struct{}

func (noopLogger) Debug(string, ...any)                     {}
func (noopLogger) Info(string, ...any)                      {}
func (noopLogger) Warn(string, ...any)                      {}
func (noopLogger) Error(string, ...any)                     {}
func (noopLogger) With(...any) basedomain.Logger             { return noopLogger{} }
func (noopLogger) WithContext(context.Context) basedomain.Logger { return noopLogger{} }

type stubAccountRepo struct {
	findResult port.Account
	findErr    error
	deleteN    int
	deleteErr  error
}

func (s *stubAccountRepo) FindOne(ctx context.Context, criteria dafi.Criteria) (port.Account, error) {
	_ = ctx
	_ = criteria
	return s.findResult, s.findErr
}

func (s *stubAccountRepo) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[port.Account], error) {
	_ = ctx
	_ = criteria
	return nil, nil
}

func (s *stubAccountRepo) Create(ctx context.Context, input port.CreateAccount) error {
	_ = ctx
	_ = input
	return nil
}

func (s *stubAccountRepo) CreateBulk(ctx context.Context, inputs basedomain.List[port.CreateAccount]) error {
	_ = ctx
	_ = inputs
	return nil
}

func (s *stubAccountRepo) Update(ctx context.Context, input port.UpdateAccount, filters ...dafi.Filter) error {
	_ = ctx
	_ = input
	_ = filters
	return nil
}

func (s *stubAccountRepo) Delete(ctx context.Context, filters ...dafi.Filter) error {
	_ = ctx
	_ = filters
	s.deleteN++
	return s.deleteErr
}

func (s *stubAccountRepo) WithTx(basedomain.Transaction) port.Repository { return s }

type stubTxnRepo struct {
	count    int64
	countErr error
}

func (s *stubTxnRepo) CountByAccountID(ctx context.Context, accountID uuid.UUID) (int64, error) {
	_ = ctx
	_ = accountID
	return s.count, s.countErr
}

func (s *stubTxnRepo) FindOne(ctx context.Context, criteria dafi.Criteria) (transactionport.Transaction, error) {
	_ = ctx
	_ = criteria
	return transactionport.Transaction{}, nil
}

func (s *stubTxnRepo) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[transactionport.Transaction], error) {
	_ = ctx
	_ = criteria
	return nil, nil
}

func (s *stubTxnRepo) Create(ctx context.Context, input transactionport.CreateTransaction) error {
	_ = ctx
	_ = input
	return nil
}

func (s *stubTxnRepo) CreateBulk(ctx context.Context, inputs basedomain.List[transactionport.CreateTransaction]) error {
	_ = ctx
	_ = inputs
	return nil
}

func (s *stubTxnRepo) Update(ctx context.Context, input transactionport.UpdateTransaction, filters ...dafi.Filter) error {
	_ = ctx
	_ = input
	_ = filters
	return nil
}

func (s *stubTxnRepo) Delete(ctx context.Context, filters ...dafi.Filter) error {
	_ = ctx
	_ = filters
	return nil
}

func (s *stubTxnRepo) WithTx(basedomain.Transaction) transactionport.Repository { return s }

func TestService_Delete_NoTransactions_Deletes(t *testing.T) {
	t.Parallel()

	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	acctRepo := &stubAccountRepo{
		findResult: port.Account{
			ID:             id,
			OrganizationID: "org-1",
			Name:           "Checking",
			Type:           "CHECKING",
			CurrencyCode:   "USD",
			CurrentBalance: money.Minor(0),
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}
	txnRepo := &stubTxnRepo{count: 0}

	svc := New(acctRepo, txnRepo, noopLogger{})

	err := svc.Delete(context.Background(), dafi.FilterBy("id", dafi.Equal, id.String())...)
	require.NoError(t, err)
	assert.Equal(t, 1, acctRepo.deleteN)
}

func TestService_Delete_HasTransactions_Conflict(t *testing.T) {
	t.Parallel()

	id := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	acctRepo := &stubAccountRepo{
		findResult: port.Account{
			ID:             id,
			OrganizationID: "org-1",
			Name:           "Savings",
			Type:           "SAVINGS",
			CurrencyCode:   "USD",
			CurrentBalance: money.Minor(0),
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}
	txnRepo := &stubTxnRepo{count: 3}

	svc := New(acctRepo, txnRepo, noopLogger{})

	err := svc.Delete(context.Background(), dafi.FilterBy("id", dafi.Equal, id.String())...)
	require.Error(t, err)
	assert.Equal(t, 0, acctRepo.deleteN)

	oopsErr, ok := oops.AsOops(err)
	require.True(t, ok)
	assert.Equal(t, apperrors.CodeConflict, oopsErr.Code())
	assert.NotEmpty(t, oopsErr.Public())
}

func TestService_Delete_FindOneError_NoDelete(t *testing.T) {
	t.Parallel()

	id := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	acctRepo := &stubAccountRepo{findErr: oops.Code(apperrors.CodeNotFound).Errorf("missing")}
	txnRepo := &stubTxnRepo{}

	svc := New(acctRepo, txnRepo, noopLogger{})

	err := svc.Delete(context.Background(), dafi.FilterBy("id", dafi.Equal, id.String())...)
	require.Error(t, err)
	assert.Equal(t, 0, acctRepo.deleteN)
}
