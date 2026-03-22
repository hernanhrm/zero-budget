interface ResetEmailSentProps {
	email: string
}

export function ResetEmailSent({ email }: ResetEmailSentProps) {
	return (
		<div className="w-full max-w-sm space-y-6">
			<div className="space-y-2">
				<h1 className="font-mono text-lg font-bold uppercase tracking-wider">
					Check your email
				</h1>
				<p className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
					We sent a password reset link to{" "}
					<span className="text-foreground">{email}</span>. Click the
					link to reset your password.
				</p>
			</div>
			<p className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
				<a
					href="/sign-in"
					className="text-primary underline underline-offset-4"
				>
					Back to sign in
				</a>
			</p>
		</div>
	)
}
