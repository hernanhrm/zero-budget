import type { ErrorComponentProps } from "@tanstack/react-router"
import { Button } from "@workspace/ui/components/button"
import { toast } from "@workspace/ui/lib/toast"
import { useEffect } from "react"

export function AccountsRouteError({ error, reset }: ErrorComponentProps) {
	useEffect(() => {
		const message =
			error instanceof Error ? error.message : "Failed to load accounts."
		toast.error(message)
	}, [error])

	return (
		<div className="flex h-full flex-col gap-6 overflow-auto p-10">
			<p className="font-space-grotesk text-sm text-destructive">
				{error instanceof Error ? error.message : "Something went wrong."}
			</p>
			<div>
				<Button
					type="button"
					variant="outline"
					className="rounded-none"
					onClick={() => reset()}
				>
					TRY AGAIN
				</Button>
			</div>
		</div>
	)
}
