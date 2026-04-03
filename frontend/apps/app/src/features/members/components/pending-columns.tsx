import type { ColumnDef } from "@tanstack/react-table"
import { TableCellValueStack } from "@workspace/ui/components/table-cell-values"
import { TableColumnHeader } from "@workspace/ui/components/table-column-header"
import { X } from "lucide-react"
import type { PendingInvitation } from "../types"

interface PendingColumnsOptions {
	onResend: (invitation: PendingInvitation) => void
	onCancel: (invitationId: string) => void
}

export function createPendingColumns(
	options: PendingColumnsOptions,
): ColumnDef<PendingInvitation>[] {
	return [
		{
			accessorKey: "email",
			header: () => <TableColumnHeader>EMAIL</TableColumnHeader>,
			cell: ({ row }) => {
				const invitation = row.original
				return (
					<div className="flex items-center gap-3">
						<div className="flex h-8 w-8 items-center justify-center bg-border">
							<span className="font-space-grotesk text-[11px] font-bold text-muted-foreground">
								{invitation.initials}
							</span>
						</div>
						<TableCellValueStack
							primary={invitation.email}
							secondary={invitation.invitedAgo}
						/>
					</div>
				)
			},
		},
		{
			accessorKey: "role",
			header: () => <TableColumnHeader>ROLE</TableColumnHeader>,
			cell: ({ getValue }) => (
				<div className="inline-flex h-5 items-center border border-border px-2">
					<span className="font-space-grotesk text-[10px] font-bold tracking-[1px] text-muted-foreground">
						{getValue() as string}
					</span>
				</div>
			),
		},
		{
			id: "actions",
			header: "",
			cell: ({ row }) => (
				<div className="flex items-center gap-3">
					<button
						type="button"
						className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-primary hover:opacity-80"
						onClick={() => options.onResend(row.original)}
					>
						RESEND
					</button>
					<button
						type="button"
						className="text-muted-foreground hover:text-destructive"
						onClick={() => options.onCancel(row.original.id)}
					>
						<X className="size-3.5" />
					</button>
				</div>
			),
		},
	]
}
