import { useForm } from "@tanstack/react-form"
import { Button } from "@workspace/ui/components/button"
import { Dialog, DialogContent } from "@workspace/ui/components/dialog"
import { DialogPanelHeader } from "@workspace/ui/components/dialog-panel-header"
import { Field, FieldGroup } from "@workspace/ui/components/field"
import { FormFieldLabel } from "@workspace/ui/components/form-field-label"
import { Input } from "@workspace/ui/components/input"

interface AddCurrencyModalProps {
	open: boolean
	onOpenChange: (open: boolean) => void
	onAdd: (currency: { name: string; code: string; rate: string }) => void
}

export function AddCurrencyModal({
	open,
	onOpenChange,
	onAdd,
}: AddCurrencyModalProps) {
	const form = useForm({
		defaultValues: {
			name: "",
			code: "",
			rate: "",
		},
		onSubmit: async ({ value }) => {
			if (!value.name || !value.code || !value.rate) return
			onAdd({
				name: value.name.toUpperCase(),
				code: value.code.toUpperCase(),
				rate: value.rate,
			})
			form.reset()
			onOpenChange(false)
		},
	})

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent className="gap-0 p-0 sm:max-w-md" showCloseButton={false}>
				<DialogPanelHeader
					title="ADD CURRENCY"
					titleClassName="text-sm uppercase"
					description="ADD A NEW CURRENCY AND EXCHANGE RATE"
				/>

				<form
					onSubmit={(e) => {
						e.preventDefault()
						form.handleSubmit()
					}}
				>
					<FieldGroup className="form-panel-body">
						<form.Field name="name">
							{(field) => (
								<Field>
									<FormFieldLabel>CURRENCY NAME</FormFieldLabel>
									<Input
										placeholder="e.g. Japanese Yen"
										value={field.state.value}
										onChange={(e) => field.handleChange(e.target.value)}
										className="h-10 px-4 font-space-grotesk text-sm"
									/>
								</Field>
							)}
						</form.Field>

						<form.Field name="code">
							{(field) => (
								<Field>
									<FormFieldLabel>CURRENCY CODE</FormFieldLabel>
									<Input
										placeholder="e.g. JPY"
										value={field.state.value}
										onChange={(e) => field.handleChange(e.target.value)}
										className="h-10 px-4 font-space-grotesk text-sm uppercase"
										maxLength={3}
									/>
								</Field>
							)}
						</form.Field>

						<form.Field name="rate">
							{(field) => (
								<Field>
									<FormFieldLabel>EXCHANGE RATE</FormFieldLabel>
									<Input
										placeholder="e.g. 149.50"
										value={field.state.value}
										onChange={(e) => field.handleChange(e.target.value)}
										className="h-10 px-4 font-ibm-plex-mono text-sm"
										type="number"
										step="any"
									/>
								</Field>
							)}
						</form.Field>

						<Button
							type="submit"
							className="w-full font-space-grotesk text-xs font-bold tracking-[1px]"
						>
							ADD CURRENCY
						</Button>
					</FieldGroup>
				</form>
			</DialogContent>
		</Dialog>
	)
}
