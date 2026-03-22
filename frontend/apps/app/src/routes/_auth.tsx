import { createFileRoute, Outlet, redirect } from "@tanstack/react-router"

import { authClient } from "#/lib/auth-client"

export const Route = createFileRoute("/_auth")({
	beforeLoad: async () => {
		const { data } = (await authClient.getSession()) ?? {}
		if (data) {
			throw redirect({ to: "/" })
		}
	},
	component: AuthLayout,
})

function AuthLayout() {
	return <Outlet />
}
