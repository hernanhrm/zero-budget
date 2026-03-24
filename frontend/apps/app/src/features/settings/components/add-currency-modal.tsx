import { useForm } from "@tanstack/react-form"
import { Button } from "@workspace/ui/components/button"
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	DialogDescription,
} from "@workspace/ui/components/dialog"
import {
	Field,
	FieldGroup,
	FieldLabel,
} from "@workspace/ui/components/field"
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
			<DialogContent className="sm:max-w-md">
				<DialogHeader>
					<DialogTitle className="font-space-grotesk text-sm font-bold uppercase tracking-[1px]">
						ADD CURRENCY
					</DialogTitle>
					<DialogDescription className="font-ibm-plex-mono text-[10px] uppercase tracking-[1px]">
						ADD A NEW CURRENCY AND EXCHANGE RATE
					</DialogDescription>
				</DialogHeader>

				<form
					onSubmit={(e) => {
						e.preventDefault()
						form.handleSubmit()
					}}
				>
					<FieldGroup>
						<form.Field
							name="name"
							children={(field) => (
								<Field>
									<FieldLabel className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
										CURRENCY NAME
									</FieldLabel>
									<Input
										placeholder="e.g. Japanese Yen"
										value={field.state.value}
										onChange={(e) => field.handleChange(e.target.value)}
										className="h-10 px-4 font-space-grotesk text-sm"
									/>
								</Field>
							)}
						/>

						<form.Field
							name="code"
							children={(field) => (
								<Field>
									<FieldLabel className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
										CURRENCY CODE
									</FieldLabel>
									<Input
										placeholder="e.g. JPY"
										value={field.state.value}
										onChange={(e) => field.handleChange(e.target.value)}
										className="h-10 px-4 font-space-grotesk text-sm uppercase"
										maxLength={3}
									/>
								</Field>
							)}
						/>

						<form.Field
							name="rate"
							children={(field) => (
								<Field>
									<FieldLabel className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
										EXCHANGE RATE
									</FieldLabel>
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
						/>

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
