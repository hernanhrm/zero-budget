import { useForm } from "@tanstack/react-form"
import type { Currency } from "@workspace/api"
import { usePostV1OrganizationCurrencies } from "@workspace/api/hooks/organization-currencies/organization-currencies"
import { Button } from "@workspace/ui/components/button"
import { Dialog, DialogContent } from "@workspace/ui/components/dialog"
import { DialogPanelHeader } from "@workspace/ui/components/dialog-panel-header"
import { Field, FieldGroup } from "@workspace/ui/components/field"
import { FormFieldLabel } from "@workspace/ui/components/form-field-label"
import { Input } from "@workspace/ui/components/input"
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@workspace/ui/components/select"
import { toast } from "@workspace/ui/lib/toast"
import { useEffect, useMemo } from "react"

export interface AddCurrencyModalProps {
	open: boolean
	onOpenChange: (open: boolean) => void
	organizationId: string | undefined
	/** Catalog currencies not yet linked to the org. */
	availableCurrencies: Currency[]
	/** When true, this will be the first org currency (base, rate 1). */
	isFirstOrgCurrency: boolean
	onSuccess: () => void
}

export function AddCurrencyModal({
	open,
	onOpenChange,
	organizationId,
	availableCurrencies,
	isFirstOrgCurrency,
	onSuccess,
}: AddCurrencyModalProps) {
	const postCurrency = usePostV1OrganizationCurrencies({
		fetch: { credentials: "include" },
	})

	const sortedCatalog = useMemo(() => {
		return [...availableCurrencies].sort((a, b) =>
			(a.code ?? "").localeCompare(b.code ?? ""),
		)
	}, [availableCurrencies])

	const form = useForm({
		defaultValues: {
			currencyCode: "",
			rate: "1",
		},
		onSubmit: async ({ value }) => {
			if (!organizationId) {
				toast.error(
					"No active organization. Open the app from your workspace and try again.",
				)
				return
			}

			const code = value.currencyCode?.trim().toUpperCase() ?? ""
			if (code.length !== 3) {
				toast.error("Select a currency.")
				return
			}

			const rateRaw = value.rate.replace(/,/g, "").trim()
			const rateNum = Number.parseFloat(rateRaw === "" ? "0" : rateRaw)
			if (!Number.isFinite(rateNum)) {
				toast.error("Enter a valid exchange rate.")
				return
			}

			if (isFirstOrgCurrency) {
				if (rateNum !== 1) {
					toast.error("The first currency must be base; rate must be 1.")
					return
				}
			} else if (rateNum <= 0) {
				toast.error("Rate must be greater than zero.")
				return
			}

			try {
				const res = await postCurrency.mutateAsync({
					data: {
						id: crypto.randomUUID(),
						organizationId,
						currencyCode: code,
						isBase: isFirstOrgCurrency,
						rate: rateNum,
					},
				})
				if (res.status !== 201) {
					toast.error("Could not add currency.")
					return
				}
				toast.success("Currency added.")
				onSuccess()
				form.reset()
				onOpenChange(false)
			} catch {
				toast.error("Could not add currency.")
			}
		},
	})

	useEffect(() => {
		if (!open) return
		const first = sortedCatalog[0]?.code ?? ""
		form.setFieldValue("currencyCode", first)
		form.setFieldValue("rate", isFirstOrgCurrency ? "1" : "")
	}, [open, sortedCatalog, isFirstOrgCurrency, form.setFieldValue])

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent className="gap-0 p-0 sm:max-w-md" showCloseButton={false}>
				<DialogPanelHeader
					title="ADD CURRENCY"
					titleClassName="text-sm uppercase"
					description={
						isFirstOrgCurrency
							? "ADD YOUR ORGANIZATION BASE CURRENCY (RATE 1)"
							: "ADD A CURRENCY AND ITS RATE VS THE BASE"
					}
				/>

				<form
					onSubmit={(e) => {
						e.preventDefault()
						form.handleSubmit()
					}}
				>
					<FieldGroup className="form-panel-body">
						<form.Field name="currencyCode">
							{(field) => (
								<Field>
									<FormFieldLabel>CURRENCY</FormFieldLabel>
									<Select
										value={field.state.value}
										onValueChange={(v) => field.handleChange(v)}
										disabled={sortedCatalog.length === 0}
									>
										<SelectTrigger className="h-10 w-full px-4 font-space-grotesk text-sm">
											<SelectValue placeholder="Select currency" />
										</SelectTrigger>
										<SelectContent>
											{sortedCatalog.map((c) => {
												const code = c.code ?? ""
												const name = c.name ?? code
												const sym = c.symbol ?? ""
												return (
													<SelectItem key={code} value={code}>
														{`${name} (${code}) ${sym}`}
													</SelectItem>
												)
											})}
										</SelectContent>
									</Select>
								</Field>
							)}
						</form.Field>

						<form.Field name="rate">
							{(field) => (
								<Field>
									<FormFieldLabel>
										{isFirstOrgCurrency
											? "RATE (BASE = 1)"
											: "RATE (PER 1 BASE UNIT)"}
									</FormFieldLabel>
									<Input
										placeholder={isFirstOrgCurrency ? "1" : "e.g. 0.92"}
										value={field.state.value}
										onChange={(e) => field.handleChange(e.target.value)}
										className="h-10 px-4 font-ibm-plex-mono text-sm"
										type="text"
										inputMode="decimal"
										disabled={isFirstOrgCurrency}
									/>
								</Field>
							)}
						</form.Field>

						<Button
							type="submit"
							disabled={
								postCurrency.isPending ||
								sortedCatalog.length === 0 ||
								!organizationId
							}
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
