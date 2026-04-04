import type { ColumnDef } from "@tanstack/react-table"
import {
	TableCellAmountValue,
	TableCellMonoValue,
	TableCellStatusBadge,
	TableCellValueStack,
} from "@workspace/ui/components/table-cell-values"
import { TableColumnHeader } from "@workspace/ui/components/table-column-header"
import type { AccountRow } from "../account-row"

type AccountColumnMeta = {
	headerClassName?: string
	cellClassName?: string
}

export const accountsColumns: ColumnDef<AccountRow, unknown>[] = [
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
			<TableCellMonoValue>{row.original.type}</TableCellMonoValue>
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
		cell: ({ row }) => <TableCellStatusBadge active={row.original.isActive} />,
	},
]
