import type { ColumnDef } from "@tanstack/react-table"
import type { OrganizationCurrency } from "@workspace/api"
import { Button } from "@workspace/ui/components/button"
import { Input } from "@workspace/ui/components/input"
import {
	TableCellMonoValue,
	TableCellValuePrimary,
} from "@workspace/ui/components/table-cell-values"
import { TableColumnHeader } from "@workspace/ui/components/table-column-header"
import { Trash2 } from "lucide-react"
import { useEffect, useState } from "react"

type ColumnMeta = {
	headerClassName?: string
	cellClassName?: string
}

function formatRate(rate: number | undefined): string {
	if (rate === undefined || Number.isNaN(rate)) return "—"
	return new Intl.NumberFormat(undefined, {
		maximumFractionDigits: 10,
		useGrouping: true,
	}).format(rate)
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

export interface OrganizationCurrencyColumnsOptions {
	onRateCommit: (row: OrganizationCurrency, raw: string) => void
	onDelete: (row: OrganizationCurrency) => void
	rateUpdatePending: boolean
	deletePending: boolean
}

export function createOrganizationCurrencyColumns(
	options: OrganizationCurrencyColumnsOptions,
): ColumnDef<OrganizationCurrency, unknown>[] {
	return [
		{
			accessorKey: "currency",
			meta: {
				headerClassName: "min-w-[200px] px-6",
				cellClassName: "min-w-0 px-6 py-4",
			} satisfies ColumnMeta,
			header: () => <TableColumnHeader>CURRENCY</TableColumnHeader>,
			cell: ({ row }) => {
				const r = row.original
				const name = r.currency?.name ?? r.currencyCode ?? "—"
				return <TableCellValuePrimary>{name}</TableCellValuePrimary>
			},
		},
		{
			accessorKey: "currencyCode",
			meta: {
				headerClassName: "w-[100px] px-6",
				cellClassName: "w-[100px] px-6 py-4",
			} satisfies ColumnMeta,
			header: () => <TableColumnHeader>CODE</TableColumnHeader>,
			cell: ({ row }) => (
				<TableCellMonoValue>
					{row.original.currencyCode ?? "—"}
				</TableCellMonoValue>
			),
		},
		{
			id: "rate",
			meta: {
				headerClassName: "min-w-[160px] px-6",
				cellClassName: "min-w-[160px] px-6 py-4",
			} satisfies ColumnMeta,
			header: () => <TableColumnHeader>RATE (PER 1 BASE)</TableColumnHeader>,
			cell: ({ row }) => {
				const r = row.original
				if (r.isBase) {
					return (
						<TableCellMonoValue className="text-[13px] text-foreground">
							{formatRate(1)}
						</TableCellMonoValue>
					)
				}
				return (
					<RateCell
						initial={r.rate}
						disabled={options.rateUpdatePending}
						onCommit={(raw) => options.onRateCommit(r, raw)}
					/>
				)
			},
		},
		{
			id: "base",
			meta: {
				headerClassName: "w-[100px] px-6 text-center",
				cellClassName: "w-[100px] px-6 py-4 text-center",
			} satisfies ColumnMeta,
			header: () => <TableColumnHeader>BASE</TableColumnHeader>,
			cell: ({ row }) => (
				<TableCellMonoValue className="inline-block">
					{row.original.isBase ? "YES" : "—"}
				</TableCellMonoValue>
			),
		},
		{
			id: "actions",
			meta: {
				headerClassName: "w-[100px] px-6 text-center",
				cellClassName: "w-[100px] px-6 py-4 text-center",
			} satisfies ColumnMeta,
			header: () => <TableColumnHeader>ACTIONS</TableColumnHeader>,
			cell: ({ row }) => {
				const r = row.original
				return (
					<div className="flex items-center justify-center">
						<Button
							type="button"
							variant="destructive"
							size="icon"
							disabled={r.isBase || options.deletePending}
							aria-label={`Remove ${r.currencyCode ?? "currency"}`}
							onClick={() => options.onDelete(r)}
						>
							<Trash2 className="size-4" />
						</Button>
					</div>
				)
			},
		},
	]
}
