import { appendFileSync } from "node:fs";
import { Hono } from "hono";
import { cors } from "hono/cors";
import { serve } from "@hono/node-server";
import { eq } from "drizzle-orm";
import { auth } from "./auth.js";
import { db } from "./db.js";
import { invitations } from "./schema.js";
import { publishOrganizationInvitationCreated } from "./lib/events.js";

// #region agent log
const AGENT_DEBUG_LOG =
  "/Users/hernanreyes/Documents/dev/github.com/hernanhrm/zero-budget/.cursor/debug-741baa.log";

function agentDebugLog(entry: {
  hypothesisId: string;
  location: string;
  message: string;
  data: Record<string, unknown>;
  runId?: string;
}) {
  const line =
    JSON.stringify({
      sessionId: "741baa",
      timestamp: Date.now(),
      runId: entry.runId ?? "pre-fix",
      ...entry,
    }) + "\n";
  try {
    appendFileSync(AGENT_DEBUG_LOG, line);
  } catch {
    /* Fly.io / read-only fs */
  }
  console.error("AGENT_DEBUG_SESSION", line.trim());
}

function betterAuthHostHint(): string | null {
  const raw = process.env.BETTER_AUTH_URL;
  if (!raw?.trim()) return null;
  try {
    return new URL(raw.trim()).host;
  } catch {
    return "invalid-url";
  }
}
// #endregion

function parseCommaOrigins(value: string | undefined): string[] {
  if (!value?.trim()) {
    return [];
  }
  return value.split(",").map((o) => o.trim()).filter(Boolean);
}

const trustedOriginsParsed = parseCommaOrigins(process.env.TRUSTED_ORIGINS);
const trustedOriginsList =
  trustedOriginsParsed.length > 0
    ? trustedOriginsParsed
    : ["http://localhost:3000"];

const corsOverride = parseCommaOrigins(process.env.CORS_ORIGIN);
const corsAllowedOrigins =
  corsOverride.length > 0 ? corsOverride : trustedOriginsList;

const app = new Hono();

const corsConfig = cors({
  origin:
    corsAllowedOrigins.length === 1
      ? corsAllowedOrigins[0]
      : corsAllowedOrigins,
  allowHeaders: ["Content-Type", "Authorization"],
  allowMethods: ["POST", "GET", "DELETE", "OPTIONS"],
  exposeHeaders: ["Content-Length"],
  maxAge: 600,
  credentials: true,
});

app.use("/api/auth/*", corsConfig);

app.use("/api/invitations/*", corsConfig);

app.on(["POST", "GET"], "/api/auth/*", async (c) => {
  const req = c.req.raw;
  const url = new URL(req.url);
  const pathname = url.pathname;
  const cookieHeader = req.headers.get("cookie") ?? "";
  const origin = req.headers.get("origin");
  const host = req.headers.get("host");

  const isGetSessionProbe =
    req.method === "GET" &&
    (pathname.includes("get-session") || pathname.endsWith("/session"));

  const isSignInPost =
    req.method === "POST" &&
    (pathname.includes("sign-in") || pathname.includes("sign_in"));

  // #region agent log
  let preHandlerSessionNull: boolean | undefined;
  if (isGetSessionProbe) {
    const session = await auth.api.getSession({ headers: req.headers });
    preHandlerSessionNull = session == null;
  }

  if (isGetSessionProbe || isSignInPost) {
    agentDebugLog({
      hypothesisId: isGetSessionProbe ? "H1-samesite-cookie" : "H5-sign-in",
      location: "identity/index.ts:api/auth",
      message: isGetSessionProbe ? "before_handler_get_session" : "before_handler_sign_in",
      data: {
        pathname,
        method: req.method,
        hasAnyCookie: cookieHeader.length > 0,
        cookieMentionsPrefix: cookieHeader.includes("zero-budget"),
        origin,
        host,
        requestHost: host,
        betterAuthHostEnv: betterAuthHostHint(),
        preHandlerGetSessionNull: preHandlerSessionNull,
        cookiePolicyExpectedSameSite: (
          process.env.BETTER_AUTH_URL ?? ""
        ).startsWith("https://")
          ? "none"
          : "lax",
      },
    });
  }
  // #endregion

  const res = await auth.handler(req);

  // #region agent log
  if (isSignInPost) {
    const setCookie = res.headers.get("set-cookie");
    agentDebugLog({
      hypothesisId: "H5-set-cookie",
      location: "identity/index.ts:api/auth",
      message: "after_handler_sign_in",
      data: {
        pathname,
        status: res.status,
        setCookiePresent: setCookie != null && setCookie.length > 0,
        origin,
        host,
        betterAuthHostEnv: betterAuthHostHint(),
      },
    });
  }
  // #endregion

  return res;
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
// Bind all interfaces so Docker / platform health checks can reach the server.
const hostname = process.env.HOST ?? "0.0.0.0";

console.log(`Identity server listening on http://${hostname}:${port}`);

serve({ fetch: app.fetch, port, hostname });
