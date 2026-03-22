import { createFileRoute, Outlet, redirect } from "@tanstack/react-router"

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
		<div className="min-h-screen bg-background">
			<nav className="border-b p-4">
				<span className="font-semibold">Zero Budget</span>
			</nav>
			<main className="p-6">
				<Outlet />
			</main>
		</div>
	)
}
