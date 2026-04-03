import { Cancel01Icon } from "@hugeicons/core-free-icons"
import { HugeiconsIcon } from "@hugeicons/react"
import { Button } from "@workspace/ui/components/button"
import {
	DialogClose,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from "@workspace/ui/components/dialog"
import { cn } from "@workspace/ui/lib/utils"
import type * as React from "react"

export interface DialogPanelHeaderProps {
	title: React.ReactNode
	titleClassName?: string
	className?: string
	description?: React.ReactNode
	descriptionClassName?: string
}

function DialogPanelHeader({
	title,
	titleClassName,
	className,
	description,
	descriptionClassName,
}: DialogPanelHeaderProps) {
	const hasDescription = description != null && description !== false

	return (
		<DialogHeader
			data-slot="dialog-panel-header"
			className={cn(
				"gap-0 space-y-0 border-b border-border bg-card p-0 text-left",
				className,
			)}
		>
			<div className="flex items-center justify-between px-6 py-4">
				<div className="flex items-center gap-3">
					<div className="h-5 w-1 bg-primary" />
					<DialogTitle
						className={cn(
							"font-space-grotesk text-base font-bold tracking-[1px]",
							titleClassName,
						)}
					>
						{title}
					</DialogTitle>
				</div>
				<DialogClose asChild>
					<Button
						type="button"
						variant="outline"
						size="icon"
						className="size-8 shrink-0 rounded-none border-border"
					>
						<HugeiconsIcon
							icon={Cancel01Icon}
							strokeWidth={2}
							className="size-4 text-muted-foreground"
						/>
						<span className="sr-only">Close</span>
					</Button>
				</DialogClose>
			</div>
			{hasDescription ? (
				<DialogDescription
					className={cn(
						"px-6 pb-4 pt-0 font-ibm-plex-mono text-[10px] uppercase tracking-[1px] text-muted-foreground",
						descriptionClassName,
					)}
				>
					{description}
				</DialogDescription>
			) : null}
		</DialogHeader>
	)
}

export { DialogPanelHeader }
