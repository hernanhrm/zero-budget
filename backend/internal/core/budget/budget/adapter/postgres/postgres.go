package postgres

import (
	"context"
	"errors"
	"time"

	"backend/core/budget/budget/port"
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

const tableName = "budget.budgets"

var columns = []string{
	"id",
	"organization_id",
	"name",
	"month",
	"year",
	"currency_code",
	"is_active",
	"created_at",
	"updated_at",
}

var sqlColumnByDomainField = map[string]string{
	"id":             "id",
	"organizationId": "organization_id",
	"name":           "name",
	"month":          "month",
	"year":           "year",
	"currencyCode":   "currency_code",
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
		logger: logger.With("component", "budget.repository"),
	}
}

func (r postgres) WithTx(tx basedomain.Transaction) port.Repository {
	return postgres{
		db:     tx.GetTx(),
		logger: r.logger,
	}
}

func (r postgres) FindOne(ctx context.Context, criteria dafi.Criteria) (port.Budget, error) {
	query := sqlcraft.Select(columns...).
		From(tableName).
		Where(criteria.Filters...).
		SQLColumnByDomainField(sqlColumnByDomainField).
		Limit(1)

	result, err := query.ToSQL()
	if err != nil {
		return port.Budget{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	row := r.db.QueryRow(ctx, result.SQL, result.Args...)

	var b port.Budget
	err = row.Scan(
		&b.ID,
		&b.OrganizationID,
		&b.Name,
		&b.Month,
		&b.Year,
		&b.CurrencyCode,
		&b.IsActive,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return port.Budget{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Code(apperrors.CodeNotFound).Wrap(err)
		}
		return port.Budget{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	return b, nil
}

func (r postgres) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[port.Budget], error) {
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

	var budgets basedomain.List[port.Budget]
	for rows.Next() {
		var b port.Budget
		err = rows.Scan(
			&b.ID,
			&b.OrganizationID,
			&b.Name,
			&b.Month,
			&b.Year,
			&b.CurrencyCode,
			&b.IsActive,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			return nil, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
		}
		budgets = append(budgets, b)
	}

	return budgets, nil
}

func (r postgres) Create(ctx context.Context, input port.CreateBudget) error {
	now := time.Now()

	query := sqlcraft.InsertInto(tableName).
		WithColumns(columns...).
		WithValues(
			input.ID,
			input.OrganizationID,
			input.Name,
			input.Month,
			input.Year,
			input.CurrencyCode,
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

func (r postgres) CreateBulk(ctx context.Context, inputs basedomain.List[port.CreateBudget]) error {
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
			input.Month,
			input.Year,
			input.CurrencyCode,
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

func (r postgres) Update(ctx context.Context, input port.UpdateBudget, filters ...dafi.Filter) error {
	query := sqlcraft.Update(tableName).
		WithColumns("name", "is_active", "updated_at").
		WithValues(
			input.Name,
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
