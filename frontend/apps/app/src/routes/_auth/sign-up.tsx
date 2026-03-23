import { createFileRoute } from "@tanstack/react-router"
import { z } from "zod"
import { SignUpPage } from "#/features/auth/sign-up/sign-up-page"

const searchSchema = z.object({
	redirect: z.string().optional(),
})

export const Route = createFileRoute("/_auth/sign-up")({
	component: RouteComponent,
	validateSearch: (search) => searchSchema.parse(search),
})

function RouteComponent() {
	const { redirect } = Route.useSearch()
	return <SignUpPage redirect={redirect} />
}
