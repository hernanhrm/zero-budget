import { createFileRoute } from "@tanstack/react-router"
import { authClient } from "#/lib/auth-client"
import { MembersPage } from "#/features/members/members-page"

export const Route = createFileRoute("/_protected/members")({
	loader: async () => {
		const { data } = await authClient.organization.listMembers({
			query: { limit: 50 },
		})
		return { members: data }
	},
	component: MembersPage,
})
