import { Dialog, DialogContent } from "@workspace/ui/components/dialog"
import { DialogPanelHeader } from "@workspace/ui/components/dialog-panel-header"
import { AddAccountForm } from "./add-account-form"

interface AddAccountModalProps {
	open: boolean
	onOpenChange: (open: boolean) => void
	organizationId: string | undefined
}

export function AddAccountModal({
	open,
	onOpenChange,
	organizationId,
}: AddAccountModalProps) {
	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent
				className="max-h-[min(90vh,900px)] gap-0 overflow-y-auto border border-border bg-background p-0 sm:max-w-[520px]"
				showCloseButton={false}
			>
				<DialogPanelHeader title="ADD ACCOUNT" />
				<AddAccountForm
					open={open}
					onComplete={() => onOpenChange(false)}
					organizationId={organizationId}
				/>
			</DialogContent>
		</Dialog>
	)
}
