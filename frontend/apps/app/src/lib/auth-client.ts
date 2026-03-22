import { createAuthClient } from "better-auth/react"
import { organizationClient } from "better-auth/client/plugins"

export const authClient = createAuthClient({
	baseURL: import.meta.env.VITE_IDENTITY_URL ?? "http://localhost:8081",
	plugins: [organizationClient()],
})
