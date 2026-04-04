import type { Account } from "@workspace/api"
import { formatMinorUnits, minorUnitsFromApi } from "@workspace/money"

export interface AccountRow {
	id: string
	name: string
	mask: string
	type: string
	institution: string
	balance: string
	balanceIsPrimary: boolean
	isActive: boolean
}

export function mapAccountsToRows(
	accounts: Account[] | null | undefined,
): AccountRow[] {
	if (!accounts || !Array.isArray(accounts) || accounts.length === 0) {
		return []
	}

	let maxBalance = Number.NEGATIVE_INFINITY
	let maxId: string | undefined
	for (const a of accounts) {
		const bal = a.currentBalance ?? 0
		if (bal > maxBalance) {
			maxBalance = bal
			maxId = a.id
		}
	}

	return accounts.map((account) => {
		const currencyCode = account.currencyCode ?? "USD"
		const minor = minorUnitsFromApi(account.currentBalance ?? 0)
		const institution =
			account.institution != null && account.institution !== ""
				? account.institution
				: "—"
		const mask =
			account.accountNumber != null && account.accountNumber !== ""
				? account.accountNumber
				: "—"
		return {
			id: account.id ?? "",
			name: account.name ?? "",
			mask,
			type: (account.type ?? "").toUpperCase(),
			institution,
			balance: formatMinorUnits(minor, currencyCode),
			balanceIsPrimary: account.id === maxId,
			isActive: account.isActive ?? true,
		}
	})
}
