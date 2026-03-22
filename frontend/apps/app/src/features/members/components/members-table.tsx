import type { MembersTableProps } from "../types"
import { mapApiMember } from "../utils"
import { DataTable } from "./data-table"
import { membersColumns } from "./members-columns"

export function MembersTable({
	members: membersData,
	isLoading,
	error,
}: MembersTableProps) {
	const members = membersData.map(mapApiMember)

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
				<DataTable columns={membersColumns} data={members} />
			)}
		</div>
	)
}
