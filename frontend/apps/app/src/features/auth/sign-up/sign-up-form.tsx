import { useForm } from "@tanstack/react-form"
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
import { authClient } from "#/lib/auth-client"
import { signUpSchema } from "./schema"

interface SignUpFormProps {
	onSuccess: (email: string) => void
	serverError: string
	onServerError: (error: string) => void
}

export function SignUpForm({
	onSuccess,
	serverError,
	onServerError,
}: SignUpFormProps) {
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
			onServerError("")
			try {
				const { error } = await authClient.signUp.email({
					name: `${value.firstName} ${value.lastName}`,
					email: value.email,
					password: value.password,
					callbackURL: window.location.origin,
					rememberMe: true,
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
	)
}
