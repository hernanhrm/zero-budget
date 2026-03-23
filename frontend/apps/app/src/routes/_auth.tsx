import { createFileRoute, Outlet, redirect } from "@tanstack/react-router"

import { authClient } from "#/lib/auth-client"

export const Route = createFileRoute("/_auth")({
	beforeLoad: async ({ location }) => {
		const { data } = (await authClient.getSession()) ?? {}
		if (data) {
			const redirectTo =
				(location.search as Record<string, unknown>).redirect as
					| string
					| undefined
			throw redirect({ to: redirectTo ?? "/" })
		}
	},
	component: AuthLayout,
})

function AuthLayout() {
	return <Outlet />
}
