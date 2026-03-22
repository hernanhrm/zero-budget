import { createFileRoute } from "@tanstack/react-router"
import { authClient } from "#/lib/auth-client"
import { MembersPage } from "#/features/members/members-page"

export const Route = createFileRoute("/_protected/members")({
	loader: async () => {
		const [membersRes, invitationsRes] = await Promise.all([
			authClient.organization.listMembers({
				query: { limit: 50 },
			}),
			authClient.organization.listInvitations(),
		])
		return { members: membersRes.data, invitations: invitationsRes.data }
	},
	component: MembersPage,
})
