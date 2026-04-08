import type { Account } from "@workspace/api"
import { Dialog, DialogContent } from "@workspace/ui/components/dialog"
import { DialogPanelHeader } from "@workspace/ui/components/dialog-panel-header"
import { EditAccountForm } from "./edit-account-form"

interface EditAccountModalProps {
	account: Account | null
	open: boolean
	onOpenChange: (open: boolean) => void
}

export function EditAccountModal({
	account,
	open,
	onOpenChange,
}: EditAccountModalProps) {
	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent
				className="max-h-[min(90vh,900px)] gap-0 overflow-y-auto border border-border bg-background p-0 sm:max-w-[520px]"
				showCloseButton={false}
			>
				<DialogPanelHeader title="EDIT ACCOUNT" />
				{account ? (
					<EditAccountForm
						open={open}
						account={account}
						onComplete={() => onOpenChange(false)}
					/>
				) : null}
			</DialogContent>
		</Dialog>
	)
}
