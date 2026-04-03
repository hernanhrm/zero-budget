export function MetricCards() {
	return (
		<div className="grid w-full gap-5 md:grid-cols-3">
			<div className="flex flex-col gap-4 border border-border p-6">
				<p className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
					TOTAL BALANCE
				</p>
				<p className="font-space-grotesk text-4xl font-bold text-foreground">
					$12,847.52
				</p>
				<div className="flex items-center gap-2">
					<span className="font-space-grotesk text-[10px] font-bold text-success">
						▲
					</span>
					<span className="font-ibm-plex-mono text-[11px] tracking-[1px] text-success">
						+3.2% FROM LAST MONTH
					</span>
				</div>
			</div>
			<div className="flex flex-col gap-4 border border-border p-6">
				<p className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
					CHECKING ACCOUNTS
				</p>
				<p className="font-space-grotesk text-4xl font-bold text-foreground">
					$8,234.18
				</p>
				<p className="font-ibm-plex-mono text-[11px] tracking-[1px] text-muted-foreground">
					2 ACCOUNTS
				</p>
			</div>
			<div className="flex flex-col gap-4 border border-border p-6">
				<p className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
					SAVINGS ACCOUNTS
				</p>
				<p className="font-space-grotesk text-4xl font-bold text-primary">
					$4,613.34
				</p>
				<p className="font-ibm-plex-mono text-[11px] tracking-[1px] text-muted-foreground">
					1 ACCOUNT
				</p>
			</div>
		</div>
	)
}
