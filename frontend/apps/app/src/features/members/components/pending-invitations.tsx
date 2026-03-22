import type { PendingInvitation } from "../types"
import { DataTable } from "./data-table"
import { pendingColumns } from "./pending-columns"

const pendingInvitations: PendingInvitation[] = [
	{
		email: "SARAH.MILLER@EMAIL.COM",
		initials: "SM",
		role: "EDITOR",
		invitedAgo: "INVITED 2 DAYS AGO",
	},
	{
		email: "MIKE.REYES@EMAIL.COM",
		initials: "MR",
		role: "VIEWER",
		invitedAgo: "INVITED 5 DAYS AGO",
	},
]

export function PendingInvitations() {
	return (
		<div className="w-full border border-border">
			<div className="flex h-14 items-center justify-between bg-card px-6 border-b border-border">
				<div className="flex items-center gap-3">
					<div className="h-5 w-1 bg-primary" />
					<span className="font-space-grotesk text-sm font-bold tracking-[1px] text-foreground">
						PENDING INVITATIONS
					</span>
					<div className="flex h-5 w-6 items-center justify-center bg-primary">
						<span className="font-space-grotesk text-[11px] font-bold text-primary-foreground">
							{pendingInvitations.length}
						</span>
					</div>
				</div>
			</div>
			<DataTable columns={pendingColumns} data={pendingInvitations} />
		</div>
	)
}
