package core

import (
	"context"

	currencypkg "backend/core/budget/currency/port"
	"backend/core/budget/organization_currency/port"
	basedomain "backend/port"
	apperrors "backend/port/errors"
	"backend/infra/dafi"
	"github.com/samber/oops"
)

var allowedOrganizationCurrencyRelations = map[string]struct{}{
	port.RelationCurrencies: {},
}

type service struct {
	repo         port.Repository
	currencyRepo currencypkg.Repository
	logger       basedomain.Logger
}

func New(repo port.Repository, currencyRepo currencypkg.Repository, logger basedomain.Logger) port.Service {
	return service{
		repo:         repo,
		currencyRepo: currencyRepo,
		logger:       logger.With("component", "organization_currency.service"),
	}
}

func (s service) WithTx(tx basedomain.Transaction) port.Service {
	return service{
		repo:         s.repo.WithTx(tx),
		currencyRepo: s.currencyRepo,
		logger:       s.logger,
	}
}

func (s service) FindOne(ctx context.Context, criteria dafi.Criteria) (port.OrganizationCurrency, error) {
	if err := dafi.ValidateRelations(criteria.Relations, allowedOrganizationCurrencyRelations); err != nil {
		return port.OrganizationCurrency{}, oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	repoCrit := criteria
	repoCrit.Relations = nil

	oc, err := s.repo.FindOne(ctx, repoCrit)
	if err != nil {
		return port.OrganizationCurrency{}, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	if !relationsWantCurrencies(criteria.Relations) {
		oc.Currency = nil
		return oc, nil
	}

	cur, err := s.currencyRepo.FindOne(ctx, dafi.Where("code", dafi.Equal, oc.CurrencyCode))
	if err != nil {
		oc.Currency = nil
		return oc, nil
	}
	oc.Currency = organizationCurrencyCurrencyFromRepo(cur)
	return oc, nil
}

func (s service) FindAll(ctx context.Context, criteria dafi.Criteria) (basedomain.List[port.OrganizationCurrency], error) {
	if err := dafi.ValidateRelations(criteria.Relations, allowedOrganizationCurrencyRelations); err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	repoCrit := criteria
	repoCrit.Relations = nil

	ocs, err := s.repo.FindAll(ctx, repoCrit)
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	if !relationsWantCurrencies(criteria.Relations) {
		for i := range ocs {
			ocs[i].Currency = nil
		}
		return ocs, nil
	}

	codes := distinctCurrencyCodes(ocs)
	if len(codes) == 0 {
		for i := range ocs {
			ocs[i].Currency = nil
		}
		return ocs, nil
	}

	currencies, err := s.currencyRepo.FindAll(ctx, dafi.Where("code", dafi.In, codes))
	if err != nil {
		return nil, oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	byCode := make(map[string]currencypkg.Currency, len(currencies))
	for _, c := range currencies {
		byCode[c.Code] = c
	}

	for i := range ocs {
		c, ok := byCode[ocs[i].CurrencyCode]
		if !ok {
			ocs[i].Currency = nil
			continue
		}
		ocs[i].Currency = organizationCurrencyCurrencyFromRepo(c)
	}

	return ocs, nil
}

func distinctCurrencyCodes(ocs []port.OrganizationCurrency) []string {
	seen := make(map[string]struct{}, len(ocs))
	out := make([]string, 0, len(ocs))
	for _, oc := range ocs {
		if oc.CurrencyCode == "" {
			continue
		}
		if _, ok := seen[oc.CurrencyCode]; ok {
			continue
		}
		seen[oc.CurrencyCode] = struct{}{}
		out = append(out, oc.CurrencyCode)
	}
	return out
}

func (s service) Create(ctx context.Context, input port.CreateOrganizationCurrency) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Create(ctx, input); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("organization currency created", "currencyCode", input.CurrencyCode)

	return nil
}

func (s service) CreateBulk(ctx context.Context, inputs basedomain.List[port.CreateOrganizationCurrency]) error {
	if err := s.repo.CreateBulk(ctx, inputs); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("organization currencies created", "count", len(inputs))

	return nil
}

func (s service) Update(ctx context.Context, input port.UpdateOrganizationCurrency, filters ...dafi.Filter) error {
	if err := input.Validate(ctx); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Code(apperrors.CodeValidation).Wrap(err)
	}

	if err := s.repo.Update(ctx, input, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("organization currency updated")

	return nil
}

func (s service) Delete(ctx context.Context, filters ...dafi.Filter) error {
	if err := s.repo.Delete(ctx, filters...); err != nil {
		return oops.WithContext(ctx).In(apperrors.LayerService).Wrap(err)
	}

	s.logger.WithContext(ctx).Info("organization currency deleted")

	return nil
}
