import { QueryClientProvider } from "@tanstack/react-query"
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router"
import { Toaster } from "@workspace/ui/components/toaster"
import { ThemeProvider } from "next-themes"
import { lazy, Suspense } from "react"
import { authClient } from "#/lib/auth-client"
import { queryClient } from "#/lib/query-client"
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
		<ThemeProvider
			attribute="class"
			defaultTheme="dark"
			enableSystem
			disableTransitionOnChange
		>
			<QueryClientProvider client={queryClient}>
				<Outlet />
				<Toaster />
				{RootDevtools ? (
					<Suspense fallback={null}>
						<RootDevtools />
					</Suspense>
				) : null}
			</QueryClientProvider>
		</ThemeProvider>
	)
}
