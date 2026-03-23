import { useMemo } from "react"
import type { PendingInvitation } from "../types"
import { DataTable } from "./data-table"
import { createPendingColumns } from "./pending-columns"

interface PendingInvitationsProps {
	invitations: PendingInvitation[]
}

export function PendingInvitations({ invitations }: PendingInvitationsProps) {
	const handleResend = async (invitationId: string) => {
		try {
			const res = await fetch(
				`${import.meta.env.VITE_IDENTITY_URL}/api/invitations/${invitationId}/resend`,
				{
					method: "POST",
					credentials: "include",
				},
			)

			if (!res.ok) {
				const body = await res.json()
				console.error("Failed to resend invitation:", body.error)
			}
		} catch (err) {
			console.error("Failed to resend invitation:", err)
		}
	}

	const columns = useMemo(
		() => createPendingColumns({ onResend: handleResend }),
		[],
	)

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
			<DataTable columns={columns} data={invitations} />
		</div>
	)
}
