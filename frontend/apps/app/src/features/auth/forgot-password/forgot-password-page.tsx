import { useState } from "react"
import { BrandingPanel } from "../components/branding-panel"
import { ForgotPasswordForm } from "./forgot-password-form"
import { ResetEmailSent } from "./reset-email-sent"

export function ForgotPasswordPage() {
	const [serverError, setServerError] = useState("")
	const [successEmail, setSuccessEmail] = useState("")

	return (
		<div className="grid min-h-screen grid-cols-1 lg:grid-cols-2">
			<div className="hidden flex-col justify-between bg-neutral-950 text-neutral-50 p-10 lg:flex">
				<BrandingPanel />
			</div>
			<div className="flex items-center justify-center bg-card p-6">
				{successEmail ? (
					<ResetEmailSent email={successEmail} />
				) : (
					<ForgotPasswordForm
						onSuccess={setSuccessEmail}
						serverError={serverError}
						onServerError={setServerError}
					/>
				)}
			</div>
		</div>
	)
}
