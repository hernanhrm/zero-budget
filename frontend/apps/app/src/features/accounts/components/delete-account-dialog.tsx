import { useQueryClient } from "@tanstack/react-query"
import type { Account } from "@workspace/api"
import {
	getGetV1AccountsQueryKey,
	useDeleteV1AccountsId,
} from "@workspace/api/hooks/accounts/accounts"
import {
	AlertDialog,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from "@workspace/ui/components/alert-dialog"
import { Button } from "@workspace/ui/components/button"
import { toast } from "@workspace/ui/lib/toast"

function apiProblemDetail(data: unknown): string | undefined {
	if (data !== null && typeof data === "object" && "detail" in data) {
		const d = (data as { detail?: unknown }).detail
		return typeof d === "string" ? d : undefined
	}
	return undefined
}

interface DeleteAccountDialogProps {
	account: Account | null
	open: boolean
	onOpenChange: (open: boolean) => void
}

export function DeleteAccountDialog({
	account,
	open,
	onOpenChange,
}: DeleteAccountDialogProps) {
	const queryClient = useQueryClient()

	const deleteAccount = useDeleteV1AccountsId({
		fetch: { credentials: "include" },
	})

	const handleConfirm = async () => {
		if (!account?.id) {
			onOpenChange(false)
			return
		}

		try {
			const result = await deleteAccount.mutateAsync({ id: account.id })

			if (result.status === 204) {
				await queryClient.invalidateQueries({
					queryKey: getGetV1AccountsQueryKey(),
				})
				toast.success("Account deleted.")
				onOpenChange(false)
				return
			}

			if (result.status === 409) {
				const detail =
					apiProblemDetail(result.data) ??
					"This account has transactions. Disable it instead of deleting."
				toast.error(detail)
				onOpenChange(false)
				return
			}

			const msg =
				apiProblemDetail(result.data) ??
				`Could not delete account (${result.status}).`
			toast.error(msg)
			return
		} catch (e) {
			toast.error(e instanceof Error ? e.message : "Something went wrong.")
		}
	}

	return (
		<AlertDialog open={open} onOpenChange={onOpenChange}>
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle className="font-space-grotesk tracking-[1px]">
						DELETE ACCOUNT
					</AlertDialogTitle>
					<AlertDialogDescription className="font-ibm-plex-mono text-xs tracking-[1px]">
						DELETE {account?.name ?? "THIS ACCOUNT"}? THIS CANNOT BE UNDONE. IF
						THIS ACCOUNT HAS TRANSACTIONS, DELETION WILL BE BLOCKED — USE EDIT
						TO DISABLE IT INSTEAD.
					</AlertDialogDescription>
				</AlertDialogHeader>
				<AlertDialogFooter>
					<AlertDialogCancel>CANCEL</AlertDialogCancel>
					<Button
						type="button"
						variant="destructive"
						disabled={deleteAccount.isPending}
						onClick={() => void handleConfirm()}
					>
						DELETE
					</Button>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	)
}
