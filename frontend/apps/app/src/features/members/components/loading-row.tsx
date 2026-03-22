export function LoadingRow() {
	return (
		<tr className="flex h-16 items-center px-6 border-b border-border">
			<td className="flex-1">
				<div className="flex items-center gap-3">
					<div className="h-9 w-9 animate-pulse rounded bg-skeleton" />
					<div className="flex flex-col gap-0.5">
						<div className="h-4 w-24 animate-pulse rounded bg-skeleton" />
						<div className="h-3 w-32 animate-pulse rounded bg-skeleton" />
					</div>
				</div>
			</td>
			<td className="flex w-40 items-center">
				<div className="h-6 w-24 animate-pulse rounded bg-skeleton" />
			</td>
			<td className="flex w-[140px] items-center">
				<div className="h-4 w-24 animate-pulse rounded bg-skeleton" />
			</td>
			<td className="flex w-20 items-center justify-center"></td>
		</tr>
	)
}
