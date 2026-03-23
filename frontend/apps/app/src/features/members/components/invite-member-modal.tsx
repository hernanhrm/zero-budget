import { useForm } from "@tanstack/react-form"
import { Button } from "@workspace/ui/components/button"
import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@workspace/ui/components/dialog"
import {
	Field,
	FieldError,
	FieldGroup,
	FieldLabel,
} from "@workspace/ui/components/field"
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
				setServerError(
					e instanceof Error ? e.message : "Something went wrong",
				)
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
			<DialogContent className="sm:max-w-[480px] gap-0 p-0">
				<DialogHeader className="border-b border-border px-6 py-4">
					<DialogTitle className="font-space-grotesk text-sm font-bold tracking-[1px]">
						INVITE MEMBER
					</DialogTitle>
				</DialogHeader>

				<form
					onSubmit={(e) => {
						e.preventDefault()
						form.handleSubmit()
					}}
				>
					<div className="px-6 py-5">
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
												className="font-space-grotesk text-[10px] font-bold tracking-[1px] text-muted-foreground"
											>
												EMAIL ADDRESS
											</FieldLabel>
											<Input
												id={field.name}
												name={field.name}
												type="email"
												placeholder="ENTER EMAIL ADDRESS"
												value={field.state.value}
												onBlur={field.handleBlur}
												onChange={(e) =>
													field.handleChange(e.target.value)
												}
												aria-invalid={isInvalid || undefined}
											/>
											{isInvalid && (
												<FieldError
													errors={field.state.meta.errors}
												/>
											)}
										</Field>
									)
								}}
							/>

							<form.Field
								name="role"
								children={(field) => (
									<Field>
										<FieldLabel className="font-space-grotesk text-[10px] font-bold tracking-[1px] text-muted-foreground">
											ROLE
										</FieldLabel>
										<Select
											value={field.state.value}
											onValueChange={(value) =>
												field.handleChange(
													value as
														| "member"
														| "admin"
														| "owner",
												)
											}
										>
											<SelectTrigger className="w-full">
												<SelectValue />
											</SelectTrigger>
											<SelectContent>
												<SelectItem value="member">
													MEMBER
												</SelectItem>
												<SelectItem value="admin">
													ADMIN
												</SelectItem>
												<SelectItem value="owner">
													OWNER
												</SelectItem>
											</SelectContent>
										</Select>
									</Field>
								)}
							/>

							<form.Field
								name="message"
								children={(field) => (
									<Field>
										<FieldLabel className="font-space-grotesk text-[10px] font-bold tracking-[1px] text-muted-foreground">
											MESSAGE (OPTIONAL)
										</FieldLabel>
										<Textarea
											id={field.name}
											name={field.name}
											placeholder="Add a personal note..."
											value={field.state.value}
											onBlur={field.handleBlur}
											onChange={(e) =>
												field.handleChange(e.target.value)
											}
											rows={3}
										/>
									</Field>
								)}
							/>

							{serverError && (
								<p className="text-xs text-destructive">
									{serverError}
								</p>
							)}
						</FieldGroup>
					</div>

					<DialogFooter className="border-t border-border px-6 py-4">
						<DialogClose asChild>
							<Button
								type="button"
								variant="outline"
								className="font-space-grotesk text-xs font-bold tracking-[1px]"
							>
								CANCEL
							</Button>
						</DialogClose>
						<form.Subscribe
							selector={(state) => state.isSubmitting}
							children={(isSubmitting) => (
								<Button
									type="submit"
									disabled={isSubmitting}
									className="gap-2 font-space-grotesk text-xs font-bold tracking-[1px]"
								>
									<Send className="size-3.5" />
									{isSubmitting
										? "SENDING..."
										: "SEND INVITE"}
								</Button>
							)}
						/>
					</DialogFooter>
				</form>
			</DialogContent>
		</Dialog>
	)
}
