import { authClient } from "./lib/auth-client";
console.log(authClient.useActiveOrganization);
console.log(authClient.organization.useActiveOrganization);
authClient.organization.listMembers({ query: { limit: 50 } }).then(res => res.data);
