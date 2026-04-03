import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from "@workspace/ui/components/alert-dialog"
import { DataTableSectionHeader } from "@workspace/ui/components/data-table-section-header"
import { useMemo, useState } from "react"
import { authClient } from "#/lib/auth-client"
import type { Member, MembersTableProps } from "../types"
import { DataTable } from "./data-table"
import { createMembersColumns } from "./members-columns"

export function MembersTable({
	members,
	isLoading,
	error,
	currentUserId,
	onSuccess,
}: MembersTableProps) {
	const [memberToRemove, setMemberToRemove] = useState<Member | null>(null)

	const handleConfirmRemove = async () => {
		if (!memberToRemove) return

		const { error } = await authClient.organization.removeMember({
			memberIdOrEmail: memberToRemove.id,
		})

		if (error) {
			console.error("Failed to remove member:", error)
			setMemberToRemove(null)
			return
		}

		setMemberToRemove(null)
		onSuccess()
	}

	const columns = useMemo(
		() => createMembersColumns({ currentUserId, onRemove: setMemberToRemove }),
		[currentUserId],
	)

	return (
		<div className="w-full border border-border">
			<DataTableSectionHeader title="ACTIVE MEMBERS" count={members.length} />
			{error ? (
				<div className="flex h-16 items-center justify-center px-6">
					<span className="font-ibm-plex-mono text-xs text-destructive">
						{error}
					</span>
				</div>
			) : isLoading ? (
				<div className="h-16 animate-pulse bg-muted" />
			) : (
				<DataTable columns={columns} data={members} />
			)}
			<AlertDialog
				open={!!memberToRemove}
				onOpenChange={(open) => !open && setMemberToRemove(null)}
			>
				<AlertDialogContent>
					<AlertDialogHeader>
						<AlertDialogTitle className="font-space-grotesk tracking-[1px]">
							REMOVE MEMBER
						</AlertDialogTitle>
						<AlertDialogDescription className="font-ibm-plex-mono text-xs tracking-[1px]">
							ARE YOU SURE YOU WANT TO REMOVE {memberToRemove?.name}? THIS
							ACTION CANNOT BE UNDONE.
						</AlertDialogDescription>
					</AlertDialogHeader>
					<AlertDialogFooter>
						<AlertDialogCancel>CANCEL</AlertDialogCancel>
						<AlertDialogAction
							variant="destructive"
							onClick={handleConfirmRemove}
						>
							REMOVE
						</AlertDialogAction>
					</AlertDialogFooter>
				</AlertDialogContent>
			</AlertDialog>
		</div>
	)
}
