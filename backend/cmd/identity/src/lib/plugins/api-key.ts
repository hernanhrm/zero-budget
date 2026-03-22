import { apiKey } from "@better-auth/api-key";

export const apiKeyPlugin = apiKey(undefined, {
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
});
