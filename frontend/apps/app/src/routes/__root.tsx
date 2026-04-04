import { lazy, Suspense } from "react"
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router"

import { authClient } from "#/lib/auth-client"
import "@workspace/ui/styles/globals.css"

const RootDevtools = import.meta.env.DEV
	? lazy(() => import("./-__root-devtools"))
	: null

export const Route = createRootRouteWithContext<{
	session: Awaited<ReturnType<typeof authClient.getSession>> | null
}>()({
	beforeLoad: async () => {
		const session = await authClient.getSession()
		return { session }
	},
	component: RootComponent,
})

function RootComponent() {
	return (
		<>
			<Outlet />
			{RootDevtools ? (
				<Suspense fallback={null}>
					<RootDevtools />
				</Suspense>
			) : null}
		</>
	)
}
