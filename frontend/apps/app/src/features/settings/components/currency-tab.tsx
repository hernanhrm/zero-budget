import { useQueryClient } from "@tanstack/react-query"
import type { Currency, OrganizationCurrency } from "@workspace/api"
import {
	getGetV1CurrenciesQueryKey,
	useGetV1Currencies,
} from "@workspace/api/hooks/currencies/currencies"
import {
	getGetV1OrganizationCurrenciesQueryKey,
	useDeleteV1OrganizationCurrenciesId,
	useGetV1OrganizationCurrencies,
	usePutV1OrganizationCurrenciesId,
} from "@workspace/api/hooks/organization-currencies/organization-currencies"
import {
	getGetV1TransactionsQueryKey,
	useGetV1Transactions,
} from "@workspace/api/hooks/transactions/transactions"
import { Button } from "@workspace/ui/components/button"
import { Input } from "@workspace/ui/components/input"
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@workspace/ui/components/select"
import { toast } from "@workspace/ui/lib/toast"
import { Trash2 } from "lucide-react"
import { useCallback, useEffect, useMemo, useState } from "react"
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@workspace/ui/components/table"
import { normalizeListPayload } from "#/lib/normalize-list-payload"
import { AddCurrencyModal } from "./add-currency-modal"

const orgCurrencyQueryParams = { relations: "currencies" } as const

function formatRate(rate: number | undefined): string {
	if (rate === undefined || Number.isNaN(rate)) return "—"
	return new Intl.NumberFormat(undefined, {
		maximumFractionDigits: 10,
		useGrouping: true,
	}).format(rate)
}

export interface CurrencyTabProps {
	organizationId: string | undefined
}

