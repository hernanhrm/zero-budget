const originalEnv = process.env;
globalThis.import = { meta: { env: { VITE_IDENTITY_URL: "http://localhost:8081" } } } as any;

import { authClient } from "./frontend/apps/app/src/lib/auth-client"

console.log(typeof authClient.organization.activeOrganization)
