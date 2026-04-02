package postgres

import (
	"context"
	"errors"

	"backend/core/budget/currency/port"
	basedomain "backend/port"
	apperrors "backend/port/errors"
	"backend/infra/dafi"
	"backend/adapter/database"
	"backend/infra/sqlcraft"

	"github.com/jackc/pgx/v5"
	"github.com/samber/oops"
)

type dbConn interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

const tableName = "budget.currencies"

var columns = []string{
	"code",
	"name",
	"symbol",
	"decimal_places",
}

var sqlColumnByDomainField = map[string]string{
	"code":          "code",
	"name":          "name",
	"symbol":        "symbol",
	"decimalPlaces": "decimal_places",
}

type postgres struct {
	db     dbConn
	logger basedomain.Logger
}

func NewPostgres(db database.PoolInterface, logger basedomain.Logger) port.Repository {
	return postgres{
		db:     db,
		logger: logger.With("component", "currency.repository"),
	}
}

func (r postgres) FindOne(ctx context.Context, criteria dafi.Criteria) (port.Currency, error) {
	query := sqlcraft.Select(columns...).
		From(tableName).
		Where(criteria.Filters...).
		SQLColumnByDomainField(sqlColumnByDomainField).
		Limit(1)

	result, err := query.ToSQL()
	if err != nil {
		return port.Currency{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	row := r.db.QueryRow(ctx, result.SQL, result.Args...)

	var currency port.Currency
	err = row.Scan(
		&currency.Code,
		&currency.Name,
		&currency.Symbol,
		&currency.DecimalPlaces,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return port.Currency{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Code(apperrors.CodeNotFound).Wrap(err)
		}
		return port.Currency{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	return currency, nil
}

func (r postgres) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[port.Currency], error) {
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

	var currencies basedomain.List[port.Currency]
	for rows.Next() {
		var currency port.Currency
		err = rows.Scan(
			&currency.Code,
			&currency.Name,
			&currency.Symbol,
			&currency.DecimalPlaces,
		)
		if err != nil {
			return nil, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
		}
		currencies = append(currencies, currency)
	}

	return currencies, nil
}
