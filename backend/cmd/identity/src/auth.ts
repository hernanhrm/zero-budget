import { betterAuth } from "better-auth";
import { openAPI, organization, twoFactor } from "better-auth/plugins";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { dash } from "@better-auth/infra";
import { apiKey } from "@better-auth/api-key";
import { db } from "./db.js";

export const auth = betterAuth({
  database: drizzleAdapter(db, {
    provider: "pg",
    usePlural: true,
  }),
  experimental: {
    joins: true, // Enable database joins for better performance
  },
  advanced: {
    ipAddress: {
      // For Cloudflare
      ipAddressHeaders: ["cf-connecting-ip", "x-forwarded-for"],

      // For Vercel
      // ipAddressHeaders: ["x-vercel-forwarded-for", "x-forwarded-for"],

      // For AWS/Generic
      // ipAddressHeaders: ["x-forwarded-for"],
    },
  },
  appName: "Zero Budget",
  emailAndPassword: {
    enabled: true,
    requireEmailVerification: true,
  },
  emailVerification: {
    sendOnSignUp: true,
    autoSignInAfterVerification: true,
    async afterEmailVerification(user) {
      const goApiUrl = process.env.GO_API_URL || "http://localhost:8080";
      const internalApiKey = process.env.INTERNAL_API_KEY;
      if (!internalApiKey) return;

      try {
        await fetch(`${goApiUrl}/v1/events/publish`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "X-API-Key": internalApiKey,
          },
          body: JSON.stringify({
            event: "user.signed_up",
            payload: {
              userId: user.id,
              email: user.email,
              name: user.name,
            },
          }),
        });
      } catch (err) {
        console.error("Failed to publish user.signed_up event:", err);
      }
    },
    sendVerificationEmail: async ({ user, url }: { user: { id: string; email: string; name: string }; url: string }) => {
      const goApiUrl = process.env.GO_API_URL || "http://localhost:8080";
      const internalApiKey = process.env.INTERNAL_API_KEY;

      console.log("[identity] sendVerificationEmail triggered", {
        userId: user.id,
        email: user.email,
        name: user.name,
        url,
        goApiUrl,
        hasApiKey: !!internalApiKey,
      });

      if (!internalApiKey) {
        console.warn("[identity] INTERNAL_API_KEY not set, skipping event publish");
        return;
      }

      try {
        const eventPayload = {
          event: "user.verification_email",
          payload: {
            userId: user.id,
            email: user.email,
            name: user.name,
            verificationUrl: url,
          },
        };
        console.log("[identity] publishing event to Go API", eventPayload);

        const res = await fetch(`${goApiUrl}/v1/events/publish`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "X-API-Key": internalApiKey,
          },
          body: JSON.stringify(eventPayload),
        });

        console.log("[identity] Go API response", {
          status: res.status,
          statusText: res.statusText,
        });

        if (!res.ok) {
          const body = await res.text();
          console.error("[identity] Go API error response body:", body);
        }
      } catch (err) {
        console.error("[identity] Failed to publish verification email event:", err);
      }
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
    openAPI(),
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
    apiKey(undefined, {
      schema: {
        apikey: {
          fields: {
            configId: "config_id",
            referenceId: "reference_id",
            createdAt: "created_at",
            updatedAt: "updated_at",
            expiresAt: "expires_at",
            rateLimitEnabled: "rate_limit_enabled",
            rateLimitTimeWindow: "rate_limit_time_window",
            rateLimitMax: "rate_limit_max",
            requestCount: "request_count",
            lastRefillAt: "last_refill_at",
            lastRequest: "last_request",
            refillAmount: "refill_amount",
            refillInterval: "refill_interval",
          },
        },
      },
    }),
  ],

});
