import { useForm } from "@tanstack/react-form"
import { Button } from "@workspace/ui/components/button"
import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogFooter,
} from "@workspace/ui/components/dialog"
import { DialogPanelHeader } from "@workspace/ui/components/dialog-panel-header"
import { Field, FieldError, FieldGroup } from "@workspace/ui/components/field"
import { FormFieldLabel } from "@workspace/ui/components/form-field-label"
import { Input } from "@workspace/ui/components/input"
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@workspace/ui/components/select"
import { Textarea } from "@workspace/ui/components/textarea"
import { Send } from "lucide-react"
import { useState } from "react"
import { authClient } from "#/lib/auth-client"
import { inviteSchema } from "../invite-schema"

interface InviteMemberModalProps {
	open: boolean
	onOpenChange: (open: boolean) => void
	onSuccess: () => void
}

export function InviteMemberModal({
	open,
	onOpenChange,
	onSuccess,
}: InviteMemberModalProps) {
	const [serverError, setServerError] = useState("")

	const form = useForm({
		defaultValues: {
			email: "",
			role: "member" as "member" | "admin" | "owner",
			message: "",
		},
		validators: {
			onSubmit: inviteSchema,
		},
		onSubmit: async ({ value }) => {
			setServerError("")
			try {
				const { error } = await authClient.organization.inviteMember({
					email: value.email,
					role: value.role,
				})
				if (error) {
					setServerError(error.message ?? "Something went wrong")
					return
				}
				form.reset()
				setServerError("")
				onSuccess()
				onOpenChange(false)
			} catch (e) {
				setServerError(e instanceof Error ? e.message : "Something went wrong")
			}
		},
	})

	return (
		<Dialog
			open={open}
			onOpenChange={(value) => {
				if (!value) {
					form.reset()
					setServerError("")
				}
				onOpenChange(value)
			}}
		>
			<DialogContent
				className="gap-0 p-0 sm:max-w-[480px]"
				showCloseButton={false}
			>
				<DialogPanelHeader title="INVITE MEMBER" titleClassName="text-sm" />

				<form
					onSubmit={(e) => {
						e.preventDefault()
						form.handleSubmit()
					}}
				>
					<div className="form-panel-body">
						<FieldGroup>
							<form.Field name="email">
								{(field) => {
									const isInvalid =
										field.state.meta.isTouched &&
										field.state.meta.errors.length > 0
									return (
										<Field data-invalid={isInvalid || undefined}>
											<FormFieldLabel htmlFor={field.name} size="sm">
												EMAIL ADDRESS
											</FormFieldLabel>
											<Input
												id={field.name}
												name={field.name}
												type="email"
												placeholder="ENTER EMAIL ADDRESS"
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
							</form.Field>

							<form.Field name="role">
								{(field) => (
									<Field>
										<FormFieldLabel size="sm">ROLE</FormFieldLabel>
										<Select
											value={field.state.value}
											onValueChange={(value) =>
												field.handleChange(
													value as "member" | "admin" | "owner",
												)
											}
										>
											<SelectTrigger className="w-full">
												<SelectValue />
											</SelectTrigger>
											<SelectContent>
												<SelectItem value="member">MEMBER</SelectItem>
												<SelectItem value="admin">ADMIN</SelectItem>
												<SelectItem value="owner">OWNER</SelectItem>
											</SelectContent>
										</Select>
									</Field>
								)}
							</form.Field>

							<form.Field name="message">
								{(field) => (
									<Field>
										<FormFieldLabel size="sm">
											MESSAGE (OPTIONAL)
										</FormFieldLabel>
										<Textarea
											id={field.name}
											name={field.name}
											placeholder="Add a personal note..."
											value={field.state.value}
											onBlur={field.handleBlur}
											onChange={(e) => field.handleChange(e.target.value)}
											rows={3}
										/>
									</Field>
								)}
							</form.Field>

							{serverError && (
								<p className="text-xs text-destructive">{serverError}</p>
							)}
						</FieldGroup>
					</div>

					<DialogFooter className="border-t border-border px-6 py-4">
						<DialogClose asChild>
							<Button type="button" variant="outline">
								CANCEL
							</Button>
						</DialogClose>
						<form.Subscribe selector={(state) => state.isSubmitting}>
							{(isSubmitting) => (
								<Button type="submit" disabled={isSubmitting} className="gap-2">
									<Send className="size-3.5" />
									{isSubmitting ? "SENDING..." : "SEND INVITE"}
								</Button>
							)}
						</form.Subscribe>
					</DialogFooter>
				</form>
			</DialogContent>
		</Dialog>
	)
}
