import { createFileRoute } from "@tanstack/react-router"
import { z } from "zod"
import { SignInPage } from "#/features/auth/sign-in/sign-in-page"

const searchSchema = z.object({
	redirect: z.string().optional(),
})

export const Route = createFileRoute("/_auth/sign-in")({
	component: RouteComponent,
	validateSearch: (search) => searchSchema.parse(search),
})

function RouteComponent() {
	const { redirect } = Route.useSearch()
	return <SignInPage redirect={redirect} />
}
