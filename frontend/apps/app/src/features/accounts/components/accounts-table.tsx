import { DataTableSectionHeader } from "@workspace/ui/components/data-table-section-header"
import { useMemo } from "react"
import type { AccountRow } from "../account-row"
import {
	type AccountsColumnsOptions,
	createAccountsColumns,
} from "./accounts-columns"
import { AccountsDataTable } from "./accounts-data-table"

interface AccountsTableProps {
	rows: AccountRow[]
	columnOptions: AccountsColumnsOptions
}

export function AccountsTable({ rows, columnOptions }: AccountsTableProps) {
	const columns = useMemo(
		() => createAccountsColumns(columnOptions),
		[columnOptions],
	)

	return (
		<div className="w-full border border-border">
			<DataTableSectionHeader title="ALL ACCOUNTS" count={rows.length} />
			<AccountsDataTable columns={columns} data={rows} />
		</div>
	)
}
