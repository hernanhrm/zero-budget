import type { Account } from "@workspace/api"
import { formatMinorUnits, minorUnitsFromApi } from "@workspace/money"

function isCheckingType(type: string | undefined): boolean {
	return (type ?? "").toUpperCase() === "CHECKING"
}

function isSavingsType(type: string | undefined): boolean {
	return (type ?? "").toUpperCase() === "SAVINGS"
}

interface MetricCardsProps {
	accounts: Account[]
}

export function MetricCards({ accounts }: MetricCardsProps) {
	const list = Array.isArray(accounts) ? accounts : []

	const displayCurrency = list[0]?.currencyCode ?? "USD"

	const checkingAccounts = list.filter((a) => isCheckingType(a.type))
	const savingsAccounts = list.filter((a) => isSavingsType(a.type))

	const totalMinor = minorUnitsFromApi(
		list.reduce((sum, a) => sum + (a.currentBalance ?? 0), 0),
	)

	const checkingMinor = minorUnitsFromApi(
		checkingAccounts.reduce((sum, a) => sum + (a.currentBalance ?? 0), 0),
	)

	const savingsMinor = minorUnitsFromApi(
		savingsAccounts.reduce((sum, a) => sum + (a.currentBalance ?? 0), 0),
	)

	return (
		<div className="grid w-full gap-5 md:grid-cols-3">
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
					CHECKING ACCOUNTS
				</p>
				<p className="font-space-grotesk text-4xl font-bold text-foreground">
					{formatMinorUnits(checkingMinor, displayCurrency)}
				</p>
				<p className="font-ibm-plex-mono text-[11px] tracking-[1px] text-muted-foreground">
					{checkingAccounts.length}{" "}
					{checkingAccounts.length === 1 ? "ACCOUNT" : "ACCOUNTS"}
				</p>
			</div>
			<div className="flex flex-col gap-4 border border-border p-6">
				<p className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
					SAVINGS ACCOUNTS
				</p>
				<p className="font-space-grotesk text-4xl font-bold text-primary">
					{formatMinorUnits(savingsMinor, displayCurrency)}
				</p>
				<p className="font-ibm-plex-mono text-[11px] tracking-[1px] text-muted-foreground">
					{savingsAccounts.length}{" "}
					{savingsAccounts.length === 1 ? "ACCOUNT" : "ACCOUNTS"}
				</p>
			</div>
		</div>
	)
}
