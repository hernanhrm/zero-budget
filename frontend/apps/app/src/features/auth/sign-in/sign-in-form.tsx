import { useForm } from "@tanstack/react-form"
import { useNavigate } from "@tanstack/react-router"
import { Button } from "@workspace/ui/components/button"
import {
	Field,
	FieldError,
	FieldGroup,
	FieldSeparator,
	FieldLabel,
} from "@workspace/ui/components/field"
import { Input } from "@workspace/ui/components/input"
import { authClient } from "#/lib/auth-client"
import { signInSchema } from "./schema"

interface SignInFormProps {
	serverError: string
	onServerError: (error: string) => void
}

export function SignInForm({ serverError, onServerError }: SignInFormProps) {
	const navigate = useNavigate()

	const form = useForm({
		defaultValues: {
			email: "",
			password: "",
		},
		validators: {
			onSubmit: signInSchema,
		},
		onSubmit: async ({ value }) => {
			onServerError("")
			try {
				const { error } = await authClient.signIn.email({
					email: value.email,
					password: value.password,
					callbackURL: window.location.origin,
					rememberMe: true,
				})
				if (error) {
					onServerError(error.message ?? "Something went wrong")
					return
				}
				navigate({ to: "/" })
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
					Welcome back
				</h1>
				<p className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
					Sign in to your account
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

					<form.Field
						name="password"
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
										Password
									</FieldLabel>
									<Input
										id={field.name}
										name={field.name}
										type="password"
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
								{isSubmitting ? "Signing in..." : "Sign in"}
							</Button>
						)}
					/>

					<FieldSeparator>
						<span className="font-mono text-[10px] uppercase tracking-wider">
							Or
						</span>
					</FieldSeparator>
				</FieldGroup>
			</form>

			<p className="text-center font-mono text-xs uppercase tracking-wider text-muted-foreground">
				Don't have an account?{" "}
				<a
					href="/sign-up"
					className="text-primary underline underline-offset-4"
				>
					Sign up
				</a>
			</p>
		</div>
	)
}
