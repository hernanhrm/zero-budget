import { useForm } from "@tanstack/react-form"
import { createFileRoute } from "@tanstack/react-router"
import { Button } from "@workspace/ui/components/button"
import { Checkbox } from "@workspace/ui/components/checkbox"
import {
	Field,
	FieldError,
	FieldGroup,
	FieldLabel,
	FieldSeparator,
} from "@workspace/ui/components/field"
import { Input } from "@workspace/ui/components/input"
import { useState } from "react"
import { z } from "zod"
import { authClient } from "#/lib/auth-client"

export const Route = createFileRoute("/sign-up")({ component: SignUp })

const signUpSchema = z
	.object({
		firstName: z.string().min(1, "First name is required"),
		lastName: z.string().min(1, "Last name is required"),
		email: z
			.string()
			.min(1, "Email is required")
			.email("Invalid email address"),
		password: z.string().min(8, "Password must be at least 8 characters"),
		confirmPassword: z.string().min(1, "Please confirm your password"),
		terms: z.literal(true, {
			errorMap: () => ({ message: "You must accept the terms" }),
		}),
	})
	.refine((data) => data.password === data.confirmPassword, {
		message: "Passwords don't match",
		path: ["confirmPassword"],
	})

function SignUp() {
	const [serverError, setServerError] = useState("")
	const [successEmail, setSuccessEmail] = useState("")
	const [resendStatus, setResendStatus] = useState<
		"idle" | "sending" | "sent" | "error"
	>("idle")

	const form = useForm({
		defaultValues: {
			firstName: "",
			lastName: "",
			email: "",
			password: "",
			confirmPassword: "",
			terms: false as boolean,
		},
		validators: {
			onSubmit: signUpSchema,
		},
		onSubmit: async ({ value }) => {
			setServerError("")
			try {
				const { error } = await authClient.signUp.email({
					name: `${value.firstName} ${value.lastName}`,
					email: value.email,
					password: value.password,
					callbackURL: window.location.origin,
					rememberMe: true,
				})
				if (error) {
					setServerError(error.message ?? "Something went wrong")
					return
				}
				setSuccessEmail(value.email)
			} catch (e) {
				setServerError(
					e instanceof Error ? e.message : "Something went wrong",
				)
			}
		},
	})

	async function handleResend() {
		setResendStatus("sending")
		try {
			const { error } = await authClient.sendVerificationEmail({
				email: successEmail,
				callbackURL: window.location.origin,
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

	if (successEmail) {
		return (
			<div className="grid min-h-screen grid-cols-1 lg:grid-cols-2">
				<div className="hidden flex-col justify-between bg-neutral-950 text-neutral-50 p-10 lg:flex">
					<BrandingPanel />
				</div>
				<div className="flex items-center justify-center bg-card p-6">
					<div className="w-full max-w-sm space-y-6">
						<div className="space-y-2">
							<h1 className="font-mono text-lg font-bold uppercase tracking-wider">
								Check your email
							</h1>
							<p className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
								We sent a verification link to{" "}
								<span className="text-foreground">{successEmail}</span>. Click
								the link to verify your account.
							</p>
						</div>
						<div className="space-y-2">
							<p className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
								Didn't receive the email?{" "}
								{resendStatus === "sent" ? (
									<span className="text-foreground">
										Verification email resent
									</span>
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
							Create account
						</h1>
						<p className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
							Enter your details to get started
						</p>
					</div>

					<form
						onSubmit={(e) => {
							e.preventDefault()
							form.handleSubmit()
						}}
					>
						<FieldGroup>
							<div className="grid grid-cols-2 gap-4">
								<form.Field
									name="firstName"
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
													First Name
												</FieldLabel>
												<Input
													id={field.name}
													name={field.name}
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
									name="lastName"
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
													Last Name
												</FieldLabel>
												<Input
													id={field.name}
													name={field.name}
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
							</div>

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

							<form.Field
								name="confirmPassword"
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
												Confirm Password
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

							<form.Field
								name="terms"
								children={(field) => {
									const isInvalid =
										field.state.meta.isTouched &&
										field.state.meta.errors.length > 0
									return (
										<Field
											orientation="horizontal"
											data-invalid={isInvalid || undefined}
										>
											<Checkbox
												id={field.name}
												checked={field.state.value}
												onCheckedChange={(checked) =>
													field.handleChange(checked === true)
												}
												onBlur={field.handleBlur}
												aria-invalid={isInvalid || undefined}
											/>
											<div className="space-y-1">
												<label
													htmlFor={field.name}
													className="font-mono text-xs uppercase tracking-wider leading-none cursor-pointer select-none"
												>
													I agree to the Terms of Service and Privacy Policy
												</label>
												{isInvalid && (
													<FieldError errors={field.state.meta.errors} />
												)}
											</div>
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
										{isSubmitting ? "Creating account..." : "Create account"}
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
						Already have an account?{" "}
						<a
							href="/sign-in"
							className="text-primary underline underline-offset-4"
						>
							Sign in
						</a>
					</p>
				</div>
			</div>
		</div>
	)
}

function BrandingPanel() {
	return (
		<>
			<div className="flex items-center gap-2">
				<div className="flex size-8 items-center justify-center bg-primary font-mono text-sm font-bold text-primary-foreground">
					Z
				</div>
				<span className="font-mono text-xs font-bold uppercase tracking-wider">
					Zero Budget
				</span>
			</div>

			<div className="space-y-6">
				<h2 className="font-mono text-4xl font-bold uppercase leading-tight tracking-wider">
					Take control
					<br />
					of your
					<br />
					finances
				</h2>
				<div className="space-y-3">
					<p className="font-mono text-xs uppercase tracking-wider text-neutral-400">
						Track every currency. Manage every account.
						<br />
						Master every category.
					</p>
					<div className="h-0.5 w-16 bg-primary" />
				</div>
			</div>

			<div className="flex items-center justify-between font-mono text-[10px] uppercase tracking-wider text-neutral-400">
				<div className="flex items-center gap-2">
					<span className="inline-block size-1.5 bg-primary" />
					<span>System online</span>
				</div>
				<span>v0.1.0</span>
			</div>
		</>
	)
}
