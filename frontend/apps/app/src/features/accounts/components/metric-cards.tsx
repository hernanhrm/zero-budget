import type { Account } from "@workspace/api"
import { formatMinorUnits, minorUnitsFromApi } from "@workspace/money"
import {
	accountTypeBucket,
	bankingCashTypes,
	creditDebtTypes,
	investmentTypes,
} from "../account-type-options"

function sumBalances(
	accounts: Account[],
	inBucket: (t: string | undefined) => boolean,
): number {
	return accounts.reduce((sum, a) => {
		if (!inBucket(a.type)) {
			return sum
		}
		return sum + (a.currentBalance ?? 0)
	}, 0)
}

interface MetricCardsProps {
	accounts: Account[]
}

export function MetricCards({ accounts }: MetricCardsProps) {
	const list = Array.isArray(accounts) ? accounts : []

	const displayCurrency = list[0]?.currencyCode ?? "USD"

	const totalMinor = minorUnitsFromApi(
		list.reduce((sum, a) => sum + (a.currentBalance ?? 0), 0),
	)

	const bankingMinor = minorUnitsFromApi(
		sumBalances(list, (t) => bankingCashTypes.has((t ?? "").toUpperCase())),
	)
	const investmentMinor = minorUnitsFromApi(
		sumBalances(list, (t) => investmentTypes.has((t ?? "").toUpperCase())),
	)
	const creditMinor = minorUnitsFromApi(
		sumBalances(list, (t) => creditDebtTypes.has((t ?? "").toUpperCase())),
	)

	const bankingCount = list.filter(
		(a) => accountTypeBucket(a.type) === "banking",
	).length
	const investmentCount = list.filter(
		(a) => accountTypeBucket(a.type) === "investment",
	).length
	const creditCount = list.filter(
		(a) => accountTypeBucket(a.type) === "credit",
	).length

	return (
		<div className="grid w-full gap-5 md:grid-cols-2 xl:grid-cols-4">
			<div className="flex flex-col gap-4 border border-border p-6">
				<p className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
					TOTAL BALANCE
				</p>
				<p className="font-space-grotesk text-4xl font-bold text-foreground">
					{formatMinorUnits(totalMinor, displayCurrency)}
				</p>
				<p className="font-ibm-plex-mono text-[11px] tracking-[1px] text-muted-foreground">
					{list.length} {list.length === 1 ? "ACCOUNT" : "ACCOUNTS"}
				</p>
			</div>
			<div className="flex flex-col gap-4 border border-border p-6">
				<p className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
					BANKING &amp; CASH
				</p>
				<p className="font-space-grotesk text-4xl font-bold text-foreground">
					{formatMinorUnits(bankingMinor, displayCurrency)}
				</p>
				<p className="font-ibm-plex-mono text-[11px] tracking-[1px] text-muted-foreground">
					{bankingCount} {bankingCount === 1 ? "ACCOUNT" : "ACCOUNTS"}
				</p>
			</div>
			<div className="flex flex-col gap-4 border border-border p-6">
				<p className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
					INVESTMENTS
				</p>
				<p className="font-space-grotesk text-4xl font-bold text-foreground">
					{formatMinorUnits(investmentMinor, displayCurrency)}
				</p>
				<p className="font-ibm-plex-mono text-[11px] tracking-[1px] text-muted-foreground">
					{investmentCount} {investmentCount === 1 ? "ACCOUNT" : "ACCOUNTS"}
				</p>
			</div>
			<div className="flex flex-col gap-4 border border-border p-6">
				<p className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
					CREDIT &amp; DEBT
				</p>
				<p className="font-space-grotesk text-4xl font-bold text-primary">
					{formatMinorUnits(creditMinor, displayCurrency)}
				</p>
				<p className="font-ibm-plex-mono text-[11px] tracking-[1px] text-muted-foreground">
					{creditCount} {creditCount === 1 ? "ACCOUNT" : "ACCOUNTS"}
				</p>
			</div>
		</div>
	)
}
