import { createAuthClient } from "better-auth/react"

export const authClient = createAuthClient({
	baseURL: import.meta.env.VITE_IDENTITY_URL ?? "http://localhost:8081",
})
