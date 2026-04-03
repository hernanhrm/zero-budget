import { useForm } from "@tanstack/react-form"
import { useRouter } from "@tanstack/react-router"
import { Button } from "@workspace/ui/components/button"
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
} from "@workspace/ui/components/card"
import { Field } from "@workspace/ui/components/field"
import { FormFieldLabel } from "@workspace/ui/components/form-field-label"
import { Input } from "@workspace/ui/components/input"
import { Download, Trash2, TriangleAlert } from "lucide-react"
import { useState } from "react"
import { authClient } from "#/lib/auth-client"

export function DangerZoneTab() {
	const router = useRouter()
	const [serverError, setServerError] = useState("")

	const form = useForm({
		defaultValues: {
			confirmation: "",
		},
		onSubmit: async ({ value }) => {
			setServerError("")
			if (value.confirmation !== "DELETE") {
				setServerError("Please type DELETE to confirm")
				return
			}
			try {
				await authClient.deleteUser()
				router.navigate({ to: "/sign-in" })
			} catch (e) {
				setServerError(
					e instanceof Error ? e.message : "Failed to delete account",
				)
			}
		},
	})

	return (
		<div className="flex flex-col gap-8">
			{/* Warning Banner */}
			<div className="flex items-center gap-3 border border-[#FF6B35] bg-[#2A1A1A] px-5 py-4">
				<TriangleAlert className="size-5 shrink-0 text-[#FF6B35]" />
				<span className="font-ibm-plex-mono text-xs tracking-[1px] text-[#FF6B35]">
					ACTIONS IN THIS SECTION ARE IRREVERSIBLE. PLEASE PROCEED WITH CAUTION.
				</span>
			</div>

			{/* Delete Account Section */}
			<Card className="gap-0 border-[#FF6B35] p-0 ring-[#FF6B35]">
				<CardHeader className="border-b border-[#FF6B35] bg-[#2A1A1A] px-6 py-4">
					<div className="flex items-center gap-3">
						<div className="h-5 w-1 bg-[#FF6B35]" />
						<CardTitle className="font-space-grotesk text-sm font-bold tracking-[1px] text-[#FF6B35]">
							DELETE ACCOUNT
						</CardTitle>
					</div>
				</CardHeader>
				<CardContent className="flex flex-col gap-6 p-6">
					<p className="font-ibm-plex-mono text-xs leading-[1.6] tracking-[0.5px] text-foreground">
						ONCE YOU DELETE YOUR ACCOUNT, THERE IS NO GOING BACK. THIS WILL
						PERMANENTLY DELETE YOUR ACCOUNT, ALL YOUR BUDGETS, TRANSACTIONS, AND
						REMOVE ALL ORGANIZATION MEMBERSHIPS.
					</p>

					<form
						onSubmit={(e) => {
							e.preventDefault()
							form.handleSubmit()
						}}
					>
						<div className="flex flex-col gap-6">
							<form.Field name="confirmation">
								{(field) => (
									<Field>
										<FormFieldLabel htmlFor={field.name}>
											TYPE &quot;DELETE&quot; TO CONFIRM
										</FormFieldLabel>
										<Input
											id={field.name}
											name={field.name}
											value={field.state.value}
											onBlur={field.handleBlur}
											onChange={(e) => field.handleChange(e.target.value)}
											placeholder="DELETE"
											className="h-11 border-[#FF6B35] px-4 font-space-grotesk text-sm"
										/>
									</Field>
								)}
							</form.Field>

							{serverError && (
								<p className="text-xs text-destructive">{serverError}</p>
							)}

							<div className="flex items-center justify-between">
								<span className="font-ibm-plex-mono text-[10px] tracking-[1px] text-muted-foreground">
									THIS ACTION CANNOT BE UNDONE
								</span>
								<form.Subscribe
									selector={(state) => ({
										isSubmitting: state.isSubmitting,
										confirmation: state.values.confirmation,
									})}
								>
									{({ isSubmitting, confirmation }) => (
										<Button
											type="submit"
											disabled={isSubmitting || confirmation !== "DELETE"}
											className="gap-2 bg-[#FF6B35] font-space-grotesk text-xs font-bold tracking-[1px] text-white hover:bg-[#FF6B35]/90"
										>
											<Trash2 data-icon="inline-start" />
											{isSubmitting ? "DELETING..." : "DELETE MY ACCOUNT"}
										</Button>
									)}
								</form.Subscribe>
							</div>
						</div>
					</form>
				</CardContent>
			</Card>

			{/* Export Data Section */}
			<Card className="gap-0 p-0">
				<CardHeader className="border-b border-border bg-[#232323] px-6 py-4">
					<div className="flex items-center gap-3">
						<div className="h-5 w-1 bg-primary" />
						<CardTitle className="font-space-grotesk text-sm font-bold tracking-[1px] text-foreground">
							EXPORT DATA
						</CardTitle>
					</div>
				</CardHeader>
				<CardContent className="flex items-center justify-between p-6">
					<div className="flex flex-col gap-1">
						<span className="font-space-grotesk text-[13px] font-bold tracking-[1px] text-foreground">
							DOWNLOAD YOUR DATA
						</span>
						<span className="font-ibm-plex-mono text-[10px] tracking-[1px] text-muted-foreground">
							EXPORT ALL YOUR BUDGETS, TRANSACTIONS AND ACCOUNT DATA AS CSV
						</span>
					</div>
					<Button
						variant="outline"
						className="gap-2 border-primary font-space-grotesk text-xs font-bold tracking-[1px] text-primary hover:bg-primary/10 hover:text-primary"
					>
						<Download data-icon="inline-start" />
						EXPORT DATA
					</Button>
				</CardContent>
			</Card>
		</div>
	)
}
