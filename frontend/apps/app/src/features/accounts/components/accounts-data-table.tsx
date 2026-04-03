import {
	type ColumnDef,
	flexRender,
	getCoreRowModel,
	useReactTable,
} from "@tanstack/react-table"
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@workspace/ui/components/table"
import { cn } from "@workspace/ui/lib/utils"

type ColumnMeta = {
	headerClassName?: string
	cellClassName?: string
}

interface AccountsDataTableProps<TData, TValue> {
	columns: ColumnDef<TData, TValue>[]
	data: TData[]
}

export function AccountsDataTable<TData, TValue>({
	columns,
	data,
}: AccountsDataTableProps<TData, TValue>) {
	const table = useReactTable({
		data,
		columns,
		getCoreRowModel: getCoreRowModel(),
	})

	return (
		<div className="w-full">
			<Table className="table-fixed">
				<TableHeader>
					{table.getHeaderGroups().map((headerGroup) => (
						<TableRow
							key={headerGroup.id}
							className="border-border hover:bg-transparent"
						>
							{headerGroup.headers.map((header) => {
								const meta = header.column.columnDef.meta as
									| ColumnMeta
									| undefined
								return (
									<TableHead
										key={header.id}
										className={cn(
											"h-10 border-b border-border align-middle",
											meta?.headerClassName,
										)}
									>
										{header.isPlaceholder
											? null
											: flexRender(
													header.column.columnDef.header,
													header.getContext(),
												)}
									</TableHead>
								)
							})}
						</TableRow>
					))}
				</TableHeader>
				<TableBody>
					{table.getRowModel().rows?.length ? (
						table.getRowModel().rows.map((row) => (
							<TableRow key={row.id} className="border-border">
								{row.getVisibleCells().map((cell) => {
									const meta = cell.column.columnDef.meta as
										| ColumnMeta
										| undefined
									return (
										<TableCell
											key={cell.id}
											className={cn(meta?.cellClassName)}
										>
											{flexRender(
												cell.column.columnDef.cell,
												cell.getContext(),
											)}
										</TableCell>
									)
								})}
							</TableRow>
						))
					) : (
						<TableRow>
							<TableCell
								colSpan={columns.length}
								className="h-24 px-6 text-center"
							>
								No results.
							</TableCell>
						</TableRow>
					)}
				</TableBody>
			</Table>
		</div>
	)
}
