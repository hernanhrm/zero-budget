import { betterAuth } from "better-auth";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { db } from "./db.js";
import { plugins } from "./lib/plugins/index.js";
import {
  publishPasswordReset,
  publishVerificationEmail,
  publishUserSignedUp,
} from "./lib/events.js";
import { organizations, members, sessions } from "./schema.js";
import { eq } from "drizzle-orm";

export const auth = betterAuth({
  trustedOrigins: (
    process.env.TRUSTED_ORIGINS ?? "http://localhost:3000"
  ).split(","),
  database: drizzleAdapter(db, {
    provider: "pg",
    usePlural: true,
  }),
  experimental: {
    joins: true,
  },
  advanced: {
    cookiePrefix: "zero-budget",
    defaultCookieAttributes: {
      sameSite: "lax",
      path: "/",
      secure: false,
    },
    ipAddress: {
      ipAddressHeaders: ["cf-connecting-ip", "x-forwarded-for"],
    },
  },
  appName: "Zero Budget",
  emailAndPassword: {
    enabled: true,
    requireEmailVerification: false,
    sendResetPassword: async ({ user, url }) => {
      await publishPasswordReset(user, url);
    },
  },
  emailVerification: {
    sendOnSignUp: true,
    autoSignInAfterVerification: true,
    async afterEmailVerification(user) {
      await publishUserSignedUp(user);
    },
    sendVerificationEmail: async ({ user, url }) => {
      await publishVerificationEmail(user, url);
    },
  },
  socialProviders: {
    ...(process.env.GOOGLE_CLIENT_ID && {
      google: {
        clientId: process.env.GOOGLE_CLIENT_ID,
        clientSecret: process.env.GOOGLE_CLIENT_SECRET!,
      },
    }),
    ...(process.env.GITHUB_CLIENT_ID && {
      github: {
        clientId: process.env.GITHUB_CLIENT_ID,
        clientSecret: process.env.GITHUB_CLIENT_SECRET!,
      },
    }),
  },
  user: {
    fields: {
      emailVerified: "email_verified",
      createdAt: "created_at",
      updatedAt: "updated_at",
    },
  },
  session: {
    fields: {
      userId: "user_id",
      expiresAt: "expires_at",
      ipAddress: "ip_address",
      userAgent: "user_agent",
      createdAt: "created_at",
      updatedAt: "updated_at",
    },
  },
  account: {
    fields: {
      userId: "user_id",
      accountId: "account_id",
      providerId: "provider_id",
      accessToken: "access_token",
      refreshToken: "refresh_token",
      accessTokenExpiresAt: "access_token_expires_at",
      refreshTokenExpiresAt: "refresh_token_expires_at",
      idToken: "id_token",
      createdAt: "created_at",
      updatedAt: "updated_at",
    },
  },
  verification: {
    fields: {
      expiresAt: "expires_at",
      createdAt: "created_at",
      updatedAt: "updated_at",
    },
  },
  databaseHooks: {
    user: {
      create: {
        after: async (user, context) => {
          if (!context) return;

          const baseName = user.name ?? user.email ?? "user";
          const slugPart = baseName
            .toLowerCase()
            .replace(/\s+/g, "-")
            .replace(/[^a-z0-9-]/g, "");
          const randomSuffix = Math.random().toString(36).substring(2, 8);
          const slug = `${slugPart}-${randomSuffix}`;
          const orgId = crypto.randomUUID();
          const memberId = crypto.randomUUID();
          const now = new Date();

          await db.insert(organizations).values({
            id: orgId,
            name: baseName,
            slug,
            created_at: now,
          });

          await db.insert(members).values({
            id: memberId,
            organization_id: orgId,
            user_id: user.id,
            role: "owner",
            created_at: now,
          });

          await db.update(sessions)
            .set({ active_organization_id: orgId })
            .where(eq(sessions.user_id, user.id));
        },
      },
    },
    session: {
      create: {
        after: async (session) => {
          if (session.activeOrganizationId) return;

          const userMemberships = await db.select()
            .from(members)
            .where(eq(members.user_id, session.userId))
            .limit(1);

          if (userMemberships.length > 0) {
            await db.update(sessions)
              .set({ active_organization_id: userMemberships[0].organization_id })
              .where(eq(sessions.id, session.id));
          }
        },
      },
    },
  },
  plugins,
});
