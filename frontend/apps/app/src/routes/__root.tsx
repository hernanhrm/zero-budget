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
		// #region agent log
		void fetch(
			"http://127.0.0.1:7850/ingest/3d0fb17f-0745-41d0-a602-4fb6cb6404fd",
			{
				method: "POST",
				headers: {
					"Content-Type": "application/json",
					"X-Debug-Session-Id": "741baa",
				},
				body: JSON.stringify({
					sessionId: "741baa",
					location: "__root.tsx:beforeLoad",
					message: "client getSession",
					data: {
						hasSessionData:
							session != null &&
							typeof session === "object" &&
							"data" in session &&
							(session as { data?: unknown }).data != null,
						viteIdentityUrl:
							import.meta.env.VITE_IDENTITY_URL ?? "unset",
					},
					timestamp: Date.now(),
					hypothesisId: "H3-vite-baseurl",
					runId: "pre-fix",
				}),
			},
		).catch(() => {})
		// #endregion
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
