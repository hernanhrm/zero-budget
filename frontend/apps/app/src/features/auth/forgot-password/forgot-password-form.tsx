import { useForm } from "@tanstack/react-form"
import { Button } from "@workspace/ui/components/button"
import {
	Field,
	FieldError,
	FieldGroup,
	FieldLabel,
} from "@workspace/ui/components/field"
import { Input } from "@workspace/ui/components/input"
import { authClient } from "#/lib/auth-client"
import { forgotPasswordSchema } from "./schema"

interface ForgotPasswordFormProps {
	serverError: string
	onServerError: (error: string) => void
	onSuccess: (email: string) => void
}

export function ForgotPasswordForm({
	serverError,
	onServerError,
	onSuccess,
}: ForgotPasswordFormProps) {
	const form = useForm({
		defaultValues: {
			email: "",
		},
		validators: {
			onSubmit: forgotPasswordSchema,
		},
		onSubmit: async ({ value }) => {
			onServerError("")
			try {
				const { error } = await authClient.requestPasswordReset({
					email: value.email,
					redirectTo: `${window.location.origin}/reset-password`,
				})
				if (error) {
					onServerError(error.message ?? "Something went wrong")
					return
				}
				onSuccess(value.email)
			} catch (e) {
				onServerError(
					e instanceof Error ? e.message : "Something went wrong",
				)
			}
		},
	})

	return (
		<div className="w-full max-w-sm space-y-6">
			<div className="space-y-1">
				<div className="flex items-center gap-2 lg:hidden mb-4">
					<div className="flex size-7 items-center justify-center bg-primary font-mono text-sm font-bold text-primary-foreground">
						Z
					</div>
					<span className="font-mono text-xs font-bold uppercase tracking-wider">
						Zero Budget
					</span>
				</div>
				<h1 className="font-mono text-lg font-bold uppercase tracking-wider">
					Forgot password?
				</h1>
				<p className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
					Enter your email to receive a reset link
				</p>
			</div>

			<form
				onSubmit={(e) => {
					e.preventDefault()
					form.handleSubmit()
				}}
			>
				<FieldGroup>
					<form.Field
						name="email"
						children={(field) => {
							const isInvalid =
								field.state.meta.isTouched &&
								field.state.meta.errors.length > 0
							return (
								<Field data-invalid={isInvalid || undefined}>
									<FieldLabel
										htmlFor={field.name}
										className="font-mono uppercase tracking-wider"
									>
										Email
									</FieldLabel>
									<Input
										id={field.name}
										name={field.name}
										type="email"
										value={field.state.value}
										onBlur={field.handleBlur}
										onChange={(e) => field.handleChange(e.target.value)}
										aria-invalid={isInvalid || undefined}
									/>
									{isInvalid && (
										<FieldError errors={field.state.meta.errors} />
									)}
								</Field>
							)
						}}
					/>

					{serverError && (
						<p className="text-xs text-destructive">{serverError}</p>
					)}

					<form.Subscribe
						selector={(state) => state.isSubmitting}
						children={(isSubmitting) => (
							<Button
								type="submit"
								disabled={isSubmitting}
								className="w-full font-mono uppercase tracking-wider"
							>
								{isSubmitting ? "Sending..." : "Send reset link"}
							</Button>
						)}
					/>
				</FieldGroup>
			</form>

			<p className="text-center font-mono text-xs uppercase tracking-wider text-muted-foreground">
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
