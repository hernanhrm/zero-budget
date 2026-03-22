import { createFileRoute } from "@tanstack/react-router"
import { z } from "zod"
import { ResetPasswordPage } from "#/features/auth/reset-password/reset-password-page"

const searchSchema = z.object({
	token: z.string().optional(),
	error: z.string().optional(),
})

export const Route = createFileRoute("/reset-password")({
	component: ResetPasswordPage,
	validateSearch: (search) => searchSchema.parse(search),
})
