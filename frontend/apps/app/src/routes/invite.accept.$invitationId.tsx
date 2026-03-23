import { useEffect, useState } from "react"
import { createFileRoute, Link, redirect, useNavigate } from "@tanstack/react-router"

import { authClient } from "#/lib/auth-client"

export const Route = createFileRoute("/invite/accept/$invitationId")({
	beforeLoad: async ({ params }) => {
		const { data } = (await authClient.getSession()) ?? {}
		if (!data) {
			throw redirect({
				to: "/sign-in",
				search: { redirect: `/invite/accept/${params.invitationId}` },
			})
		}
	},
	component: AcceptInvitationPage,
})

function AcceptInvitationPage() {
	const { invitationId } = Route.useParams()
	const navigate = useNavigate()
	const [status, setStatus] = useState<"loading" | "error">("loading")
	const [errorMessage, setErrorMessage] = useState("")

	useEffect(() => {
		let cancelled = false

		async function accept() {
			const { error } = await authClient.organization.acceptInvitation({
				invitationId,
			})

			if (cancelled) return

			if (error) {
				setErrorMessage(error.message ?? "Failed to accept invitation")
				setStatus("error")
				return
			}

			navigate({ to: "/" })
		}

		accept()

		return () => {
			cancelled = true
		}
	}, [invitationId, navigate])

	return (
		<div className="flex min-h-screen items-center justify-center bg-card p-6">
			<div className="w-full max-w-sm space-y-6 text-center">
				<div className="flex items-center justify-center gap-2">
					<div className="flex size-7 items-center justify-center bg-primary font-mono text-sm font-bold text-primary-foreground">
						Z
					</div>
					<span className="font-mono text-xs font-bold uppercase tracking-wider">
						Zero Budget
					</span>
				</div>

				{status === "loading" && (
					<div className="space-y-2">
						<h1 className="font-mono text-lg font-bold uppercase tracking-wider">
							Accepting invitation
						</h1>
						<p className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
							Please wait...
						</p>
					</div>
				)}

				{status === "error" && (
					<div className="space-y-4">
						<h1 className="font-mono text-lg font-bold uppercase tracking-wider">
							Invitation failed
						</h1>
						<p className="font-mono text-xs uppercase tracking-wider text-destructive">
							{errorMessage}
						</p>
						<div className="flex flex-col gap-2">
							<Link
								to="/"
								className="font-mono text-xs uppercase tracking-wider text-primary underline underline-offset-4"
							>
								Go to dashboard
							</Link>
							<Link
								to="/sign-in"
								className="font-mono text-xs uppercase tracking-wider text-muted-foreground underline underline-offset-4"
							>
								Sign in with a different account
							</Link>
						</div>
					</div>
				)}
			</div>
		</div>
	)
}
