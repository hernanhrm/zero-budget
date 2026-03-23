import { useState } from "react"
import { BrandingPanel } from "../components/branding-panel"
import { SignInForm } from "./sign-in-form"

interface SignInPageProps {
	redirect?: string
}

export function SignInPage({ redirect }: SignInPageProps) {
	const [serverError, setServerError] = useState("")

	return (
		<div className="grid min-h-screen grid-cols-1 lg:grid-cols-2">
			<div className="hidden flex-col justify-between bg-neutral-950 text-neutral-50 p-10 lg:flex">
				<BrandingPanel />
			</div>
			<div className="flex items-center justify-center bg-card p-6">
				<SignInForm
					redirect={redirect}
					serverError={serverError}
					onServerError={setServerError}
				/>
			</div>
		</div>
	)
}
