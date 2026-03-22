import { UserPlus } from "lucide-react"
import { Route } from "#/routes/_protected/members"
import { MembersTable } from "./components/members-table"
import { PendingInvitations } from "./components/pending-invitations"

export function MembersPage() {
	const loaderData = Route.useLoaderData()

	const members = loaderData?.members?.members || []
	const isLoading = !loaderData
	const error = loaderData?.members?.error
		? String(loaderData.members.error)
		: null

	return (
		<div className="flex h-full flex-col gap-8 overflow-auto p-10">
			<div className="flex w-full items-center justify-between">
				<div className="flex flex-col gap-3">
					<h1 className="font-space-grotesk text-4xl font-bold tracking-[1px] text-foreground">
						MEMBERS
					</h1>
					<p className="font-ibm-plex-mono text-[13px] tracking-[1px] text-muted-foreground">
						INVITE PEOPLE AND MANAGE ROLES FOR YOUR BUDGET
					</p>
				</div>
				<button
					type="button"
					className="flex h-10 items-center gap-2 bg-primary px-4 font-space-grotesk text-xs font-bold tracking-[1px] text-primary-foreground hover:opacity-90"
				>
					<UserPlus className="size-3.5" />
					INVITE MEMBER
				</button>
			</div>
			<PendingInvitations />
			<MembersTable members={members} isLoading={isLoading} error={error} />
		</div>
	)
}

