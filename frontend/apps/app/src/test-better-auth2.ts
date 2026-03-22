import { authClient } from "./lib/auth-client";
type Res = Awaited<ReturnType<typeof authClient.organization.listMembers>>;
type Data = NonNullable<Res["data"]>[number];

declare const d: Data;
d.user.name;
d.user.email;
d.role;
d.createdAt;
