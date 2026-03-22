import { UserPlus } from "lucide-react"
import { MembersTable } from "./components/members-table"
import { PendingInvitations } from "./components/pending-invitations"

export function MembersPage() {
	return (
		<div className="flex h-full flex-col gap-8 overflow-auto p-10">
			<div className="flex w-full items-center justify-between">
				<div className="flex flex-col gap-3">
					<h1 className="font-space-grotesk text-4xl font-bold tracking-[1px] text-[var(--zb-text-primary)]">
						MEMBERS
					</h1>
					<p className="font-ibm-plex-mono text-[13px] tracking-[1px] text-[var(--zb-text-secondary)]">
						INVITE PEOPLE AND MANAGE ROLES FOR YOUR BUDGET
					</p>
				</div>
				<button
					type="button"
					className="flex h-10 items-center gap-2 bg-[var(--zb-accent)] px-4 font-space-grotesk text-xs font-bold tracking-[1px] text-[var(--zb-bg)] hover:opacity-90"
				>
					<UserPlus className="size-3.5" />
					INVITE MEMBER
				</button>
			</div>
			<PendingInvitations />
			<MembersTable />
		</div>
	)
}
