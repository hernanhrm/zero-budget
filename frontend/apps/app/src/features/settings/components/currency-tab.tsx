import { Trash2 } from "lucide-react"
import { useState } from "react"
import { Button } from "@workspace/ui/components/button"
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@workspace/ui/components/select"
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@workspace/ui/components/table"
import { AddCurrencyModal } from "./add-currency-modal"

interface ExchangeRate {
	id: string
	currency: string
	code: string
	rate: string
}

const initialRates: ExchangeRate[] = [
	{ id: "1", currency: "EURO", code: "EUR", rate: "0.9265" },
	{ id: "2", currency: "BRITISH POUND", code: "GBP", rate: "0.7905" },
	{ id: "3", currency: "MEXICAN PESO", code: "MXN", rate: "17.3180" },
	{ id: "4", currency: "CANADIAN DOLLAR", code: "CAD", rate: "1.3645" },
]

const currencies = [
	{ value: "USD", label: "USD — US DOLLAR ($)" },
	{ value: "EUR", label: "EUR — EURO (€)" },
	{ value: "GBP", label: "GBP — BRITISH POUND (£)" },
	{ value: "MXN", label: "MXN — MEXICAN PESO ($)" },
	{ value: "CAD", label: "CAD — CANADIAN DOLLAR ($)" },
]

const displayFormats = [
	{ value: "symbol-first", label: "$1,234.56" },
	{ value: "code-first", label: "USD 1,234.56" },
	{ value: "symbol-last", label: "1,234.56$" },
]

export function CurrencyTab() {
	const [rates, setRates] = useState(initialRates)
	const [primaryCurrency, setPrimaryCurrency] = useState("USD")
	const [displayFormat, setDisplayFormat] = useState("symbol-first")
	const [addModalOpen, setAddModalOpen] = useState(false)

	const handleDeleteRate = (id: string) => {
		setRates((prev) => prev.filter((r) => r.id !== id))
	}

	return (
		<div className="flex flex-col gap-8">
			<div className="flex gap-6">
				<div className="flex flex-1 flex-col gap-2">
					<span className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
						PRIMARY CURRENCY
					</span>
					<Select value={primaryCurrency} onValueChange={setPrimaryCurrency}>
						<SelectTrigger className="h-11 w-full px-4 font-space-grotesk text-sm">
							<SelectValue />
						</SelectTrigger>
						<SelectContent>
							{currencies.map((c) => (
								<SelectItem key={c.value} value={c.value}>
									{c.label}
								</SelectItem>
							))}
						</SelectContent>
					</Select>
				</div>

				<div className="flex flex-1 flex-col gap-2">
					<span className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
						DISPLAY FORMAT
					</span>
					<Select value={displayFormat} onValueChange={setDisplayFormat}>
						<SelectTrigger className="h-11 w-full px-4 font-space-grotesk text-sm">
							<SelectValue />
						</SelectTrigger>
						<SelectContent>
							{displayFormats.map((f) => (
								<SelectItem key={f.value} value={f.value}>
									{f.label}
								</SelectItem>
							))}
						</SelectContent>
					</Select>
				</div>
			</div>

			<div className="flex flex-col gap-4">
				<div className="flex items-center justify-between">
					<span className="font-space-grotesk text-sm font-bold tracking-[1px] text-foreground">
						EXCHANGE RATES
					</span>
					<Button
						variant="outline"
						size="sm"
						onClick={() => setAddModalOpen(true)}
						className="border-primary font-space-grotesk text-[10px] font-bold tracking-[1px] text-primary hover:bg-primary/10 hover:text-primary"
					>
						<span className="text-sm">+</span>
						ADD CURRENCY
					</Button>
				</div>

				<Table>
					<TableHeader>
						<TableRow className="bg-[#232323]">
							<TableHead className="font-space-grotesk text-[11px] font-bold tracking-[1px]">
								CURRENCY
							</TableHead>
							<TableHead className="w-[100px] font-space-grotesk text-[11px] font-bold tracking-[1px]">
								CODE
							</TableHead>
							<TableHead className="w-[120px] text-right font-space-grotesk text-[11px] font-bold tracking-[1px]">
								RATE
							</TableHead>
							<TableHead className="w-10" />
						</TableRow>
					</TableHeader>
					<TableBody>
						{rates.map((rate) => (
							<TableRow key={rate.id}>
								<TableCell className="font-space-grotesk text-[13px] font-semibold text-foreground">
									{rate.currency}
								</TableCell>
								<TableCell className="font-ibm-plex-mono text-xs text-muted-foreground">
									{rate.code}
								</TableCell>
								<TableCell className="text-right font-ibm-plex-mono text-[13px] text-foreground">
									{rate.rate}
								</TableCell>
								<TableCell>
									<Button
										variant="ghost"
										size="icon-sm"
										onClick={() => handleDeleteRate(rate.id)}
										className="text-muted-foreground hover:text-destructive"
									>
										<Trash2 />
									</Button>
								</TableCell>
							</TableRow>
						))}
					</TableBody>
				</Table>

				<span className="font-ibm-plex-mono text-[10px] tracking-[1px] text-[#3D3D3D]">
					RATES LAST UPDATED: MARCH 23, 2026 AT 09:00 UTC
				</span>
			</div>

			<AddCurrencyModal
				open={addModalOpen}
				onOpenChange={setAddModalOpen}
				onAdd={(currency) => {
					setRates((prev) => [
						...prev,
						{
							id: crypto.randomUUID(),
							currency: currency.name,
							code: currency.code,
							rate: currency.rate,
						},
					])
				}}
			/>
		</div>
	)
}
