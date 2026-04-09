import { useForm } from "@tanstack/react-form"
import { useQueryClient } from "@tanstack/react-query"
import {
	getGetV1AccountsQueryKey,
	usePostV1Accounts,
} from "@workspace/api/hooks/accounts/accounts"
import {
	getGetV1OrganizationCurrenciesQueryKey,
	useGetV1OrganizationCurrencies,
} from "@workspace/api/hooks/organization-currencies/organization-currencies"
import { Button } from "@workspace/ui/components/button"
import { DialogClose, DialogFooter } from "@workspace/ui/components/dialog"
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
import { toast } from "@workspace/ui/lib/toast"
import { Check } from "lucide-react"
import { useEffect, useLayoutEffect } from "react"

export interface AddAccountFormProps {
	open: boolean
	onComplete: () => void
	organizationId: string | undefined
}

export function AddAccountForm({
	open,
	onComplete,
	organizationId,
}: AddAccountFormProps) {
	const queryClient = useQueryClient()

	const postAccount = usePostV1Accounts({
		fetch: { credentials: "include" },
	})

	const orgCurrenciesQuery = useGetV1OrganizationCurrencies(
		{ relations: "currencies" },
		{
			query: {
				enabled: open && Boolean(organizationId),
				queryKey: getGetV1OrganizationCurrenciesQueryKey({
					relations: "currencies",
				}),
			},
			fetch: { credentials: "include" },
		},
	)

	const orgCurrencies = orgCurrenciesQuery.data?.data ?? []
	const currenciesReady = orgCurrenciesQuery.isSuccess
	const noOrgCurrencies = currenciesReady && orgCurrencies.length === 0

	const form = useForm({
		defaultValues: {
			accountName: "",
			accountType: "checking" as "checking" | "savings",
			institution: "",
			startingBalance: "0.00",
			accountNumber: "",
			currency: "",
		},
		onSubmit: async ({ value }) => {
			if (!organizationId) {
				toast.error(
					"No active organization. Open the app from your workspace and try again.",
				)
				return
			}

			const name = value.accountName.trim()
			if (name.length < 2) {
				toast.error("Account name must be at least 2 characters.")
				return
			}

			const type = value.accountType === "checking" ? "CHECKING" : "SAVINGS"
			const currencyCode = value.currency?.trim() ?? ""
			if (currencyCode.length !== 3) {
				toast.error("Select a currency for this account.")
				return
			}

			const institution = value.institution.trim()
			const accountNumber = value.accountNumber.trim()

			const balanceRaw = value.startingBalance.replace(/,/g, "").trim()
			const balanceNum = Number.parseFloat(balanceRaw === "" ? "0" : balanceRaw)
			if (!Number.isFinite(balanceNum)) {
				toast.error("Enter a valid starting balance.")
				return
			}

			try {
				const result = await postAccount.mutateAsync({
					data: {
						id: crypto.randomUUID(),
						organizationId,
						name,
						type,
						...(institution !== "" ? { institution } : {}),
						...(accountNumber !== "" ? { accountNumber } : {}),
						currencyCode,
						currentBalance: balanceNum,
						isActive: true,
					},
				})

				if (result.status !== 201) {
					toast.error(`Could not create account (${result.status}).`)
					return
				}

				await queryClient.invalidateQueries({
					queryKey: getGetV1AccountsQueryKey(),
				})
				form.reset()
				toast.success("Account created.")
				onComplete()
			} catch (e) {
				toast.error(e instanceof Error ? e.message : "Something went wrong.")
			}
		},
	})

	useEffect(() => {
		if (!open) {
			form.reset()
		}
	}, [open, form.reset])

	useLayoutEffect(() => {
		if (!open || !orgCurrenciesQuery.isSuccess) {
			return
		}
		const list = orgCurrenciesQuery.data?.data ?? []
		if (list.length === 0) {
			return
		}
		const code =
			list.find((c) => c.isBase)?.currencyCode ??
			list[0]?.currencyCode ??
			""
		if (code) {
			form.setFieldValue("currency", code)
		}
	}, [open, orgCurrenciesQuery.isSuccess, orgCurrenciesQuery.data, form.setFieldValue])

	return (
		<form
			onSubmit={(e) => {
				e.preventDefault()
				form.handleSubmit()
			}}
		>
			<div className="form-panel-body">
				<FieldGroup>
					<form.Field name="accountName">
						{(field) => {
							const isInvalid =
								field.state.meta.isTouched && field.state.meta.errors.length > 0
							return (
								<Field data-invalid={isInvalid || undefined}>
									<FormFieldLabel htmlFor={field.name}>
										ACCOUNT NAME
									</FormFieldLabel>
									<Input
										id={field.name}
										name={field.name}
										className="h-10 rounded-none px-4"
										placeholder="E.G. PRIMARY CHECKING..."
										value={field.state.value}
										onBlur={field.handleBlur}
										onChange={(e) => field.handleChange(e.target.value)}
										autoComplete="off"
										aria-invalid={isInvalid || undefined}
									/>
									{isInvalid && <FieldError errors={field.state.meta.errors} />}
								</Field>
							)
						}}
					</form.Field>

					<div className="grid gap-4 sm:grid-cols-2">
						<form.Field name="accountType">
							{(field) => (
								<Field>
									<FormFieldLabel>ACCOUNT TYPE</FormFieldLabel>
									<Select
										value={field.state.value}
										onValueChange={(value) =>
											field.handleChange(value as "checking" | "savings")
										}
									>
										<SelectTrigger
											size="default"
											className="h-10 w-full rounded-none px-4"
										>
											<SelectValue />
										</SelectTrigger>
										<SelectContent>
											<SelectItem value="checking">CHECKING</SelectItem>
											<SelectItem value="savings">SAVINGS</SelectItem>
										</SelectContent>
									</Select>
								</Field>
							)}
						</form.Field>
						<form.Field name="institution">
							{(field) => {
								const isInvalid =
									field.state.meta.isTouched &&
									field.state.meta.errors.length > 0
								return (
									<Field data-invalid={isInvalid || undefined}>
										<FormFieldLabel htmlFor={field.name}>
											INSTITUTION
										</FormFieldLabel>
										<Input
											id={field.name}
											name={field.name}
											className="h-10 rounded-none px-4"
											placeholder="BANK NAME"
											value={field.state.value}
											onBlur={field.handleBlur}
											onChange={(e) => field.handleChange(e.target.value)}
											autoComplete="organization"
											aria-invalid={isInvalid || undefined}
										/>
										{isInvalid && (
											<FieldError errors={field.state.meta.errors} />
										)}
									</Field>
								)
							}}
						</form.Field>
					</div>

					<form.Subscribe
						selector={(state) => ({ currencyCode: state.values.currency })}
					>
						{({ currencyCode }) => {
							const sym =
								orgCurrencies.find((c) => c.currencyCode === currencyCode)
									?.currency?.symbol ?? "$"
							return (
								<form.Field name="startingBalance">
									{(field) => {
										const isInvalid =
											field.state.meta.isTouched &&
											field.state.meta.errors.length > 0
										return (
											<Field data-invalid={isInvalid || undefined}>
												<FormFieldLabel htmlFor={field.name}>
													STARTING BALANCE
												</FormFieldLabel>
												<div className="flex h-12 items-center gap-3 border border-ring px-4">
													<span className="font-space-grotesk text-2xl font-bold text-primary">
														{sym}
													</span>
													<Input
														id={field.name}
														name={field.name}
														className="h-auto flex-1 border-0 bg-transparent p-0 font-ibm-plex-mono text-2xl font-bold text-foreground shadow-none focus-visible:ring-0"
														value={field.state.value}
														onBlur={field.handleBlur}
														onChange={(e) => field.handleChange(e.target.value)}
														inputMode="decimal"
														aria-invalid={isInvalid || undefined}
													/>
												</div>
												{isInvalid && (
													<FieldError errors={field.state.meta.errors} />
												)}
											</Field>
										)
									}}
								</form.Field>
							)
						}}
					</form.Subscribe>

					<form.Field name="accountNumber">
						{(field) => (
							<Field>
								<FormFieldLabel htmlFor={field.name}>
									ACCOUNT NUMBER (OPTIONAL)
								</FormFieldLabel>
								<Input
									id={field.name}
									name={field.name}
									className="h-10 rounded-none px-4"
									placeholder="****1234"
									value={field.state.value}
									onBlur={field.handleBlur}
									onChange={(e) => field.handleChange(e.target.value)}
									autoComplete="off"
								/>
							</Field>
						)}
					</form.Field>

					<form.Field name="currency">
						{(field) => (
							<Field>
								<FormFieldLabel>CURRENCY</FormFieldLabel>
								<Select
									value={field.state.value}
									onValueChange={(value) => field.handleChange(value)}
									disabled={
										orgCurrenciesQuery.isLoading ||
										noOrgCurrencies ||
										orgCurrencies.length === 0 ||
										(currenciesReady && !field.state.value)
									}
								>
									<SelectTrigger
										size="default"
										className="h-10 w-full rounded-none px-4"
									>
										<SelectValue
											placeholder={
												orgCurrenciesQuery.isLoading
													? "Loading currencies…"
													: noOrgCurrencies
														? "No currencies for this organization"
														: "Select currency"
											}
										/>
									</SelectTrigger>
									<SelectContent>
										{orgCurrencies.map((c) => {
											const code = c.currencyCode ?? ""
											const name = c.currency?.name ?? code
											const symbol = c.currency?.symbol ?? ""
											const base = c.isBase ? " — base" : ""
											return (
												<SelectItem key={code} value={code}>
													{`${code} — ${name}${symbol ? ` (${symbol})` : ""}${base}`}
												</SelectItem>
											)
										})}
									</SelectContent>
								</Select>
							</Field>
						)}
					</form.Field>
				</FieldGroup>
			</div>

			<DialogFooter className="border-t border-border px-6 py-4 sm:justify-end">
				<DialogClose asChild>
					<Button
						type="button"
						variant="outline"
						className="rounded-none border-border text-muted-foreground"
					>
						CANCEL
					</Button>
				</DialogClose>
				<form.Subscribe
					selector={(state) => ({
						isSubmitting: state.isSubmitting,
						currency: state.values.currency,
					})}
				>
					{({ isSubmitting, currency }) => (
						<Button
							type="submit"
							disabled={
								isSubmitting ||
								postAccount.isPending ||
								!organizationId ||
								orgCurrenciesQuery.isLoading ||
								noOrgCurrencies ||
								orgCurrencies.length === 0 ||
								!currency
							}
							className="gap-2 rounded-none"
						>
							<Check className="size-3.5" />
							ADD ACCOUNT
						</Button>
					)}
				</form.Subscribe>
			</DialogFooter>
		</form>
	)
}
