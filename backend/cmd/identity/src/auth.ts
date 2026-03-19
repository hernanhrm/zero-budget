import { betterAuth } from "better-auth";
import { organization, twoFactor } from "better-auth/plugins";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { dash } from "@better-auth/infra";
import { db } from "./db.js";

export const auth = betterAuth({
  database: drizzleAdapter(db, {
    provider: "pg",
    usePlural: true,
  }),
  experimental: {
    joins: true, // Enable database joins for better performance
  },
  appName: "Zero Budget",
  emailAndPassword: {
    enabled: true,
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
  // Core table snake_case mappings
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
  plugins: [
    dash(),
    organization({
      allowUserToCreateOrganization: true,
      schema: {
        organization: {
          fields: {
            createdAt: "created_at",
          },
        },
        member: {
          fields: {
            organizationId: "organization_id",
            userId: "user_id",
            createdAt: "created_at",
          },
        },
        invitation: {
          fields: {
            organizationId: "organization_id",
            inviterId: "inviter_id",
            expiresAt: "expires_at",
          },
        },
      },
    }),
    twoFactor({
      issuer: "Zero Budget",
      schema: {
        twoFactor: {
          fields: {
            userId: "user_id",
          },
        },
        user: {
          fields: {
            twoFactorEnabled: "two_factor_enabled",
          },
        },
      },
    }),
  ],
});
