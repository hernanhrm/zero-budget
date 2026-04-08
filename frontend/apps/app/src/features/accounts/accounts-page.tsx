import type { Account } from "@workspace/api"
import { Button } from "@workspace/ui/components/button"
import { ModulePageHeader } from "@workspace/ui/components/module-page-header"
import { Plus } from "lucide-react"
import { Suspense, useCallback, useState } from "react"
import { getActiveOrganizationId } from "#/lib/session-org"
import { Route } from "#/routes/__root"
import { AccountsContent } from "./components/accounts-content"
import { AccountsDataFallback } from "./components/accounts-data-fallback"
import { AddAccountModal } from "./components/add-account-modal"
import { DeleteAccountDialog } from "./components/delete-account-dialog"
import { EditAccountModal } from "./components/edit-account-modal"

export function AccountsPage() {
	const [addOpen, setAddOpen] = useState(false)
	const [editAccount, setEditAccount] = useState<Account | null>(null)
	const [deleteAccount, setDeleteAccount] = useState<Account | null>(null)
	const { session } = Route.useRouteContext()
	const organizationId = getActiveOrganizationId(session)

	const handleEditAccount = useCallback((account: Account) => {
		setEditAccount(account)
	}, [])

	const handleDeleteAccount = useCallback((account: Account) => {
		setDeleteAccount(account)
	}, [])

	return (
		<div className="flex h-full flex-col gap-8 overflow-auto p-10">
			<ModulePageHeader
				title="ACCOUNTS"
				description="MANAGE YOUR BANK ACCOUNTS AND BALANCES"
			>
				<Button type="button" onClick={() => setAddOpen(true)}>
					<Plus className="size-4" strokeWidth={2.5} />
					ADD ACCOUNT
				</Button>
			</ModulePageHeader>

			<Suspense fallback={<AccountsDataFallback />}>
				<AccountsContent
					onEditAccount={handleEditAccount}
					onDeleteAccount={handleDeleteAccount}
				/>
			</Suspense>

			<AddAccountModal
				open={addOpen}
				onOpenChange={setAddOpen}
				organizationId={organizationId}
			/>

			<EditAccountModal
				account={editAccount}
				open={editAccount !== null}
				onOpenChange={(open) => {
					if (!open) {
						setEditAccount(null)
					}
				}}
			/>

			<DeleteAccountDialog
				account={deleteAccount}
				open={deleteAccount !== null}
				onOpenChange={(open) => {
					if (!open) {
						setDeleteAccount(null)
					}
				}}
			/>
		</div>
	)
}
