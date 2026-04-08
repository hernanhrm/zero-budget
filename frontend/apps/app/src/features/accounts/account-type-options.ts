import { AccountType } from "@workspace/api"

export type AccountTypeValue = (typeof AccountType)[keyof typeof AccountType]

const allowed = new Set<string>(Object.values(AccountType))

export const ACCOUNT_TYPE_OPTIONS: {
	value: AccountTypeValue
	label: string
}[] = [
	{ value: AccountType.CHECKING, label: "Checking" },
	{ value: AccountType.SAVINGS, label: "Savings" },
	{ value: AccountType.CASH, label: "Cash" },
	{ value: AccountType.CREDIT_CARD, label: "Credit card" },
	{ value: AccountType.MONEY_MARKET, label: "Money market" },
	{ value: AccountType.INVESTMENT, label: "Investment" },
	{ value: AccountType.RETIREMENT, label: "Retirement" },
	{ value: AccountType.LOAN, label: "Loan" },
	{ value: AccountType.MORTGAGE, label: "Mortgage" },
	{ value: AccountType.LINE_OF_CREDIT, label: "Line of credit" },
	{ value: AccountType.OTHER, label: "Other" },
]

export const DEFAULT_ACCOUNT_TYPE: AccountTypeValue = AccountType.CHECKING

/** Maps API / DB value to a form select value; unknown legacy values become OTHER. */
export function accountTypeFromApi(
	raw: string | null | undefined,
): AccountTypeValue {
	const u = (raw ?? "").toUpperCase().trim()
	if (allowed.has(u)) {
		return u as AccountTypeValue
	}
	return AccountType.OTHER
}

const labelByValue: Record<string, string> = Object.fromEntries(
	ACCOUNT_TYPE_OPTIONS.map((o) => [o.value, o.label]),
)

export function formatAccountTypeLabel(type: string | undefined): string {
	const u = (type ?? "").toUpperCase()
	return labelByValue[u] ?? u
}

/** Types rolled into dashboard "Banking & cash" totals. */
export const bankingCashTypes: ReadonlySet<string> = new Set([
	AccountType.CHECKING,
	AccountType.SAVINGS,
	AccountType.CASH,
	AccountType.MONEY_MARKET,
	AccountType.OTHER,
])

/** Investment bucket. */
export const investmentTypes: ReadonlySet<string> = new Set([
	AccountType.INVESTMENT,
	AccountType.RETIREMENT,
])

/** Credit and debt bucket. */
export const creditDebtTypes: ReadonlySet<string> = new Set([
	AccountType.CREDIT_CARD,
	AccountType.LOAN,
	AccountType.MORTGAGE,
	AccountType.LINE_OF_CREDIT,
])

export function accountTypeBucket(
	type: string | undefined,
): "banking" | "investment" | "credit" {
	const u = (type ?? "").toUpperCase()
	if (creditDebtTypes.has(u)) {
		return "credit"
	}
	if (investmentTypes.has(u)) {
		return "investment"
	}
	return "banking"
}
