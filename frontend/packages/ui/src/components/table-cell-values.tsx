import { cn } from "@workspace/ui/lib/utils"
import type * as React from "react"

function TableCellValueStack({
	className,
	primary,
	secondary,
	...props
}: React.ComponentProps<"div"> & {
	primary: React.ReactNode
	secondary: React.ReactNode
}) {
	return (
		<div
			data-slot="table-cell-value-stack"
			className={cn("flex min-w-0 flex-col gap-0.5", className)}
			{...props}
		>
			<TableCellValuePrimary>{primary}</TableCellValuePrimary>
			<TableCellValueSecondary>{secondary}</TableCellValueSecondary>
		</div>
	)
}

function TableCellValuePrimary({
	className,
	...props
}: React.ComponentProps<"span">) {
	return (
		<span
			data-slot="table-cell-value-primary"
			className={cn(
				"font-space-grotesk text-[13px] font-bold tracking-[1px] text-foreground",
				className,
			)}
			{...props}
		/>
	)
}

function TableCellValueSecondary({
	className,
	...props
}: React.ComponentProps<"span">) {
	return (
		<span
			data-slot="table-cell-value-secondary"
			className={cn(
				"font-ibm-plex-mono text-[10px] tracking-[1px] text-muted-foreground",
				className,
			)}
			{...props}
		/>
	)
}

function TableCellMonoValue({
	className,
	...props
}: React.ComponentProps<"span">) {
	return (
		<span
			data-slot="table-cell-mono-value"
			className={cn(
				"font-ibm-plex-mono text-xs tracking-[1px] text-muted-foreground",
				className,
			)}
			{...props}
		/>
	)
}

function TableCellAmountValue({
	className,
	emphasis = "default",
	...props
}: React.ComponentProps<"span"> & {
	emphasis?: "default" | "primary"
}) {
	return (
		<span
			data-slot="table-cell-amount-value"
			className={cn(
				"font-ibm-plex-mono text-[13px] font-bold tracking-[1px]",
				emphasis === "primary" ? "text-primary" : "text-foreground",
				className,
			)}
			{...props}
		/>
	)
}

function TableCellStatusBadge({
	active,
	className,
	...props
}: React.ComponentProps<"div"> & { active: boolean }) {
	return (
		<div
			data-slot="table-cell-status-badge"
			className={cn("flex h-6 items-center justify-center", className)}
			{...props}
		>
			<div className="flex items-center gap-1.5">
				<span
					className={cn(
						"size-1.5 shrink-0 rounded-full",
						active ? "bg-success" : "bg-muted-foreground",
					)}
				/>
				<span
					className={cn(
						"font-space-grotesk text-[10px] font-bold tracking-[1px]",
						active ? "text-success" : "text-muted-foreground",
					)}
				>
					{active ? "ACTIVE" : "INACTIVE"}
				</span>
			</div>
		</div>
	)
}

export {
	TableCellAmountValue,
	TableCellMonoValue,
	TableCellStatusBadge,
	TableCellValuePrimary,
	TableCellValueSecondary,
	TableCellValueStack,
}
