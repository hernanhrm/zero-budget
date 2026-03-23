import { Hono } from "hono";
import { cors } from "hono/cors";
import { serve } from "@hono/node-server";
import { eq } from "drizzle-orm";
import { auth } from "./auth.js";
import { db } from "./db.js";
import { invitations } from "./schema.js";
import { publishOrganizationInvitationCreated } from "./lib/events.js";

const app = new Hono();

const corsConfig = cors({
  origin: process.env.CORS_ORIGIN ?? "http://localhost:3000",
  allowHeaders: ["Content-Type", "Authorization"],
  allowMethods: ["POST", "GET", "OPTIONS"],
  exposeHeaders: ["Content-Length"],
  maxAge: 600,
  credentials: true,
});

app.use("/api/auth/*", corsConfig);

app.use("/api/invitations/*", corsConfig);

app.on(["POST", "GET"], "/api/auth/*", (c) => {
  return auth.handler(c.req.raw);
});

app.post("/api/invitations/:invitation-id/resend", async (c) => {
  const session = await auth.api.getSession({ headers: c.req.raw.headers });
  if (!session) {
    return c.json({ error: "Unauthorized" }, 401);
  }

  const invitationId = c.req.param("invitation-id");

  const invitation = await db.query.invitations.findFirst({
    where: eq(invitations.id, invitationId),
    with: {
      users: true,
      organizations: true,
    },
  });

  if (!invitation) {
    return c.json({ error: "Invitation not found" }, 404);
  }

  if (invitation.status !== "pending") {
    return c.json({ error: "Invitation is no longer pending" }, 400);
  }

  if (new Date(invitation.expires_at) < new Date()) {
    return c.json({ error: "Invitation has expired" }, 400);
  }

  const inviter = invitation.users;
  const org = invitation.organizations;
  const inviterName = inviter.name || inviter.email;
  const baseUrl = process.env.APP_URL || "http://localhost:3000";

  await publishOrganizationInvitationCreated({
    email: invitation.email,
    inviterName,
    inviterEmail: inviter.email,
    inviterInitial: inviterName.charAt(0).toUpperCase(),
    organizationName: org.name,
    acceptUrl: `${baseUrl}/invite/accept/${invitation.id}`,
    declineUrl: `${baseUrl}/invite/decline/${invitation.id}`,
  });

  return c.json({ success: true });
});

app.get("/health", (c) => {
  return c.json({ status: "ok" });
});

const port = Number(process.env.PORT) || 8081;

console.log(`Identity server running on port ${port}`);

serve({ fetch: app.fetch, port });
