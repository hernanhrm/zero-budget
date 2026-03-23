import { useMemo } from "react"
import { authClient } from "#/lib/auth-client"
import type { PendingInvitation } from "../types"
import { DataTable } from "./data-table"
import { createPendingColumns } from "./pending-columns"

interface PendingInvitationsProps {
	invitations: PendingInvitation[]
	onSuccess: () => void
}

export function PendingInvitations({ invitations, onSuccess }: PendingInvitationsProps) {
	const handleResend = async (invitation: PendingInvitation) => {
		const { error } = await authClient.organization.inviteMember({
			email: invitation.email.toLowerCase(),
			role: invitation.role.toLowerCase(),
			resend: true,
		})

		if (error) {
			console.error("Failed to resend invitation:", error)
		}
	}

	const handleCancel = async (invitationId: string) => {
		const { error } = await authClient.organization.cancelInvitation({
			invitationId,
		})

		if (error) {
			console.error("Failed to cancel invitation:", error)
			return
		}

		onSuccess()
	}

	const columns = useMemo(
		() => createPendingColumns({ onResend: handleResend, onCancel: handleCancel }),
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
