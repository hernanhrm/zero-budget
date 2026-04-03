import { cn } from "@workspace/ui/lib/utils"
import type * as React from "react"

function TableColumnHeader({
	className,
	...props
}: React.ComponentProps<"span">) {
	return (
		<span
			data-slot="table-column-header"
			className={cn(
				"font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground",
				className,
			)}
			{...props}
		/>
	)
}

export { TableColumnHeader }
