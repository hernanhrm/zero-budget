import { Button } from "@workspace/ui/components/button"
import { ModulePageHeader } from "@workspace/ui/components/module-page-header"
import { Plus } from "lucide-react"
import { useState } from "react"
import { AccountsTable } from "./components/accounts-table"
import { AddAccountModal } from "./components/add-account-modal"
import { MetricCards } from "./components/metric-cards"

export function AccountsPage() {
	const [addOpen, setAddOpen] = useState(false)

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

			<MetricCards />

			<div className="min-w-0 overflow-x-auto">
				<AccountsTable />
			</div>

			<AddAccountModal open={addOpen} onOpenChange={setAddOpen} />
		</div>
	)
}
