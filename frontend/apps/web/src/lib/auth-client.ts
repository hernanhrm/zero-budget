import { createAuthClient } from "better-auth/react"
import { organizationClient, twoFactorClient } from "better-auth/client/plugins"
import { dashClient } from "@better-auth/infra/client"

export const authClient = createAuthClient({
  baseURL: import.meta.env.VITE_AUTH_URL ?? "http://localhost:8081",
  plugins: [
    organizationClient(),
    twoFactorClient(),
    dashClient(),
  ],
})
