import { cn } from "@workspace/ui/lib/utils"
import type * as React from "react"

export interface DataTableSectionHeaderProps
	extends Omit<React.ComponentProps<"div">, "title"> {
	title: React.ReactNode
	/** Shown as mono count when `countSlot` is not set */
	count?: number
	/** When set, replaces the default mono count */
	countSlot?: React.ReactNode
	/** Right-aligned actions (e.g. primary button) */
	endSlot?: React.ReactNode
}

function DataTableSectionHeader({
	className,
	title,
	count,
	countSlot,
	endSlot,
	...props
}: DataTableSectionHeaderProps) {
	const showDefaultCount = countSlot === undefined && typeof count === "number"

	return (
		<div
			data-slot="data-table-section-header"
			className={cn(
				"flex h-14 items-center justify-between gap-4 border-b border-border bg-card px-6",
				className,
			)}
			{...props}
		>
			<div className="flex min-w-0 items-center gap-3">
				<div className="h-5 w-1 shrink-0 bg-primary" />
				<span className="font-space-grotesk text-sm font-bold tracking-[1px] text-foreground">
					{title}
				</span>
				{countSlot}
				{showDefaultCount ? (
					<span className="font-ibm-plex-mono text-xs text-muted-foreground">
						{count}
					</span>
				) : null}
			</div>
			{endSlot}
		</div>
	)
}

export { DataTableSectionHeader }
