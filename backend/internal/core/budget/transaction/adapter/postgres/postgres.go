package postgres

import (
	"context"
	"errors"
	"time"

	"backend/core/budget/transaction/port"
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

const tableName = "budget.transactions"

var columns = []string{
	"id",
	"organization_id",
	"account_id",
	"category_id",
	"subcategory_id",
	"budget_id",
	"type",
	"amount",
	"description",
	"external_reference_number",
	"date",
	"created_at",
	"updated_at",
}

var sqlColumnByDomainField = map[string]string{
	"id":                       "id",
	"organizationId":           "organization_id",
	"accountId":                "account_id",
	"categoryId":               "category_id",
	"subcategoryId":            "subcategory_id",
	"budgetId":                 "budget_id",
	"type":                     "type",
	"amount":                   "amount",
	"description":              "description",
	"externalReferenceNumber":  "external_reference_number",
	"date":                     "date",
	"createdAt":                "created_at",
	"updatedAt":                "updated_at",
}

type postgres struct {
	db     dbConn
	logger basedomain.Logger
}

func NewPostgres(db database.PoolInterface, logger basedomain.Logger) port.Repository {
	return postgres{
		db:     db,
		logger: logger.With("component", "transaction.repository"),
	}
}

func (r postgres) WithTx(tx basedomain.Transaction) port.Repository {
	return postgres{
		db:     tx.GetTx(),
		logger: r.logger,
	}
}

func (r postgres) FindOne(ctx context.Context, criteria dafi.Criteria) (port.Transaction, error) {
	query := sqlcraft.Select(columns...).
		From(tableName).
		Where(criteria.Filters...).
		SQLColumnByDomainField(sqlColumnByDomainField).
		Limit(1)

	result, err := query.ToSQL()
	if err != nil {
		return port.Transaction{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	r.logger.WithContext(ctx).Debug("executing query", "sql", result.SQL)

	row := r.db.QueryRow(ctx, result.SQL, result.Args...)

	var txn port.Transaction
	err = row.Scan(
		&txn.ID,
		&txn.OrganizationID,
		&txn.AccountID,
		&txn.CategoryID,
		&txn.SubcategoryID,
		&txn.BudgetID,
		&txn.Type,
		&txn.Amount,
		&txn.Description,
		&txn.ExternalReferenceNumber,
		&txn.Date,
		&txn.CreatedAt,
		&txn.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return port.Transaction{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Code(apperrors.CodeNotFound).Wrap(err)
		}
		return port.Transaction{}, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
	}

	return txn, nil
}

func (r postgres) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[port.Transaction], error) {
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

	var txns basedomain.List[port.Transaction]
	for rows.Next() {
		var txn port.Transaction
		err = rows.Scan(
			&txn.ID,
			&txn.OrganizationID,
			&txn.AccountID,
			&txn.CategoryID,
			&txn.SubcategoryID,
			&txn.BudgetID,
			&txn.Type,
			&txn.Amount,
			&txn.Description,
			&txn.ExternalReferenceNumber,
			&txn.Date,
			&txn.CreatedAt,
			&txn.UpdatedAt,
		)
		if err != nil {
			return nil, oops.WithContext(ctx).In(apperrors.LayerRepository).Wrap(err)
		}
		txns = append(txns, txn)
	}

	return txns, nil
}

func (r postgres) Create(ctx context.Context, input port.CreateTransaction) error {
	now := time.Now()

	query := sqlcraft.InsertInto(tableName).
		WithColumns(columns...).
		WithValues(
			input.ID,
			input.OrganizationID,
			input.AccountID,
			input.CategoryID,
			input.SubcategoryID,
			input.BudgetID,
			input.Type,
			input.Amount,
			input.Description,
			input.ExternalReferenceNumber,
			input.Date,
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

func (r postgres) CreateBulk(ctx context.Context, inputs basedomain.List[port.CreateTransaction]) error {
	if inputs.IsEmpty() {
		return nil
	}

	now := time.Now()
	query := sqlcraft.InsertInto(tableName).WithColumns(columns...)

	for _, input := range inputs {
		query = query.WithValues(
			input.ID,
			input.OrganizationID,
			input.AccountID,
			input.CategoryID,
			input.SubcategoryID,
			input.BudgetID,
			input.Type,
			input.Amount,
			input.Description,
			input.ExternalReferenceNumber,
			input.Date,
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

func (r postgres) Update(ctx context.Context, input port.UpdateTransaction, filters ...dafi.Filter) error {
	query := sqlcraft.Update(tableName).
		WithColumns("category_id", "subcategory_id", "budget_id", "type", "amount", "description", "external_reference_number", "date", "updated_at").
		WithValues(
			input.CategoryID,
			input.SubcategoryID,
			input.BudgetID,
			input.Type,
			input.Amount,
			input.Description,
			input.ExternalReferenceNumber,
			input.Date,
			time.Now(),
		).
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
