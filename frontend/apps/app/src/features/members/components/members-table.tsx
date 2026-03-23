import { useMemo, useState } from "react"
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
			<div className="flex h-14 items-center justify-between bg-card px-6 border-b border-border">
				<div className="flex items-center gap-3">
					<div className="h-5 w-1 bg-primary" />
					<span className="font-space-grotesk text-sm font-bold tracking-[1px] text-foreground">
						ACTIVE MEMBERS
					</span>
					<span className="font-ibm-plex-mono text-xs text-muted-foreground">
						{members.length}
					</span>
				</div>
			</div>
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
			<AlertDialog open={!!memberToRemove} onOpenChange={(open) => !open && setMemberToRemove(null)}>
				<AlertDialogContent>
					<AlertDialogHeader>
						<AlertDialogTitle className="font-space-grotesk tracking-[1px]">
							REMOVE MEMBER
						</AlertDialogTitle>
						<AlertDialogDescription className="font-ibm-plex-mono text-xs tracking-[1px]">
							ARE YOU SURE YOU WANT TO REMOVE {memberToRemove?.name}? THIS ACTION CANNOT BE UNDONE.
						</AlertDialogDescription>
					</AlertDialogHeader>
					<AlertDialogFooter>
						<AlertDialogCancel className="font-space-grotesk text-xs font-bold tracking-[1px]">
							CANCEL
						</AlertDialogCancel>
						<AlertDialogAction
							variant="destructive"
							className="font-space-grotesk text-xs font-bold tracking-[1px]"
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
