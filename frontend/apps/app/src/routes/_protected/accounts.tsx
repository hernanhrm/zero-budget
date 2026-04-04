import { createFileRoute } from "@tanstack/react-router"
import { AccountsPage } from "#/features/accounts/accounts-page"
import { AccountsRouteError } from "#/features/accounts/components/accounts-route-error"

export const Route = createFileRoute("/_protected/accounts")({
	component: AccountsPage,
	errorComponent: AccountsRouteError,
})
