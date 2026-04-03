import type { ColumnDef } from "@tanstack/react-table"
import {
	TableCellMonoValue,
	TableCellValueStack,
} from "@workspace/ui/components/table-cell-values"
import { TableColumnHeader } from "@workspace/ui/components/table-column-header"
import { Trash2 } from "lucide-react"
import type { Member } from "../types"
import { OwnerBadge } from "./owner-badge"
import { RoleSelect } from "./role-select"

interface MembersColumnsOptions {
	currentUserId: string
	onRemove: (member: Member) => void
}

export function createMembersColumns(
	options: MembersColumnsOptions,
): ColumnDef<Member>[] {
	return [
		{
			accessorKey: "name",
			header: () => <TableColumnHeader>MEMBER</TableColumnHeader>,
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
						<TableCellValueStack
							primary={member.name}
							secondary={member.email}
						/>
					</div>
				)
			},
		},
		{
			accessorKey: "role",
			header: () => <TableColumnHeader>ROLE</TableColumnHeader>,
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
			header: () => <TableColumnHeader>JOINED</TableColumnHeader>,
			cell: ({ getValue }) => (
				<TableCellMonoValue>{getValue() as string}</TableCellMonoValue>
			),
		},
		{
			id: "actions",
			header: () => <TableColumnHeader>ACTIONS</TableColumnHeader>,
			cell: ({ row }) => {
				const member = row.original
				const isCurrentUser = member.userId === options.currentUserId
				return member.isOwner || isCurrentUser ? (
					<div className="flex w-20 items-center justify-center">
						<span className="font-space-grotesk text-sm font-bold text-muted-foreground">
							—
						</span>
					</div>
				) : (
					<div className="flex w-20 items-center justify-center">
						<button
							type="button"
							className="text-destructive hover:opacity-80"
							onClick={() => options.onRemove(member)}
						>
							<Trash2 className="size-4" />
						</button>
					</div>
				)
			},
		},
	]
}
