import {
	type ColumnDef,
	flexRender,
	getCoreRowModel,
	useReactTable,
} from "@tanstack/react-table"
import { Trash2 } from "lucide-react"
import type { Member, MembersTableProps } from "../types"
import { mapApiMember } from "../utils"
import { LoadingRow } from "./loading-row"
import { OwnerBadge } from "./owner-badge"
import { RoleSelect } from "./role-select"

const columns: ColumnDef<Member>[] = [
	{
		accessorKey: "name",
		header: "MEMBER",
		cell: ({ row }) => {
			const member = row.original
			return (
				<div className="flex items-center gap-3">
					<div
						className={`flex h-9 w-9 items-center justify-center ${member.isOwner ? "bg-primary" : "bg-skeleton"}`}
					>
						<span
							className={`font-space-grotesk text-xs font-bold ${member.isOwner ? "text-foreground" : "text-muted-foreground"}`}
						>
							{member.initials}
						</span>
					</div>
					<div className="flex flex-col gap-0.5">
						<span className="font-space-grotesk text-[13px] font-bold tracking-[1px] text-foreground">
							{member.name}
						</span>
						<span className="font-ibm-plex-mono text-[10px] tracking-[1px] text-muted-foreground">
							{member.email}
						</span>
					</div>
				</div>
			)
		},
	},
	{
		accessorKey: "role",
		header: "ROLE",
		cell: ({ row }) => {
			const member = row.original
			return member.isOwner ? (
				<OwnerBadge role={member.role} />
			) : (
				<RoleSelect role={member.role} userId={member.userId} />
			)
		},
	},
	{
		accessorKey: "joined",
		header: "JOINED",
		cell: ({ getValue }) => (
			<span className="font-ibm-plex-mono text-xs tracking-[1px] text-muted-foreground">
				{getValue() as string}
			</span>
		),
	},
	{
		id: "actions",
		header: "ACTIONS",
		cell: ({ row }) => {
			const member = row.original
			return member.isOwner ? (
				<div className="flex w-20 items-center justify-center">
					<span className="font-space-grotesk text-sm font-bold text-muted-foreground">
						—
					</span>
				</div>
			) : (
				<div className="flex w-20 items-center justify-center">
					<button type="button" className="text-destructive hover:opacity-80">
						<Trash2 className="size-4" />
					</button>
				</div>
			)
		},
	},
]

export function MembersTable({
	members: membersData,
	isLoading,
	error,
}: MembersTableProps) {
	const members: Member[] = membersData.map(mapApiMember)

	const table = useReactTable({
		data: members,
		columns,
		getCoreRowModel: getCoreRowModel(),
	})

	return (
		<div className="w-full border border-border">
			<div className="flex h-14 items-center justify-between bg-card px-6 border-b border-border">
				<div className="flex items-center gap-3">
					<div className="h-5 w-1 bg-primary" />
					<span className="font-space-grotesk text-sm font-bold tracking-[1px] text-foreground">
						ACTIVE MEMBERS
					</span>
					<span className="font-ibm-plex-mono text-xs text-muted-foreground">
						{members.length}
					</span>
				</div>
			</div>
			<table className="w-full">
				<thead>
					<tr className="flex h-10 items-center px-6 border-b border-border">
						{table.getHeaderGroups().map((headerGroup) =>
							headerGroup.headers.map((header) => (
								<th
									key={header.id}
									className={
										header.id === "name"
											? "flex-1 text-left font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground"
											: header.id === "actions"
												? "w-20 text-center font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground"
												: header.id === "joined"
													? "w-[140px] text-left font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground"
													: "w-40 text-left font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground"
									}
								>
									{flexRender(
										header.column.columnDef.header,
										header.getContext(),
									)}
								</th>
							)),
						)}
					</tr>
				</thead>
				<tbody>
					{error && (
						<tr className="h-16 border-b border-border">
							<td colSpan={columns.length} className="px-6 text-center">
								<span className="font-ibm-plex-mono text-xs text-destructive">
									{error}
								</span>
							</td>
						</tr>
					)}
					{isLoading ? (
						<>
							<LoadingRow />
							<LoadingRow />
							<LoadingRow />
						</>
					) : (
						table.getRowModel().rows.map((row) => (
							<tr
								key={row.id}
								className="flex h-16 items-center px-6 border-b border-border last:border-b-0"
							>
								{row.getVisibleCells().map((cell) => (
									<td
										key={cell.id}
										className={
											cell.column.id === "name"
												? "flex-1"
												: cell.column.id === "actions"
													? "flex w-20 items-center justify-center"
													: cell.column.id === "joined"
														? "flex w-[140px] items-center"
														: "flex w-40 items-center"
										}
									>
										{flexRender(cell.column.columnDef.cell, cell.getContext())}
									</td>
								))}
							</tr>
						))
					)}
				</tbody>
			</table>
		</div>
	)
}
