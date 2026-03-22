import { createFileRoute } from "@tanstack/react-router"
import { Button } from "@workspace/ui/components/button"

export const Route = createFileRoute("/_protected/")({
	component: App,
})

function App() {
	return <Button>Hello</Button>
}