export function CurrencyTab({ organizationId }: CurrencyTabProps) {
	const queryClient = useQueryClient()
	const [addModalOpen, setAddModalOpen] = useState(false)
	const [displayFormat, setDisplayFormat] = useState("symbol-first")

	const orgCurrenciesQuery = useGetV1OrganizationCurrencies(orgCurrencyQueryParams, {
		query: {
			enabled: Boolean(organizationId),
			queryKey: getGetV1OrganizationCurrenciesQueryKey(orgCurrencyQueryParams),
		},
		fetch: { credentials: "include" },
	})

	const txnProbeParams = { limit: 1, offset: 0 } as const

	const catalogQuery = useGetV1Currencies(undefined, {
		query: {
			enabled: Boolean(organizationId),
			queryKey: getGetV1CurrenciesQueryKey(undefined),
		},
		fetch: { credentials: "include" },
	})

	const transactionsProbe = useGetV1Transactions(txnProbeParams, {
		query: {
			enabled: Boolean(organizationId),
			queryKey: getGetV1TransactionsQueryKey(txnProbeParams),
		},
		fetch: { credentials: "include" },
	})

	const invalidateOrgCurrencies = useCallback(() => {
		void queryClient.invalidateQueries({
			queryKey: getGetV1OrganizationCurrenciesQueryKey(orgCurrencyQueryParams),
		})
	}, [queryClient])

	const putCurrency = usePutV1OrganizationCurrenciesId({
		fetch: { credentials: "include" },
		mutation: { onSuccess: () => invalidateOrgCurrencies() },
	})

	const deleteCurrency = useDeleteV1OrganizationCurrenciesId({
		fetch: { credentials: "include" },
		mutation: { onSuccess: () => invalidateOrgCurrencies() },
	})

	const orgCurrencies = useMemo(() => {
		return normalizeListPayload<OrganizationCurrency>(
			orgCurrenciesQuery.data?.data,
		)
	}, [orgCurrenciesQuery.data])

	const catalog = useMemo(() => {
		return normalizeListPayload<Currency>(catalogQuery.data?.data)
	}, [catalogQuery.data])

	const existingCodes = useMemo(() => {
		return new Set(
			orgCurrencies.map((c) => (c.currencyCode ?? "").toUpperCase()).filter(Boolean),
		)
	}, [orgCurrencies])

	const availableCatalog = useMemo(() => {
		return catalog.filter((c) => {
			const code = (c.code ?? "").toUpperCase()
			return code.length === 3 && !existingCodes.has(code)
		})
	}, [catalog, existingCodes])

	const hasTransactions = useMemo(() => {
		return normalizeListPayload(transactionsProbe.data?.data).length > 0
	}, [transactionsProbe.data])

	const baseRow = useMemo(
		() => orgCurrencies.find((c) => c.isBase),
		[orgCurrencies],
	)

	const primaryCurrencyCode = baseRow?.currencyCode ?? ""

	const primaryOptions = useMemo(() => {
		return [...orgCurrencies].sort((a, b) =>
			(a.currencyCode ?? "").localeCompare(b.currencyCode ?? ""),
		)
	}, [orgCurrencies])

	const lastUpdatedLabel = useMemo(() => {
		let maxMs = 0
		for (const c of orgCurrencies) {
			if (!c.updatedAt) continue
			const t = Date.parse(c.updatedAt)
			if (Number.isFinite(t) && t > maxMs) maxMs = t
		}
		if (maxMs === 0) return null
		const d = new Date(maxMs)
		return d.toLocaleString(undefined, {
			year: "numeric",
			month: "long",
			day: "numeric",
			hour: "2-digit",
			minute: "2-digit",
			timeZone: "UTC",
			timeZoneName: "short",
		})
	}, [orgCurrencies])

	const handlePrimaryChange = async (nextCode: string) => {
		if (!organizationId || hasTransactions) {
			return
		}

		const next = orgCurrencies.find(
			(c) => (c.currencyCode ?? "").toUpperCase() === nextCode.toUpperCase(),
		)
		const currentBase = orgCurrencies.find((c) => c.isBase)
		if (!next?.id || !currentBase?.id) {
			toast.error("Could not update primary currency.")
			return
		}
		if (next.id === currentBase.id) return

		try {
			let r = await putCurrency.mutateAsync({
				id: currentBase.id,
				data: { isBase: false },
			})
			if (r.status !== 204) {
				toast.error("Could not update primary currency.")
				return
			}
			r = await putCurrency.mutateAsync({
				id: next.id,
				data: { isBase: true },
			})
			if (r.status !== 204) {
				toast.error("Could not set new primary currency.")
				return
			}
			toast.success("Primary currency updated.")
		} catch {
			toast.error("Could not update primary currency.")
		}
	}

	const handleRateBlur = async (row: OrganizationCurrency, raw: string) => {
		if (row.isBase) return
		const id = row.id
		if (!id) return

		const cleaned = raw.replace(/,/g, "").trim()
		const n = Number.parseFloat(cleaned === "" ? "NaN" : cleaned)
		if (!Number.isFinite(n) || n <= 0) {
			toast.error("Rate must be a positive number.")
			return
		}
		if (row.rate !== undefined && Math.abs(n - row.rate) < 1e-12) return

		const r = await putCurrency.mutateAsync({ id, data: { rate: n } })
		if (r.status !== 204) {
			toast.error("Could not update rate.")
			return
		}
		toast.success("Rate updated.")
	}

	const handleDelete = async (row: OrganizationCurrency) => {
		if (row.isBase) {
			toast.error("Remove another currency first, or change the primary currency.")
			return
		}
		const id = row.id
		if (!id) return

		const r = await deleteCurrency.mutateAsync({ id })
		if (r.status !== 204) {
			toast.error("Could not remove currency.")
			return
		}
		toast.success("Currency removed.")
	}

	const displayFormats = [
		{ value: "symbol-first", label: "$1,234.56" },
		{ value: "code-first", label: "USD 1,234.56" },
		{ value: "symbol-last", label: "1,234.56$" },
	]

	const loading =
		Boolean(organizationId) &&
		(orgCurrenciesQuery.isLoading || catalogQuery.isLoading)

	if (!organizationId) {
		return (
			<p className="font-space-grotesk text-sm text-muted-foreground">
				Select an active organization to manage currencies.
			</p>
		)
	}

	return (
		<div className="flex flex-col gap-8">
			<div className="flex gap-6">
				<div className="flex flex-1 flex-col gap-2">
					<span className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
						PRIMARY CURRENCY
					</span>
					<Select
						value={primaryCurrencyCode}
						onValueChange={(v) => void handlePrimaryChange(v)}
						disabled={
							hasTransactions ||
							primaryOptions.length === 0 ||
							putCurrency.isPending
						}
					>
						<SelectTrigger className="h-11 w-full px-4 font-space-grotesk text-sm">
							<SelectValue placeholder="No currencies yet" />
						</SelectTrigger>
						<SelectContent>
							{primaryOptions.map((c) => {
								const code = c.currencyCode ?? ""
								const name = c.currency?.name ?? code
								return (
									<SelectItem key={code} value={code}>
										{`${code} — ${name}`}
									</SelectItem>
								)
							})}
						</SelectContent>
					</Select>
					{hasTransactions ? (
						<span className="font-space-grotesk text-[11px] text-muted-foreground">
							Primary currency is locked because this organization has
							transactions.
						</span>
					) : null}
				</div>

				<div className="flex flex-1 flex-col gap-2">
					<span className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
						DISPLAY FORMAT
					</span>
					<Select value={displayFormat} onValueChange={setDisplayFormat}>
						<SelectTrigger className="h-11 w-full px-4 font-space-grotesk text-sm">
							<SelectValue />
						</SelectTrigger>
						<SelectContent>
							{displayFormats.map((f) => (
								<SelectItem key={f.value} value={f.value}>
									{f.label}
								</SelectItem>
							))}
						</SelectContent>
					</Select>
				</div>
			</div>

			<div className="flex flex-col gap-4">
				<div className="flex items-center justify-between">
					<span className="font-space-grotesk text-sm font-bold tracking-[1px] text-foreground">
						ORGANIZATION CURRENCIES
					</span>
					<Button
						variant="outline"
						size="sm"
						onClick={() => setAddModalOpen(true)}
						disabled={availableCatalog.length === 0 || loading}
						className="border-primary font-space-grotesk text-[10px] font-bold tracking-[1px] text-primary hover:bg-primary/10 hover:text-primary"
					>
						<span className="text-sm">+</span>
						ADD CURRENCY
					</Button>
				</div>

				{loading ? (
					<p className="font-space-grotesk text-sm text-muted-foreground">
						Loading currencies…
					</p>
				) : orgCurrencies.length === 0 ? (
					<p className="font-space-grotesk text-sm text-muted-foreground">
						No currencies yet. Add your base currency to get started.
					</p>
				) : (
					<Table>
						<TableHeader>
							<TableRow className="border-b border-border bg-muted/50 hover:bg-muted/50">
								<TableHead className="font-space-grotesk text-[11px] font-bold tracking-[1px]">
									CURRENCY
								</TableHead>
								<TableHead className="w-[100px] font-space-grotesk text-[11px] font-bold tracking-[1px]">
									CODE
								</TableHead>
								<TableHead className="min-w-[140px] font-space-grotesk text-[11px] font-bold tracking-[1px]">
									RATE (PER 1 BASE)
								</TableHead>
								<TableHead className="w-[100px] font-space-grotesk text-[11px] font-bold tracking-[1px]">
									BASE
								</TableHead>
								<TableHead className="w-10" />
							</TableRow>
						</TableHeader>
						<TableBody>
							{orgCurrencies.map((row) => (
								<TableRow key={row.id ?? row.currencyCode}>
									<TableCell className="font-space-grotesk text-[13px] font-semibold text-foreground">
										{row.currency?.name ?? row.currencyCode ?? "—"}
									</TableCell>
									<TableCell className="font-ibm-plex-mono text-xs text-muted-foreground">
										{row.currencyCode ?? "—"}
									</TableCell>
									<TableCell>
										{row.isBase ? (
											<span className="font-ibm-plex-mono text-[13px] text-foreground">
												{formatRate(1)}
											</span>
										) : (
											<RateCell
												initial={row.rate}
												onCommit={(raw) => void handleRateBlur(row, raw)}
												disabled={putCurrency.isPending}
											/>
										)}
									</TableCell>
									<TableCell className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
										{row.isBase ? "YES" : "—"}
									</TableCell>
									<TableCell>
										<Button
											variant="ghost"
											size="icon-sm"
											disabled={row.isBase || deleteCurrency.isPending}
											onClick={() => void handleDelete(row)}
											className="text-muted-foreground hover:text-destructive"
										>
											<Trash2 />
										</Button>
									</TableCell>
								</TableRow>
							))}
						</TableBody>
					</Table>
				)}

				{lastUpdatedLabel ? (
					<span className="font-ibm-plex-mono text-[10px] tracking-[1px] text-muted-foreground">
						{`LAST UPDATED: ${lastUpdatedLabel.toUpperCase()}`}
					</span>
				) : null}
			</div>

			<AddCurrencyModal
				open={addModalOpen}
				onOpenChange={setAddModalOpen}
				organizationId={organizationId}
				availableCurrencies={availableCatalog}
				isFirstOrgCurrency={orgCurrencies.length === 0}
				onSuccess={invalidateOrgCurrencies}
			/>
		</div>
	)
}

function RateCell(props: {
	initial: number | undefined
	onCommit: (raw: string) => void
	disabled: boolean
}) {
	const [value, setValue] = useState(() =>
		props.initial !== undefined ? formatRate(props.initial) : "",
	)

	useEffect(() => {
		setValue(props.initial !== undefined ? formatRate(props.initial) : "")
	}, [props.initial])

	return (
		<Input
			className="h-9 max-w-[180px] font-ibm-plex-mono text-[13px]"
			value={value}
			onChange={(e) => setValue(e.target.value)}
			onBlur={() => props.onCommit(value)}
			disabled={props.disabled}
			type="text"
			inputMode="decimal"
		/>
	)
}
