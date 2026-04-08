package postgres

import (
	"context"
	"errors"
	"time"

	"backend/core/budget/organization_currency/port"
	basedomain "backend/port"
	apperrors "backend/port/errors"
	"backend/infra/dafi"
	"backend/adapter/database"
	"backend/infra/sqlcraft"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/samber/oops"
)

type dbConn interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

const tableName = "budget.organization_currencies"

var columns = []string{
	"id",
	"organization_id",
	"currency_code",
	"is_base",
	"created_at",
	"updated_at",
}

var sqlColumnByDomainField = map[string]string{
	"id":             "id",
	"organizationId": "organization_id",
	"currencyCode":   "currency_code",
	"isBase":         "is_base",
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
		logger: logger.With("component", "organization_currency.repository"),
	}
}

func (r postgres) WithTx(tx basedomain.Transaction) port.Repository {
	return postgres{
		db:     tx.GetTx(),
		logger: r.logger,
	}
}

func (r postgres) FindOne(ctx context.Context, criteria dafi.Criteria) (port.OrganizationCurrency, error) {
	query := sqlcraft.Select(columns...).
		From(tableName).
		Where(criteria.Filters...).
		SQLColumnByDomainField(sqlColumnByDomainField).
		Limit(1)

	result, err := query.ToSQL()
	if err != nil {
		return port.OrganizationCurrency{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	row := r.db.QueryRow(ctx, result.SQL, result.Args...)

	var oc port.OrganizationCurrency
	err = row.Scan(
		&oc.ID,
		&oc.OrganizationID,
		&oc.CurrencyCode,
		&oc.IsBase,
		&oc.CreatedAt,
		&oc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return port.OrganizationCurrency{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Code(apperrors.CodeNotFound).Wrap(err)
		}
		return port.OrganizationCurrency{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	return oc, nil
}

func (r postgres) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[port.OrganizationCurrency], error) {
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

	var ocs basedomain.List[port.OrganizationCurrency]
	for rows.Next() {
		var oc port.OrganizationCurrency
		err = rows.Scan(
			&oc.ID,
			&oc.OrganizationID,
			&oc.CurrencyCode,
			&oc.IsBase,
			&oc.CreatedAt,
			&oc.UpdatedAt,
		)
		if err != nil {
			return nil, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
		}
		ocs = append(ocs, oc)
	}

	return ocs, nil
}

func (r postgres) Create(ctx context.Context, input port.CreateOrganizationCurrency) error {
	now := time.Now()

	query := sqlcraft.InsertInto(tableName).
		WithColumns(columns...).
		WithValues(
			input.ID,
			input.OrganizationID,
			input.CurrencyCode,
			input.IsBase,
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

func (r postgres) CreateBulk(ctx context.Context, inputs basedomain.List[port.CreateOrganizationCurrency]) error {
	if inputs.IsEmpty() {
		return nil
	}

	now := time.Now()
	query := sqlcraft.InsertInto(tableName).WithColumns(columns...)

	for _, input := range inputs {
		query = query.WithValues(
			input.ID,
			input.OrganizationID,
			input.CurrencyCode,
			input.IsBase,
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

func (r postgres) Update(ctx context.Context, input port.UpdateOrganizationCurrency, filters ...dafi.Filter) error {
	query := sqlcraft.Update(tableName).
		WithColumns("is_base", "updated_at").
		WithValues(
			input.IsBase,
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
		return oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	return nil
}
