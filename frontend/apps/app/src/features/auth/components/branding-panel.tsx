export function BrandingPanel() {
	return (
		<>
			<div className="flex items-center gap-2">
				<div className="flex size-8 items-center justify-center bg-primary font-mono text-sm font-bold text-primary-foreground">
					Z
				</div>
				<span className="font-mono text-xs font-bold uppercase tracking-wider">
					Zero Budget
				</span>
			</div>

			<div className="space-y-6">
				<h2 className="font-mono text-4xl font-bold uppercase leading-tight tracking-wider">
					Take control
					<br />
					of your
					<br />
					finances
				</h2>
				<div className="space-y-3">
					<p className="font-mono text-xs uppercase tracking-wider text-neutral-400">
						Track every currency. Manage every account.
						<br />
						Master every category.
					</p>
					<div className="h-0.5 w-16 bg-primary" />
				</div>
			</div>

			<div className="flex items-center justify-between font-mono text-[10px] uppercase tracking-wider text-neutral-400">
				<div className="flex items-center gap-2">
					<span className="inline-block size-1.5 bg-primary" />
					<span>System online</span>
				</div>
				<span>v0.1.0</span>
			</div>
		</>
	)
}
