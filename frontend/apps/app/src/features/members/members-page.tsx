import { useRouter } from "@tanstack/react-router"
import { Button } from "@workspace/ui/components/button"
import { ModulePageHeader } from "@workspace/ui/components/module-page-header"
import { UserPlus } from "lucide-react"
import { useState } from "react"
import { Route } from "#/routes/_protected/members"
import { InviteMemberModal } from "./components/invite-member-modal"
import { MembersTable } from "./components/members-table"
import { PendingInvitations } from "./components/pending-invitations"
import type { ApiInvitation } from "./types"
import { mapApiInvitation, mapApiMember } from "./utils"

export function MembersPage() {
	const loaderData = Route.useLoaderData()
	const { session } = Route.useRouteContext()
	const router = useRouter()
	const [inviteOpen, setInviteOpen] = useState(false)

	const currentUserId = session?.data?.user?.id ?? ""
	const membersResult = loaderData?.members
	const apiMembers = membersResult?.data?.members ?? []
	const members = apiMembers.map(mapApiMember)
	const isLoading = !loaderData
	const error = membersResult?.error
		? String(membersResult.error.message ?? membersResult.error)
		: null

	const apiInvitations = (loaderData?.invitations ?? []) as ApiInvitation[]
	const pendingInvitations = apiInvitations
		.filter((inv) => inv.status === "pending")
		.map(mapApiInvitation)

	return (
		<div className="flex h-full flex-col gap-8 overflow-auto p-10">
			<ModulePageHeader
				title="MEMBERS"
				description="INVITE PEOPLE AND MANAGE ROLES FOR YOUR BUDGET"
			>
				<Button
					type="button"
					onClick={() => setInviteOpen(true)}
					className="gap-2"
				>
					<UserPlus className="size-3.5" />
					INVITE MEMBER
				</Button>
			</ModulePageHeader>
			<PendingInvitations
				invitations={pendingInvitations}
				onSuccess={() => router.invalidate()}
			/>
			<MembersTable
				members={members}
				isLoading={isLoading}
				error={error}
				currentUserId={currentUserId}
				onSuccess={() => router.invalidate()}
			/>
			<InviteMemberModal
				open={inviteOpen}
				onOpenChange={setInviteOpen}
				onSuccess={() => router.invalidate()}
			/>
		</div>
	)
}
