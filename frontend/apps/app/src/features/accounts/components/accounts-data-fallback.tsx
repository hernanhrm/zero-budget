import { DataTableSectionHeader } from "@workspace/ui/components/data-table-section-header"

export function AccountsDataFallback() {
	return (
		<>
			<div className="grid w-full gap-5 md:grid-cols-3">
				{[1, 2, 3].map((k) => (
					<div key={k} className="flex flex-col gap-4 border border-border p-6">
						<p className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
							…
						</p>
						<p className="font-space-grotesk text-4xl font-bold text-muted-foreground">
							…
						</p>
					</div>
				))}
			</div>

			<div className="min-w-0 overflow-x-auto">
				<div className="w-full border border-border">
					<DataTableSectionHeader title="ALL ACCOUNTS" count={0} />
					<p className="px-6 py-8 font-space-grotesk text-sm text-muted-foreground">
						LOADING ACCOUNTS…
					</p>
				</div>
			</div>
		</>
	)
}
