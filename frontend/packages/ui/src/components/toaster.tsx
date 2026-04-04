import { cn } from "@workspace/ui/lib/utils"
import { Toaster as Sonner } from "sonner"

export function Toaster() {
	return (
		<Sonner
			position="top-center"
			richColors
			closeButton
			className="toaster group"
			toastOptions={{
				classNames: {
					toast: cn(
						"group pointer-events-auto flex w-full max-w-md items-start gap-3 border border-border bg-popover p-4 text-popover-foreground shadow-md",
						"font-space-grotesk text-sm data-[type=error]:border-destructive/50 data-[type=error]:bg-destructive/10",
					),
					title: "font-bold tracking-[0.5px]",
					description: "text-muted-foreground",
					actionButton:
						"rounded-none border border-border bg-background px-3 py-1.5 text-xs font-bold tracking-[1px]",
					cancelButton:
						"rounded-none px-3 py-1.5 text-xs font-bold tracking-[1px] text-muted-foreground",
					closeButton:
						"absolute top-3 right-3 border-0 bg-transparent text-foreground/60 hover:text-foreground",
				},
			}}
		/>
	)
}
