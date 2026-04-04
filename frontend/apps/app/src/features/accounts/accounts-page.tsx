import { Button } from "@workspace/ui/components/button"
import { ModulePageHeader } from "@workspace/ui/components/module-page-header"
import { Plus } from "lucide-react"
import { Suspense, useState } from "react"
import { Route } from "#/routes/__root"
import { AccountsContent } from "./components/accounts-content"
import { AccountsDataFallback } from "./components/accounts-data-fallback"
import { AddAccountModal } from "./components/add-account-modal"

export function AccountsPage() {
	const [addOpen, setAddOpen] = useState(false)
	const { session } = Route.useRouteContext()
	const sessionData = session?.data as
		| {
				session?: { activeOrganizationId?: string | null }
		  }
		| null
		| undefined
	const organizationId = sessionData?.session?.activeOrganizationId ?? undefined

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
				<AccountsContent />
			</Suspense>

			<AddAccountModal
				open={addOpen}
				onOpenChange={setAddOpen}
				organizationId={organizationId}
			/>
		</div>
	)
}
