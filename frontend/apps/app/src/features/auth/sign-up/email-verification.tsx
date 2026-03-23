import { useState } from "react"
import { authClient } from "#/lib/auth-client"

interface EmailVerificationProps {
	email: string
	redirect?: string
}

export function EmailVerification({
	email,
	redirect,
}: EmailVerificationProps) {
	const [resendStatus, setResendStatus] = useState<
		"idle" | "sending" | "sent" | "error"
	>("idle")

	async function handleResend() {
		setResendStatus("sending")
		try {
			const { error } = await authClient.sendVerificationEmail({
				email,
				callbackURL: redirect
					? `${window.location.origin}${redirect}`
					: window.location.origin,
			})
			if (error) {
				setResendStatus("error")
				return
			}
			setResendStatus("sent")
		} catch {
			setResendStatus("error")
		}
	}

	return (
		<div className="w-full max-w-sm space-y-6">
			<div className="space-y-2">
				<h1 className="font-mono text-lg font-bold uppercase tracking-wider">
					Check your email
				</h1>
				<p className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
					We sent a verification link to{" "}
					<span className="text-foreground">{email}</span>. Click the link to
					verify your account.
				</p>
			</div>
			<div className="space-y-2">
				<p className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
					Didn't receive the email?{" "}
					{resendStatus === "sent" ? (
						<span className="text-foreground">Verification email resent</span>
					) : (
						<button
							type="button"
							onClick={handleResend}
							disabled={resendStatus === "sending"}
							className="text-primary underline underline-offset-4 disabled:opacity-50"
						>
							{resendStatus === "sending"
								? "Sending..."
								: "Resend verification email"}
						</button>
					)}
				</p>
				{resendStatus === "error" && (
					<p className="text-xs text-destructive">
						Failed to resend. Please try again.
					</p>
				)}
			</div>
		</div>
	)
}
