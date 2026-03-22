import type { ColumnDef } from "@tanstack/react-table"
import { X } from "lucide-react"
import type { PendingInvitation } from "../types"

export const pendingColumns: ColumnDef<PendingInvitation>[] = [
	{
		accessorKey: "email",
		header: () => (
			<span className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
				EMAIL
			</span>
		),
		cell: ({ row }) => {
			const invitation = row.original
			return (
				<div className="flex items-center gap-3">
					<div className="flex h-8 w-8 items-center justify-center bg-border">
						<span className="font-space-grotesk text-[11px] font-bold text-muted-foreground">
							{invitation.initials}
						</span>
					</div>
					<div className="flex flex-col gap-0.5">
						<span className="font-space-grotesk text-[13px] font-bold tracking-[1px] text-foreground">
							{invitation.email}
						</span>
						<span className="font-ibm-plex-mono text-[10px] tracking-[1px] text-muted-foreground">
							{invitation.invitedAgo}
						</span>
					</div>
				</div>
			)
		},
	},
	{
		accessorKey: "role",
		header: () => (
			<span className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
				ROLE
			</span>
		),
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
		cell: () => (
			<div className="flex items-center gap-3">
				<button
					type="button"
					className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-primary hover:opacity-80"
				>
					RESEND
				</button>
				<button
					type="button"
					className="text-muted-foreground hover:text-destructive"
				>
					<X className="size-3.5" />
				</button>
			</div>
		),
	},
]
