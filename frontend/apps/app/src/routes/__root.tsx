import { TanStackDevtools } from "@tanstack/react-devtools"
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router"
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools"

import { authClient } from "#/lib/auth-client"
import "@workspace/ui/styles/globals.css"

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
			<TanStackDevtools
				config={{
					position: "bottom-right",
				}}
				plugins={[
					{
						name: "TanStack Router",
						render: <TanStackRouterDevtoolsPanel />,
					},
				]}
			/>
		</>
	)
}
