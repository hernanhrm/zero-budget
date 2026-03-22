import type { ColumnDef } from "@tanstack/react-table"
import { Trash2 } from "lucide-react"
import type { Member } from "../types"
import { OwnerBadge } from "./owner-badge"
import { RoleSelect } from "./role-select"

export const membersColumns: ColumnDef<Member>[] = [
	{
		accessorKey: "name",
		header: () => (
			<span className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
				MEMBER
			</span>
		),
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
		header: () => (
			<span className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
				ROLE
			</span>
		),
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
		header: () => (
			<span className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
				JOINED
			</span>
		),
		cell: ({ getValue }) => (
			<span className="font-ibm-plex-mono text-xs tracking-[1px] text-muted-foreground">
				{getValue() as string}
			</span>
		),
	},
	{
		id: "actions",
		header: () => (
			<span className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
				ACTIONS
			</span>
		),
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
