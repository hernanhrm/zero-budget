import { cn } from "@workspace/ui/lib/utils"
import type { ReactNode } from "react"

export interface ModulePageHeaderProps {
	title: string
	description: string
	children?: ReactNode
	className?: string
}

export function ModulePageHeader({
	title,
	description,
	children,
	className,
}: ModulePageHeaderProps) {
	return (
		<div
			className={cn(
				"flex w-full flex-col gap-8 sm:flex-row sm:items-center sm:justify-between",
				className,
			)}
		>
			<div className="flex flex-col gap-3">
				<h1 className="font-space-grotesk text-4xl font-bold tracking-[1px] text-foreground">
					{title}
				</h1>
				<p className="font-ibm-plex-mono text-[13px] tracking-[1px] text-muted-foreground">
					{description}
				</p>
			</div>
			{children != null ? (
				<div className="flex shrink-0 items-center self-start sm:self-auto">
					{children}
				</div>
			) : null}
		</div>
	)
}
