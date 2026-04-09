package core

import (
	currencypkg "backend/core/budget/currency/port"
	"backend/core/budget/organization_currency/port"
)

func organizationCurrencyCurrencyFromRepo(c currencypkg.Currency) *port.OrganizationCurrencyCurrency {
	return &port.OrganizationCurrencyCurrency{
		Code:          c.Code,
		Name:          c.Name,
		Symbol:        c.Symbol,
		DecimalPlaces: c.DecimalPlaces,
	}
}

func relationsWantCurrencies(relations []string) bool {
	for _, r := range relations {
		if r == port.RelationCurrencies {
			return true
		}
	}
	return false
}
