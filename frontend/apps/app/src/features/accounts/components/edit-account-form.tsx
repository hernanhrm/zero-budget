import { useForm } from "@tanstack/react-form"
import { useQueryClient } from "@tanstack/react-query"
import type { Account, OrganizationCurrency } from "@workspace/api"
import {
	getGetV1AccountsQueryKey,
	usePutV1AccountsId,
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
import { Switch } from "@workspace/ui/components/switch"
import { toast } from "@workspace/ui/lib/toast"
import { Check } from "lucide-react"
import { useEffect, useMemo } from "react"

function apiProblemDetail(data: unknown): string | undefined {
	if (data !== null && typeof data === "object" && "detail" in data) {
		const d = (data as { detail?: unknown }).detail
		return typeof d === "string" ? d : undefined
	}
	return undefined
}

function typeToForm(t: string | null | undefined): "checking" | "savings" {
	const u = (t ?? "").toUpperCase()
	return u === "SAVINGS" ? "savings" : "checking"
}

export interface EditAccountFormProps {
	open: boolean
	account: Account
	onComplete: () => void
}

export function EditAccountForm({
	open,
	account,
	onComplete,
}: EditAccountFormProps) {
	const queryClient = useQueryClient()

	const putAccount = usePutV1AccountsId({
		fetch: { credentials: "include" },
	})

	const orgId = account.organizationId
	const orgCurrenciesQuery = useGetV1OrganizationCurrencies(
		{ relations: "currencies" },
		{
			query: {
				enabled: open && Boolean(orgId),
				queryKey: getGetV1OrganizationCurrenciesQueryKey({
					relations: "currencies",
				}),
			},
			fetch: { credentials: "include" },
		},
	)
	const orgCurrencies = useMemo(() => {
		const raw = orgCurrenciesQuery.data?.data
		return Array.isArray(raw) ? raw : []
	}, [orgCurrenciesQuery.data])

	const currenciesReady = orgCurrenciesQuery.isSuccess

	const form = useForm({
		defaultValues: {
			accountName: account.name ?? "",
			accountType: typeToForm(account.type),
			institution: account.institution ?? "",
			accountNumber: account.accountNumber ?? "",
			isActive: account.isActive ?? true,
			currency: account.currencyCode ?? "",
		},
		onSubmit: async ({ value }) => {
			const id = account.id
			if (!id) {
				toast.error("Missing account id.")
				return
			}

			const name = value.accountName.trim()
			if (name.length < 2) {
				toast.error("Account name must be at least 2 characters.")
				return
			}

			const type = value.accountType === "checking" ? "CHECKING" : "SAVINGS"
			const institution = value.institution.trim()
			const accountNumber = value.accountNumber.trim()

			try {
				const currencyCode = value.currency?.trim() ?? ""
				if (currencyCode.length !== 3) {
					toast.error("Select a currency for this account.")
					return
				}

				const result = await putAccount.mutateAsync({
					id,
					data: {
						name,
						type,
						institution,
						accountNumber,
						isActive: value.isActive,
						currencyCode,
					},
				})

				if (result.status === 204) {
					await queryClient.invalidateQueries({
						queryKey: getGetV1AccountsQueryKey(),
					})
					toast.success("Account updated.")
					onComplete()
					return
				}

				const msg =
					apiProblemDetail(result.data) ??
					`Could not update account (${result.status}).`
				toast.error(msg)
			} catch (e) {
				toast.error(e instanceof Error ? e.message : "Something went wrong.")
			}
		},
	})

	useEffect(() => {
		if (!open) {
			return
		}
		form.reset({
			accountName: account.name ?? "",
			accountType: typeToForm(account.type),
			institution: account.institution ?? "",
			accountNumber: account.accountNumber ?? "",
			isActive: account.isActive ?? true,
			currency: account.currencyCode ?? "",
		})
	}, [open, account, form.reset])

	const orgOptions: OrganizationCurrency[] = (() => {
		const list = [...orgCurrencies]
		const code = account.currencyCode
		if (code && !list.some((c) => c.currencyCode === code)) {
			const synthetic: OrganizationCurrency = {
				currencyCode: code,
				isBase: false,
				currency: null,
			}
			list.unshift(synthetic)
		}
		return list
	})()

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
							{(field) => (
								<Field>
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
									/>
								</Field>
							)}
						</form.Field>
					</div>

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
										orgOptions.length === 0 ||
										(currenciesReady && orgCurrencies.length > 0 && !field.state.value)
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
													: orgOptions.length === 0
														? "No currencies for this organization"
														: "Select currency"
											}
										/>
									</SelectTrigger>
									<SelectContent>
										{orgOptions.map((c) => {
											const itemCode = c.currencyCode ?? ""
											const name = c.currency?.name ?? itemCode
											const symbol = c.currency?.symbol ?? ""
											const base = c.isBase ? " — base" : ""
											return (
												<SelectItem key={itemCode} value={itemCode}>
													{`${itemCode} — ${name}${symbol ? ` (${symbol})` : ""}${base}`}
												</SelectItem>
											)
										})}
									</SelectContent>
								</Select>
							</Field>
						)}
					</form.Field>

					<form.Field name="isActive">
						{(field) => (
							<Field className="flex flex-row items-center justify-between gap-4 rounded-none border border-border px-4 py-3">
								<div className="space-y-0.5">
									<FormFieldLabel className="text-base">
										ACCOUNT ACTIVE
									</FormFieldLabel>
									<p className="font-ibm-plex-mono text-xs text-muted-foreground">
										Inactive accounts stay in your history but can be hidden
										from primary views.
									</p>
								</div>
								<Switch
									checked={field.state.value}
									onCheckedChange={(v) => field.handleChange(v)}
								/>
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
								putAccount.isPending ||
								orgCurrenciesQuery.isLoading ||
								orgOptions.length === 0 ||
								!currency
							}
							className="gap-2 rounded-none"
						>
							<Check className="size-3.5" />
							SAVE CHANGES
						</Button>
					)}
				</form.Subscribe>
			</DialogFooter>
		</form>
	)
}
