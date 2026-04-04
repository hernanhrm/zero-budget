import { TanStackDevtools } from "@tanstack/react-devtools"
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools"

/**
 * Dev-only UI; imported lazily from __root so production builds skip this module entirely.
 */
export default function RootDevtools() {
	return (
		<TanStackDevtools
			config={{
				position: "bottom-right",
			}}
			plugins={[
				{
					name: "TanStack Router",
					render: <TanStackRouterDevtoolsPanel />,
				},
			]}
		/>
	)
}
