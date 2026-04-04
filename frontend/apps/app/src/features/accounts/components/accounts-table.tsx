import { DataTableSectionHeader } from "@workspace/ui/components/data-table-section-header"
import { useMemo } from "react"
import type { AccountRow } from "../account-row"
import { accountsColumns } from "./accounts-columns"
import { AccountsDataTable } from "./accounts-data-table"

interface AccountsTableProps {
	rows: AccountRow[]
}

export function AccountsTable({ rows }: AccountsTableProps) {
	const columns = useMemo(() => accountsColumns, [])

	return (
		<div className="w-full border border-border">
			<DataTableSectionHeader title="ALL ACCOUNTS" count={rows.length} />
			<AccountsDataTable columns={columns} data={rows} />
		</div>
	)
}
