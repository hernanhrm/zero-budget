import { organizationClient } from "better-auth/client/plugins"
import { createAuthClient } from "better-auth/react"

const client = createAuthClient({ plugins: [organizationClient()] })

console.log(Object.keys(client.organization))
