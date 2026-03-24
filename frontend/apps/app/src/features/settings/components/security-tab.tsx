import { useForm } from "@tanstack/react-form"
import { Button } from "@workspace/ui/components/button"
import {
	Field,
	FieldError,
	FieldGroup,
	FieldLabel,
} from "@workspace/ui/components/field"
import { Input } from "@workspace/ui/components/input"
import { Separator } from "@workspace/ui/components/separator"
import { Switch } from "@workspace/ui/components/switch"
import { useState } from "react"
import { authClient } from "#/lib/auth-client"

export function SecurityTab() {
	const [serverError, setServerError] = useState("")
	const [success, setSuccess] = useState("")
	const [twoFactorEnabled, setTwoFactorEnabled] = useState(false)

	const form = useForm({
		defaultValues: {
			currentPassword: "",
			newPassword: "",
			confirmPassword: "",
		},
		onSubmit: async ({ value }) => {
			setServerError("")
			setSuccess("")

			if (value.newPassword !== value.confirmPassword) {
				setServerError("Passwords do not match")
				return
			}

			try {
				const { error } = await authClient.changePassword({
					currentPassword: value.currentPassword,
					newPassword: value.newPassword,
				})
				if (error) {
					setServerError(error.message ?? "Failed to change password")
					return
				}
				setSuccess("Password changed successfully")
				form.reset()
			} catch (e) {
				setServerError(
					e instanceof Error ? e.message : "Something went wrong",
				)
			}
		},
	})

	return (
		<div className="flex max-w-2xl flex-col gap-8">
			<form
				onSubmit={(e) => {
					e.preventDefault()
					form.handleSubmit()
				}}
			>
				<FieldGroup>
					<form.Field
						name="currentPassword"
						children={(field) => {
							const isInvalid =
								field.state.meta.isTouched &&
								field.state.meta.errors.length > 0
							return (
								<Field data-invalid={isInvalid || undefined}>
									<FieldLabel
										htmlFor={field.name}
										className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground"
									>
										CURRENT PASSWORD
									</FieldLabel>
									<Input
										id={field.name}
										name={field.name}
										type="password"
										value={field.state.value}
										onBlur={field.handleBlur}
										onChange={(e) => field.handleChange(e.target.value)}
										className="h-10 px-4 font-space-grotesk text-sm"
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
						name="newPassword"
						children={(field) => {
							const isInvalid =
								field.state.meta.isTouched &&
								field.state.meta.errors.length > 0
							return (
								<Field data-invalid={isInvalid || undefined}>
									<FieldLabel
										htmlFor={field.name}
										className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground"
									>
										NEW PASSWORD
									</FieldLabel>
									<Input
										id={field.name}
										name={field.name}
										type="password"
										placeholder="Enter new password"
										value={field.state.value}
										onBlur={field.handleBlur}
										onChange={(e) => field.handleChange(e.target.value)}
										className="h-10 px-4 font-space-grotesk text-sm"
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
										className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground"
									>
										CONFIRM PASSWORD
									</FieldLabel>
									<Input
										id={field.name}
										name={field.name}
										type="password"
										placeholder="Confirm new password"
										value={field.state.value}
										onBlur={field.handleBlur}
										onChange={(e) => field.handleChange(e.target.value)}
										className="h-10 px-4 font-space-grotesk text-sm"
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
					{success && (
						<p className="text-xs text-primary">{success}</p>
					)}

					<form.Subscribe
						selector={(state) => state.isSubmitting}
						children={(isSubmitting) => (
							<Button
								type="submit"
								variant="outline"
								disabled={isSubmitting}
								className="w-fit border-primary font-space-grotesk text-xs font-bold tracking-[1px] text-primary hover:bg-primary/10 hover:text-primary"
							>
								{isSubmitting ? "CHANGING..." : "CHANGE PASSWORD"}
							</Button>
						)}
					/>
				</FieldGroup>
			</form>

			<Separator />

			<div className="flex items-center justify-between">
				<div className="flex flex-col gap-1">
					<span className="font-space-grotesk text-[13px] font-bold tracking-[1px] text-foreground">
						2FA AUTHENTICATION
					</span>
					<span className="font-ibm-plex-mono text-[10px] tracking-[1px] text-muted-foreground">
						ADD AN EXTRA LAYER OF SECURITY TO YOUR ACCOUNT
					</span>
				</div>
				<Switch
					checked={twoFactorEnabled}
					onCheckedChange={setTwoFactorEnabled}
				/>
			</div>
		</div>
	)
}
