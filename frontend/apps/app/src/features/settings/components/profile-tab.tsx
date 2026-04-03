import { useForm } from "@tanstack/react-form"
import { useRouter } from "@tanstack/react-router"
import { Button } from "@workspace/ui/components/button"
import { Field, FieldError, FieldGroup } from "@workspace/ui/components/field"
import { FormFieldLabel } from "@workspace/ui/components/form-field-label"
import { Input } from "@workspace/ui/components/input"
import { useState } from "react"
import { authClient } from "#/lib/auth-client"
import { Route } from "#/routes/_protected/settings"
import { profileSchema } from "../schema"

const inputClassName =
	"h-10 rounded-sm border-border bg-transparent px-4 font-space-grotesk text-sm text-foreground"

export function ProfileTab() {
	const { session } = Route.useRouteContext()
	const router = useRouter()
	const [serverError, setServerError] = useState("")

	const user = session?.data?.user
	const nameParts = (user?.name ?? "").split(" ")
	const firstName = nameParts[0] ?? ""
	const lastName = nameParts.slice(1).join(" ")

	const defaultValues = {
		firstName,
		lastName,
		timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
	}

	const hasChanges = (
		current: typeof defaultValues,
		initial: typeof defaultValues,
	) => {
		return (
			current.firstName !== initial.firstName ||
			current.lastName !== initial.lastName ||
			current.timezone !== initial.timezone
		)
	}

	const form = useForm({
		defaultValues,
		validators: {
			onSubmit: profileSchema,
		},
		onSubmit: async ({ value }) => {
			setServerError("")
			try {
				const fullName = `${value.firstName} ${value.lastName}`.trim()
				const { error } = await authClient.updateUser({
					name: fullName,
				})
				if (error) {
					setServerError(error.message ?? "Failed to update profile")
					return
				}
				router.invalidate()
			} catch (e) {
				setServerError(e instanceof Error ? e.message : "Something went wrong")
			}
		},
	})

	return (
		<form
			onSubmit={(e) => {
				e.preventDefault()
				form.handleSubmit()
			}}
		>
			<FieldGroup className="max-w-2xl gap-6">
				<div className="flex w-full gap-4">
					<form.Field name="firstName">
						{(field) => {
							const isInvalid =
								field.state.meta.isTouched && field.state.meta.errors.length > 0
							return (
								<Field data-invalid={isInvalid || undefined} className="flex-1">
									<FormFieldLabel htmlFor={field.name}>
										FIRST NAME
									</FormFieldLabel>
									<Input
										id={field.name}
										name={field.name}
										value={field.state.value}
										onBlur={field.handleBlur}
										onChange={(e) => field.handleChange(e.target.value)}
										className={inputClassName}
										aria-invalid={isInvalid || undefined}
									/>
									{isInvalid && <FieldError errors={field.state.meta.errors} />}
								</Field>
							)
						}}
					</form.Field>

					<form.Field name="lastName">
						{(field) => {
							const isInvalid =
								field.state.meta.isTouched && field.state.meta.errors.length > 0
							return (
								<Field data-invalid={isInvalid || undefined} className="flex-1">
									<FormFieldLabel htmlFor={field.name}>
										LAST NAME
									</FormFieldLabel>
									<Input
										id={field.name}
										name={field.name}
										value={field.state.value}
										onBlur={field.handleBlur}
										onChange={(e) => field.handleChange(e.target.value)}
										className={inputClassName}
										aria-invalid={isInvalid || undefined}
									/>
									{isInvalid && <FieldError errors={field.state.meta.errors} />}
								</Field>
							)
						}}
					</form.Field>
				</div>

				<Field>
					<FormFieldLabel>EMAIL ADDRESS</FormFieldLabel>
					<Input
						value={user?.email ?? ""}
						disabled
						className={inputClassName}
					/>
				</Field>

				{serverError && (
					<p className="text-xs text-destructive">{serverError}</p>
				)}

				<div className="flex items-center justify-between border-t border-border pt-6">
					<p className="font-space-grotesk text-xs font-normal text-muted-foreground">
						Changes will be saved to your profile
					</p>
					<form.Subscribe
						selector={(state) => ({
							isSubmitting: state.isSubmitting,
							values: state.values,
						})}
					>
						{({ isSubmitting, values }) => {
							const isDirty = hasChanges(values, defaultValues)
							return (
								<Button
									type="submit"
									disabled={isSubmitting || !isDirty}
									className="h-9 rounded px-5 font-space-grotesk text-xs font-bold tracking-[1px]"
								>
									{isSubmitting ? "SAVING..." : "SAVE CHANGES"}
								</Button>
							)
						}}
					</form.Subscribe>
				</div>
			</FieldGroup>
		</form>
	)
}
