import { DataTableSectionHeader } from "@workspace/ui/components/data-table-section-header"
import { useCallback, useMemo } from "react"
import { authClient } from "#/lib/auth-client"
import type { PendingInvitation } from "../types"
import { DataTable } from "./data-table"
import { createPendingColumns } from "./pending-columns"

interface PendingInvitationsProps {
	invitations: PendingInvitation[]
	onSuccess: () => void
}

export function PendingInvitations({
	invitations,
	onSuccess,
}: PendingInvitationsProps) {
	const handleResend = useCallback(async (invitation: PendingInvitation) => {
		const { error } = await authClient.organization.inviteMember({
			email: invitation.email.toLowerCase(),
			role: invitation.role.toLowerCase() as "member" | "admin" | "owner",
			resend: true,
		})

		if (error) {
			console.error("Failed to resend invitation:", error)
		}
	}, [])

	const handleCancel = useCallback(
		async (invitationId: string) => {
			const { error } = await authClient.organization.cancelInvitation({
				invitationId,
			})

			if (error) {
				console.error("Failed to cancel invitation:", error)
				return
			}

			onSuccess()
		},
		[onSuccess],
	)

	const columns = useMemo(
		() =>
			createPendingColumns({ onResend: handleResend, onCancel: handleCancel }),
		[handleResend, handleCancel],
	)

	if (invitations.length === 0) {
		return null
	}

	return (
		<div className="w-full border border-border">
			<DataTableSectionHeader
				title="PENDING INVITATIONS"
				countSlot={
					<div className="flex h-5 w-6 items-center justify-center bg-primary">
						<span className="font-space-grotesk text-[11px] font-bold text-primary-foreground">
							{invitations.length}
						</span>
					</div>
				}
			/>
			<DataTable columns={columns} data={invitations} />
		</div>
	)
}
