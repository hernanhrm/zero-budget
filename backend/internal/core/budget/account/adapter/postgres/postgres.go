package postgres

import (
	"context"
	"errors"
	"time"

	"backend/adapter/database"
	"backend/core/budget/account/port"
	"backend/infra/dafi"
	"backend/infra/sqlcraft"
	basedomain "backend/port"
	apperrors "backend/port/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
)

type dbConn interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

const tableName = "budget.accounts"

var columns = []string{
	"id",
	"organization_id",
	"name",
	"type",
	"institution",
	"account_number",
	"currency_code",
	"current_balance",
	"is_active",
	"created_at",
	"updated_at",
}

var sqlColumnByDomainField = map[string]string{
	"id":             "id",
	"organizationId": "organization_id",
	"name":           "name",
	"type":           "type",
	"institution":    "institution",
	"accountNumber":  "account_number",
	"currencyCode":   "currency_code",
	"currentBalance": "current_balance",
	"isActive":       "is_active",
	"createdAt":      "created_at",
	"updatedAt":      "updated_at",
}

type postgres struct {
	db     dbConn
	logger basedomain.Logger
}

func NewPostgres(db database.PoolInterface, logger basedomain.Logger) port.Repository {
	return postgres{
		db:     db,
		logger: logger.With("component", "account.repository"),
	}
}

func (r postgres) WithTx(tx basedomain.Transaction) port.Repository {
	return postgres{
		db:     tx.GetTx(),
		logger: r.logger,
	}
}

func (r postgres) FindOne(ctx context.Context, criteria dafi.Criteria) (port.Account, error) {
	query := sqlcraft.Select(columns...).
		From(tableName).
		Where(criteria.Filters...).
		SQLColumnByDomainField(sqlColumnByDomainField).
		Limit(1)

	result, err := query.ToSQL()
	if err != nil {
		return port.Account{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	row := r.db.QueryRow(ctx, result.SQL, result.Args...)

	var acct port.Account
	err = row.Scan(
		&acct.ID,
		&acct.OrganizationID,
		&acct.Name,
		&acct.Type,
		&acct.Institution,
		&acct.AccountNumber,
		&acct.CurrencyCode,
		&acct.CurrentBalance,
		&acct.IsActive,
		&acct.CreatedAt,
		&acct.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return port.Account{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Code(apperrors.CodeNotFound).Wrap(err)
		}
		return port.Account{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	return acct, nil
}

func (r postgres) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[port.Account], error) {
	query := sqlcraft.Select(columns...).
		From(tableName).
		Where(criteria.Filters...).
		OrderBy(criteria.Sorts...).
		Limit(criteria.Pagination.PageSize).
		Page(criteria.Pagination.PageNumber).
		SQLColumnByDomainField(sqlColumnByDomainField)

	result, err := query.ToSQL()
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	rows, err := r.db.Query(ctx, result.SQL, result.Args...)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}
	defer rows.Close()

	var accts basedomain.List[port.Account]
	for rows.Next() {
		var acct port.Account
		err = rows.Scan(
			&acct.ID,
			&acct.OrganizationID,
			&acct.Name,
			&acct.Type,
			&acct.Institution,
			&acct.AccountNumber,
			&acct.CurrencyCode,
			&acct.CurrentBalance,
			&acct.IsActive,
			&acct.CreatedAt,
			&acct.UpdatedAt,
		)
		if err != nil {
			return nil, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
		}
		accts = append(accts, acct)
	}

	return accts, nil
}

func (r postgres) Create(ctx context.Context, input port.CreateAccount) error {
	now := time.Now()

	query := sqlcraft.InsertInto(tableName).
		WithColumns(columns...).
		WithValues(
			input.ID,
			input.OrganizationID,
			input.Name,
			input.Type,
			input.Institution,
			input.AccountNumber,
			input.CurrencyCode,
			input.CurrentBalance.Minor,
			input.IsActive,
			now,
			now,
		)

	result, err := query.ToSQL()
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	_, err = r.db.Exec(ctx, result.SQL, result.Args...)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	return nil
}

func (r postgres) CreateBulk(ctx context.Context, inputs basedomain.List[port.CreateAccount]) error {
	if inputs.IsEmpty() {
		return nil
	}

	now := time.Now()
	query := sqlcraft.InsertInto(tableName).WithColumns(columns...)

	for _, input := range inputs {
		query = query.WithValues(
			input.ID,
			input.OrganizationID,
			input.Name,
			input.Type,
			input.Institution,
			input.AccountNumber,
			input.CurrencyCode,
			input.CurrentBalance.Minor,
			input.IsActive,
			now,
			now,
		)
	}

	result, err := query.ToSQL()
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("bulk insert", "sql", result.SQL, "count", len(inputs))

	_, err = r.db.Exec(ctx, result.SQL, result.Args...)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	return nil
}

func (r postgres) Update(ctx context.Context, input port.UpdateAccount, filters ...dafi.Filter) error {
	query := sqlcraft.Update(tableName).
		WithColumns("name", "type", "institution", "account_number", "currency_code", "current_balance", "is_active", "updated_at").
		WithValues(
			input.Name,
			input.Type,
			input.Institution,
			input.AccountNumber,
			input.CurrencyCode,
			input.CurrentBalance,
			input.IsActive,
			time.Now(),
		).
		Where(filters...).
		SQLColumnByDomainField(sqlColumnByDomainField).
		WithPartialUpdate()

	result, err := query.ToSQL()
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	_, err = r.db.Exec(ctx, result.SQL, result.Args...)
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	return nil
}

const pgErrForeignKeyViolation = "23503"

func (r postgres) Delete(ctx context.Context, filters ...dafi.Filter) error {
	query := sqlcraft.DeleteFrom(tableName).
		Where(filters...).
		SQLColumnByDomainField(sqlColumnByDomainField)

	result, err := query.ToSQL()
	if err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	_, err = r.db.Exec(ctx, result.SQL, result.Args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgErrForeignKeyViolation {
			return oops.WithContext(ctx).In(apperrors.LayerRepository).
				Code(apperrors.CodeConflict).
				Public("This account has transactions. Disable it instead of deleting.").
				Wrap(err)
		}
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	return nil
}
