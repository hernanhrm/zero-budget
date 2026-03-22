import { createFileRoute, Outlet, redirect } from "@tanstack/react-router"

import { Sidebar } from "#/components/sidebar"
import { authClient } from "#/lib/auth-client"

export const Route = createFileRoute("/_protected")({
	beforeLoad: async () => {
		const { data } = (await authClient.getSession()) ?? {}
		if (!data) {
			throw redirect({ to: "/sign-in" })
		}
	},
	component: ProtectedLayout,
})

function ProtectedLayout() {
	return (
		<div className="flex h-screen bg-[var(--zb-bg)]">
			<Sidebar />
			<main className="flex-1 overflow-hidden">
				<Outlet />
			</main>
		</div>
	)
}
