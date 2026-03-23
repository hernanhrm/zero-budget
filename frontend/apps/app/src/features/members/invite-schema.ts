import { z } from "zod"

export const inviteSchema = z.object({
	email: z
		.string()
		.min(1, "Email is required")
		.email("Invalid email address"),
	role: z.enum(["member", "admin", "owner"]),
	message: z.string().optional(),
})
