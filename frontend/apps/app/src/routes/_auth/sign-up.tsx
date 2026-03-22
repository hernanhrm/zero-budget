import { createFileRoute } from "@tanstack/react-router"
import { SignUpPage } from "#/features/auth/sign-up/sign-up-page"

export const Route = createFileRoute("/_auth/sign-up")({
	component: SignUpPage,
})
