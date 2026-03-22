import { useState } from "react"
import { Route } from "#/routes/_auth/reset-password"
import { BrandingPanel } from "../components/branding-panel"
import { ResetPasswordForm } from "./reset-password-form"

export function ResetPasswordPage() {
	const [serverError, setServerError] = useState("")
	const { token, error } = Route.useSearch()

	if (error) {
		return (
			<div className="grid min-h-screen grid-cols-1 lg:grid-cols-2">
				<div className="hidden flex-col justify-between bg-neutral-950 text-neutral-50 p-10 lg:flex">
					<BrandingPanel />
				</div>
				<div className="flex items-center justify-center bg-card p-6">
					<div className="w-full max-w-sm space-y-6">
						<div className="space-y-2">
							<h1 className="font-mono text-lg font-bold uppercase tracking-wider">
								Invalid reset link
							</h1>
							<p className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
								This password reset link is invalid or has expired.
							</p>
						</div>
						<p className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
							<a
								href="/forgot-password"
								className="text-primary underline underline-offset-4"
							>
								Request a new link
							</a>
						</p>
					</div>
				</div>
			</div>
		)
	}

	if (!token) {
		return (
			<div className="grid min-h-screen grid-cols-1 lg:grid-cols-2">
				<div className="hidden flex-col justify-between bg-neutral-950 text-neutral-50 p-10 lg:flex">
					<BrandingPanel />
				</div>
				<div className="flex items-center justify-center bg-card p-6">
					<div className="w-full max-w-sm space-y-6">
						<div className="space-y-2">
							<h1 className="font-mono text-lg font-bold uppercase tracking-wider">
								Missing reset token
							</h1>
							<p className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
								No reset token was provided. Please use the link from
								your email.
							</p>
						</div>
						<p className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
							<a
								href="/forgot-password"
								className="text-primary underline underline-offset-4"
							>
								Request a new link
							</a>
						</p>
					</div>
				</div>
			</div>
		)
	}

	return (
		<div className="grid min-h-screen grid-cols-1 lg:grid-cols-2">
			<div className="hidden flex-col justify-between bg-neutral-950 text-neutral-50 p-10 lg:flex">
				<BrandingPanel />
			</div>
			<div className="flex items-center justify-center bg-card p-6">
				<ResetPasswordForm
					token={token}
					serverError={serverError}
					onServerError={setServerError}
				/>
			</div>
		</div>
	)
}
