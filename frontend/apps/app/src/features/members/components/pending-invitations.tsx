import type { PendingInvitation } from "../types"
import { DataTable } from "./data-table"
import { pendingColumns } from "./pending-columns"

interface PendingInvitationsProps {
	invitations: PendingInvitation[]
}

export function PendingInvitations({ invitations }: PendingInvitationsProps) {
	if (invitations.length === 0) {
		return null
	}

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
							{invitations.length}
						</span>
					</div>
				</div>
			</div>
			<DataTable columns={pendingColumns} data={invitations} />
		</div>
	)
}
