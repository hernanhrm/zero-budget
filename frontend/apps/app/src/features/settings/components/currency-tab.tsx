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
import { Button } from "@workspace/ui/components/button"
import { DataTableSectionHeader } from "@workspace/ui/components/data-table-section-header"
import { toast } from "@workspace/ui/lib/toast"
import { Plus } from "lucide-react"
import { useCallback, useMemo, useState } from "react"
import { AccountsDataTable } from "#/features/accounts/components/accounts-data-table"
import { normalizeListPayload } from "#/lib/normalize-list-payload"
import { AddCurrencyModal } from "./add-currency-modal"
import { createOrganizationCurrencyColumns } from "./organization-currencies-columns"

const orgCurrencyQueryParams = { relations: "currencies" } as const

export interface CurrencyTabProps {
	organizationId: string | undefined
}

export function CurrencyTab({ organizationId }: CurrencyTabProps) {
	const queryClient = useQueryClient()
	const [addModalOpen, setAddModalOpen] = useState(false)

	const orgCurrenciesQuery = useGetV1OrganizationCurrencies(
		orgCurrencyQueryParams,
		{
			query: {
				enabled: Boolean(organizationId),
				queryKey: getGetV1OrganizationCurrenciesQueryKey(
					orgCurrencyQueryParams,
				),
			},
			fetch: { credentials: "include" },
		},
	)

	const catalogQuery = useGetV1Currencies(undefined, {
		query: {
			enabled: Boolean(organizationId),
			queryKey: getGetV1CurrenciesQueryKey(undefined),
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
			orgCurrencies
				.map((c) => (c.currencyCode ?? "").toUpperCase())
				.filter(Boolean),
		)
	}, [orgCurrencies])

	const availableCatalog = useMemo(() => {
		return catalog.filter((c) => {
			const code = (c.code ?? "").toUpperCase()
			return code.length === 3 && !existingCodes.has(code)
		})
	}, [catalog, existingCodes])

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

	const handleRateBlur = useCallback(
		async (row: OrganizationCurrency, raw: string) => {
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
		},
		[putCurrency],
	)

	const handleDelete = useCallback(
		async (row: OrganizationCurrency) => {
			if (row.isBase) {
				toast.error(
					"Remove another currency first, or change the base currency from the API.",
				)
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
		},
		[deleteCurrency],
	)

	const columnOptions = useMemo(
		() => ({
			onRateCommit: (row: OrganizationCurrency, raw: string) => {
				void handleRateBlur(row, raw)
			},
			onDelete: (row: OrganizationCurrency) => {
				void handleDelete(row)
			},
			rateUpdatePending: putCurrency.isPending,
			deletePending: deleteCurrency.isPending,
		}),
		[
			handleDelete,
			handleRateBlur,
			putCurrency.isPending,
			deleteCurrency.isPending,
		],
	)

	const columns = useMemo(
		() => createOrganizationCurrencyColumns(columnOptions),
		[columnOptions],
	)

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
		<div className="flex flex-col gap-6">
			<div className="w-full border border-border">
				<DataTableSectionHeader
					title="ORGANIZATION CURRENCIES"
					count={orgCurrencies.length}
					endSlot={
						<Button
							type="button"
							onClick={() => setAddModalOpen(true)}
							disabled={availableCatalog.length === 0 || loading}
						>
							<Plus className="size-4" strokeWidth={2.5} />
							ADD CURRENCY
						</Button>
					}
				/>
				{loading ? (
					<p className="px-6 py-8 font-space-grotesk text-sm text-muted-foreground">
						Loading currencies…
					</p>
				) : orgCurrencies.length === 0 ? (
					<p className="px-6 py-8 font-space-grotesk text-sm text-muted-foreground">
						No currencies yet. Add your base currency to get started.
					</p>
				) : (
					<AccountsDataTable columns={columns} data={orgCurrencies} />
				)}
			</div>

			{lastUpdatedLabel ? (
				<span className="font-ibm-plex-mono text-[10px] tracking-[1px] text-muted-foreground">
					{`LAST UPDATED: ${lastUpdatedLabel.toUpperCase()}`}
				</span>
			) : null}

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
