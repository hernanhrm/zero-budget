import { DataTableSectionHeader } from "@workspace/ui/components/data-table-section-header"
import { useMemo } from "react"
import { MOCK_ACCOUNT_ROWS } from "../mock-rows"
import { accountsColumns } from "./accounts-columns"
import { AccountsDataTable } from "./accounts-data-table"

export function AccountsTable() {
	const columns = useMemo(() => accountsColumns, [])

	return (
		<div className="w-full border border-border">
			<DataTableSectionHeader
				title="ALL ACCOUNTS"
				count={MOCK_ACCOUNT_ROWS.length}
			/>
			<AccountsDataTable columns={columns} data={MOCK_ACCOUNT_ROWS} />
		</div>
	)
}
