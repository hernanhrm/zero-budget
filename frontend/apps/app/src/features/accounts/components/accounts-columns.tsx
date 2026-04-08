import type { ColumnDef } from "@tanstack/react-table"
import { Button } from "@workspace/ui/components/button"
import {
	TableCellAmountValue,
	TableCellMonoValue,
	TableCellStatusBadge,
	TableCellValueStack,
} from "@workspace/ui/components/table-cell-values"
import { TableColumnHeader } from "@workspace/ui/components/table-column-header"
import { Pencil, Trash2 } from "lucide-react"
import type { AccountRow } from "../account-row"
import { formatAccountTypeLabel } from "../account-type-options"

type AccountColumnMeta = {
	headerClassName?: string
	cellClassName?: string
}

export interface AccountsColumnsOptions {
	onEdit: (row: AccountRow) => void
	onDelete: (row: AccountRow) => void
}

export function createAccountsColumns(
	options: AccountsColumnsOptions,
): ColumnDef<AccountRow, unknown>[] {
	return [
		{
			accessorKey: "name",
			meta: {
				headerClassName: "min-w-[200px] px-6",
				cellClassName: "min-w-0 px-6 py-4",
			} satisfies AccountColumnMeta,
			header: () => <TableColumnHeader>ACCOUNT NAME</TableColumnHeader>,
			cell: ({ row }) => {
				const r = row.original
				return <TableCellValueStack primary={r.name} secondary={r.mask} />
			},
		},
		{
			accessorKey: "type",
			meta: {
				headerClassName: "w-[140px] px-6",
				cellClassName: "w-[140px] px-6 py-4",
			} satisfies AccountColumnMeta,
			header: () => <TableColumnHeader>TYPE</TableColumnHeader>,
			cell: ({ row }) => (
				<TableCellMonoValue>
					{formatAccountTypeLabel(row.original.type)}
				</TableCellMonoValue>
			),
		},
		{
			accessorKey: "institution",
			meta: {
				headerClassName: "w-[180px] px-6",
				cellClassName: "w-[180px] px-6 py-4",
			} satisfies AccountColumnMeta,
			header: () => <TableColumnHeader>INSTITUTION</TableColumnHeader>,
			cell: ({ row }) => (
				<TableCellMonoValue>{row.original.institution}</TableCellMonoValue>
			),
		},
		{
			accessorKey: "balance",
			meta: {
				headerClassName: "w-[140px] px-6 text-right",
				cellClassName: "w-[140px] px-6 py-4 text-right",
			} satisfies AccountColumnMeta,
			header: () => <TableColumnHeader>BALANCE</TableColumnHeader>,
			cell: ({ row }) => {
				const r = row.original
				return (
					<TableCellAmountValue
						emphasis={r.balanceIsPrimary ? "primary" : "default"}
					>
						{r.balance}
					</TableCellAmountValue>
				)
			},
		},
		{
			id: "status",
			meta: {
				headerClassName: "w-[100px] px-6 text-center",
				cellClassName: "w-[100px] px-6 py-4 text-center",
			} satisfies AccountColumnMeta,
			header: () => <TableColumnHeader>STATUS</TableColumnHeader>,
			cell: ({ row }) => (
				<TableCellStatusBadge active={row.original.isActive} />
			),
		},
		{
			id: "actions",
			meta: {
				headerClassName: "w-[100px] px-6 text-center",
				cellClassName: "w-[100px] px-6 py-4 text-center",
			} satisfies AccountColumnMeta,
			header: () => <TableColumnHeader>ACTIONS</TableColumnHeader>,
			cell: ({ row }) => {
				const r = row.original
				return (
					<div className="flex items-center justify-center gap-1">
						<Button
							type="button"
							variant="ghost"
							size="icon"
							aria-label={`Edit ${r.name}`}
							onClick={() => options.onEdit(r)}
						>
							<Pencil className="size-4" />
						</Button>
						<Button
							type="button"
							variant="destructive"
							size="icon"
							aria-label={`Delete ${r.name}`}
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
